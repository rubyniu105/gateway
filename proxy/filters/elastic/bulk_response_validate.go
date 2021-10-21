/* ©INFINI, All Rights Reserved.
 * mail: contact#infini.ltd */

package elastic

import (
	"fmt"
	log "github.com/cihub/seelog"
	"infini.sh/framework/core/config"
	"infini.sh/framework/core/elastic"
	"infini.sh/framework/core/global"
	"infini.sh/framework/core/pipeline"
	"infini.sh/framework/core/queue"
	"infini.sh/framework/core/stats"
	"infini.sh/framework/core/util"
	"infini.sh/framework/lib/bytebufferpool"
	"infini.sh/framework/lib/fasthttp"
	"net/http"
)

type BulkResponseValidate struct {
	config *Config
}

func (this *BulkResponseValidate) Name() string {
	return "bulk_response_validate"
}

func (this *BulkResponseValidate) Filter(ctx *fasthttp.RequestCtx) {
	path := string(ctx.URI().Path())
	if string(ctx.Request.Header.Method()) != "POST" || !util.ContainStr(path, "_bulk") {
		return
	}

	if ctx.Response.StatusCode() == http.StatusOK || ctx.Response.StatusCode() == http.StatusCreated {
		var resbody = ctx.Response.GetRawBody()
		containError := util.LimitedBytesSearch(resbody, []byte("\"errors\":true"), 64)
		if containError {
			if global.Env().IsDebug {
				log.Error("error in bulk requests,", util.SubString(string(resbody), 0, 256))
			}

			//decode response
			response := elastic.BulkResponse{}
			err := response.UnmarshalJSON(resbody)
			if err != nil {
				panic(err)
			}
			var contains400Error = false
			//busyRejectOffset := map[int]elastic.BulkActionMetadata{}
			invalidOffset := map[int]elastic.BulkActionMetadata{}
			//failureOffset := map[int]elastic.BulkActionMetadata{}
			var invalidCount = 0
			var statsCodeStats = map[int]int{}
			for i, v := range response.Items {
				item := v.GetItem()

				x, ok := statsCodeStats[item.Status]
				if !ok {
					x = 0
				}
				x++
				statsCodeStats[item.Status] = x

				if item.Error != nil {
					invalidCount++
					invalidOffset[i] = v
				}
			}

			for x, y := range statsCodeStats {
				stats.IncrementBy("bulk_items", fmt.Sprintf("%v", x), int64(y))
			}

			if invalidCount > 0 {

				requestBytes := ctx.Request.GetRawBody()
				nonRetryableItems := bytebufferpool.Get()
				retryableItems := bytebufferpool.Get()

				var offset = 0
				var match = false
				var retryable = false
				var response elastic.BulkActionMetadata
				invalidCount = 0
				var failureCount = 0
				//walk bulk message, with invalid id, save to another list

				var docBuffer []byte
				docBuffer = p.Get(this.config.DocBufferSize)
				defer p.Put(docBuffer)

				WalkBulkRequests(requestBytes, docBuffer, func(eachLine []byte) (skipNextLine bool) {
					return false
				}, func(metaBytes []byte, actionStr, index, typeName, id string) (err error) {
					response, match = invalidOffset[offset]
					if match {

						//find invalid request
						if response.GetItem().Status >= 400 && response.GetItem().Status < 500 && response.GetItem().Status != 429 {
							retryable = false
							contains400Error = true
							if nonRetryableItems.Len() > 0 {
								nonRetryableItems.WriteByte('\n')
							}
							nonRetryableItems.Write(metaBytes)
							invalidCount++
						} else {
							retryable = true
							if retryableItems.Len() > 0 {
								retryableItems.WriteByte('\n')
							}
							retryableItems.Write(metaBytes)
							failureCount++
						}
					}
					offset++
					return nil
				}, func(payloadBytes []byte) {
					if match {
						if payloadBytes != nil && len(payloadBytes) > 0 {
							if retryable {
								if retryableItems.Len() > 0 {
									retryableItems.WriteByte('\n')
								}
								retryableItems.Write(payloadBytes)
							} else {
								if nonRetryableItems.Len() > 0 {
									nonRetryableItems.WriteByte('\n')
								}
								nonRetryableItems.Write(payloadBytes)
							}
						}
					}
				})

				if nonRetryableItems.Len() > 0 {
					nonRetryableItems.WriteByte('\n')
					bytes := ctx.Request.OverrideBodyEncode(nonRetryableItems.Bytes(),true)
					queue.Push(this.config.InvalidQueue, bytes)
					//send to redis channel
					nonRetryableItems.Reset()
					bytebufferpool.Put(nonRetryableItems)
				}

				if retryableItems.Len() > 0 {
					retryableItems.WriteByte('\n')
					bytes := ctx.Request.OverrideBodyEncode(retryableItems.Bytes(),true)
					queue.Push(this.config.FailureQueue, bytes)
					retryableItems.Reset()
					bytebufferpool.Put(retryableItems)
				}
			}

			if contains400Error {
				ctx.Response.SetStatusCode(this.config.InvalidStatus)
			} else {
				ctx.Response.SetStatusCode(this.config.FailureStatus)
			}

			if this.config.ContinueOnError {
				ctx.Finished()
			}

		}
	}
}

type Config struct {
	DocBufferSize   int    `config:"doc_buffer_size"`
	InvalidQueue    string `config:"invalid_queue"`
	FailureQueue    string `config:"failure_queue"`
	InvalidStatus   int    `config:"invalid_status"`
	FailureStatus   int    `config:"failure_status"`
	ContinueOnError bool   `config:"continue_on_error"`
}

func NewBulkResponseValidate(c *config.Config) (pipeline.Filter, error) {
	cfg := Config{
		DocBufferSize: 256 * 1024,
	}
	if err := c.Unpack(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unpack the filter configuration : %s", err)
	}

	runner := BulkResponseValidate{config: &cfg}

	return &runner, nil
}
