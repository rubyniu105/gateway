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

package filter

import (
	"fmt"

	log "github.com/cihub/seelog"
	"github.com/rubyniu105/framework/core/config"
	"github.com/rubyniu105/framework/core/global"
	"github.com/rubyniu105/framework/core/pipeline"
	"github.com/rubyniu105/framework/lib/fasthttp"
)

type ResponseStatusCodeFilter struct {
	genericFilter *RequestFilter
	Include       []int `config:"include"`
	Exclude       []int `config:"exclude"`
}

func (filter ResponseStatusCodeFilter) Name() string {
	return "response_status_filter"
}

func init() {
	pipeline.RegisterFilterPluginWithConfigMetadata("response_status_filter", NewResponseStatusCodeFilter, &ResponseStatusCodeFilter{})
}

func NewResponseStatusCodeFilter(c *config.Config) (pipeline.Filter, error) {

	runner := ResponseStatusCodeFilter{}
	if err := c.Unpack(&runner); err != nil {
		return nil, fmt.Errorf("failed to unpack the filter configuration : %s", err)
	}

	runner.genericFilter = &RequestFilter{
		Action: "deny",
		Status: 403,
	}

	if err := c.Unpack(runner.genericFilter); err != nil {
		return nil, fmt.Errorf("failed to unpack the filter configuration : %s", err)
	}

	return &runner, nil
}

func (filter *ResponseStatusCodeFilter) Filter(ctx *fasthttp.RequestCtx) {

	code := ctx.Response.StatusCode()

	if global.Env().IsDebug {
		log.Debug("code:", code, ",exclude:", filter.Exclude)
	}
	if len(filter.Exclude) > 0 {
		for _, x := range filter.Exclude {
			y := int(x)
			if global.Env().IsDebug {
				log.Debugf("exclude code: %v vs %v, match: %v", x, code, y == code)
			}
			if y == code {
				filter.genericFilter.Filter(ctx)
				if global.Env().IsDebug {
					log.Debugf("rule matched, this request has been filtered: %v", ctx.Request.PhantomURI().String())
				}
				return
			}
		}
	}

	if global.Env().IsDebug {
		log.Debug("include:", filter.Include)
	}
	if len(filter.Include) > 0 {
		for _, x := range filter.Include {
			y := int(x)
			if global.Env().IsDebug {
				log.Debugf("include code: %v vs %v, match: %v", x, code, y == code)
			}
			if y == code {
				if global.Env().IsDebug {
					log.Debugf("rule matched, this request has been marked as good one: %v", ctx.Request.PhantomURI().String())
				}
				return
			}
		}
		filter.genericFilter.Filter(ctx)
		if global.Env().IsDebug {
			log.Debugf("no rule matched, this request has been filtered: %v", ctx.Request.PhantomURI().String())
		}
	}
}
