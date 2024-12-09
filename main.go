// Copyright (C) INFINI Labs & INFINI LIMITED.
//
// The INFINI Framework is offered under the GNU Affero General Public License v3.0
// and as commercial software.
//
// For commercial licensing, contact us at:
//   - Website: infinilabs.com
//   - Email: hello@infini.ltd
//
// Open Source licensed under AGPL V3:
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

/* Copyright © INFINI Ltd. All rights reserved.
 * web: https://infinilabs.com
 * mail: hello#infini.ltd */

package main

import (
	_ "expvar"
	"github.com/rubyniu105/framework"
	"github.com/rubyniu105/framework/core/module"
	"github.com/rubyniu105/framework/core/util"
	"github.com/rubyniu105/framework/modules/api"
	"github.com/rubyniu105/framework/modules/elastic"
	"github.com/rubyniu105/framework/modules/metrics"
	"github.com/rubyniu105/framework/modules/pipeline"
	"github.com/rubyniu105/framework/modules/queue"
	queue2 "github.com/rubyniu105/framework/modules/queue/disk_queue"
	"github.com/rubyniu105/framework/modules/redis"
	"github.com/rubyniu105/framework/modules/s3"
	stats2 "github.com/rubyniu105/framework/modules/stats"
	"github.com/rubyniu105/framework/modules/task"
	_ "github.com/rubyniu105/framework/plugins"
	stats "github.com/rubyniu105/framework/plugins/stats_statsd"
	"github.com/rubyniu105/gateway/config"
	_ "github.com/rubyniu105/gateway/pipeline"
	"github.com/rubyniu105/gateway/proxy"
	"github.com/rubyniu105/gateway/service/floating_ip"
	"github.com/rubyniu105/gateway/service/forcemerge"
)

func setup() {
	module.RegisterSystemModule(&stats2.SimpleStatsModule{})
	module.RegisterUserPlugin(&stats.StatsDModule{})
	module.RegisterSystemModule(&s3.S3Module{})
	module.RegisterSystemModule(&queue2.DiskQueue{})
	module.RegisterSystemModule(&redis.RedisModule{})
	module.RegisterSystemModule(&elastic.ElasticModule{})
	module.RegisterSystemModule(&queue.Module{})
	module.RegisterSystemModule(&task.TaskModule{})
	module.RegisterSystemModule(&api.APIModule{})
	module.RegisterModuleWithPriority(&pipeline.PipeModule{}, 100)

	module.RegisterUserPlugin(forcemerge.ForceMergeModule{})
	module.RegisterUserPlugin(floating_ip.FloatingIPPlugin{})
	module.RegisterUserPlugin(&metrics.MetricsModule{})
	module.RegisterPluginWithPriority(&proxy.GatewayModule{}, 200)
}

func start() {
	module.Start()
}

func main() {

	terminalHeader := ("\n   ___   _   _____  __  __    __  _       \n")
	terminalHeader += ("  / _ \\ /_\\ /__   \\/__\\/ / /\\ \\ \\/_\\ /\\_/\\\n")
	terminalHeader += (" / /_\\///_\\\\  / /\\/_\\  \\ \\/  \\/ //_\\\\\\_ _/\n")
	terminalHeader += ("/ /_\\\\/  _  \\/ / //__   \\  /\\  /  _  \\/ \\ \n")
	terminalHeader += ("\\____/\\_/ \\_/\\/  \\__/    \\/  \\/\\_/ \\_/\\_/ \n\n")

	terminalFooter := ""

	app := framework.NewApp("gateway", "A light-weight, powerful and high-performance search gateway.",
		util.TrimSpaces(config.Version), util.TrimSpaces(config.BuildNumber), util.TrimSpaces(config.LastCommitLog), util.TrimSpaces(config.BuildDate), util.TrimSpaces(config.EOLDate), terminalHeader, terminalFooter)

	app.Init(nil)

	defer app.Shutdown()

	if app.Setup(setup, start, nil) {
		app.Run()
	}

}
