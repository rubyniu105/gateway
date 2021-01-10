/*
Copyright 2016 Medcl (m AT medcl.net)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package indexing

import (
	"bytes"
	"fmt"
	log "github.com/cihub/seelog"
	"infini.sh/framework/core/elastic"
	"infini.sh/framework/core/param"
	"infini.sh/framework/core/pipeline"
	"infini.sh/framework/core/queue"
	"infini.sh/framework/core/stats"
	"sync"
	"time"
)

type JsonIndexingJoint struct {
	param.Parameters
	inputQueueName string
}

//处理纯 json 格式的消息索引
func (joint JsonIndexingJoint) Name() string {
	return "json_indexing"
}

//TODO 重启子进程，当子进程挂了之后
func (joint JsonIndexingJoint) Process(c *pipeline.Context) error {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("error in process: %s\n", err)
		}
	}()

	workers, _ := joint.GetInt("worker_size", 1)
	bulkSizeInMB, _ := joint.GetInt("bulk_size_in_mb", 10)
	joint.inputQueueName = joint.GetStringOrDefault("input_queue", "es_queue")
	bulkSizeInMB = 1000000 * bulkSizeInMB

	wg := sync.WaitGroup{}
	totalSize := 0
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go joint.NewBulkWorker(&totalSize, bulkSizeInMB, &wg)
	}

	wg.Wait()

	return nil
}

func (joint JsonIndexingJoint) NewBulkWorker(count *int, bulkSizeInMB int, wg *sync.WaitGroup) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("error in bulk worker: %s", err)
			//TODO failure and save logs for later recovery
			wg.Done()
		}
	}()

	log.Trace("start bulk worker")

	mainBuf := bytes.Buffer{}
	docBuf := bytes.Buffer{}

	destIndex := joint.GetStringOrDefault("index_name", "")
	destType := joint.GetStringOrDefault("index_type", "_doc")
	esInstanceVal := joint.GetStringOrDefault("elasticsearch", "es_json_bulk")

	timeOut := joint.GetIntOrDefault("idle_timeout_in_second", 5)
	idleDuration := time.Duration(timeOut) * time.Second
	idleTimeout := time.NewTimer(idleDuration)
	defer idleTimeout.Stop()

	client := elastic.GetClient(esInstanceVal)

READ_DOCS:
	for {
		idleTimeout.Reset(idleDuration)

		select {

		case pop := <-queue.ReadChan(joint.inputQueueName):

			stats.IncrementBy("bulk", "received", int64(mainBuf.Len()))

			//TODO record ingest time,	request.LoggingTime = time.Now().UTC().Format("2006-01-02T15:04:05.000Z")

			docBuf.WriteString(fmt.Sprintf("{ \"index\" : { \"_index\" : \"%s\", \"_type\" : \"%s\" } }\n", destIndex, destType))
			docBuf.Write(pop)
			docBuf.WriteString("\n")

			mainBuf.Write(docBuf.Bytes())

			docBuf.Reset()
			(*count)++

			if mainBuf.Len() > (bulkSizeInMB) {
				log.Trace("hit buffer size, ", mainBuf.Len())
				goto CLEAN_BUFFER
			}

		case <-idleTimeout.C:
			log.Tracef("%v no message input", idleDuration)
			goto CLEAN_BUFFER
		}

		goto READ_DOCS

	CLEAN_BUFFER:

		if docBuf.Len() > 0 {
			mainBuf.Write(docBuf.Bytes())
		}

		if mainBuf.Len() > 0 {
			client.Bulk(&mainBuf)
			//TODO handle retry and fallback/over, dead letter queue
			//set services to failure, need manual restart
			//process dead letter queue first next round

			stats.IncrementBy("bulk", "processed", int64(mainBuf.Len()))
			log.Trace("clean buffer, and execute bulk insert")
		}

	}
}