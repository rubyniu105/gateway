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

type RequestAPIKeyFilter struct {
	genericFilter *RequestFilter
	Include       []string `config:"include"`
	Exclude       []string `config:"exclude"`
}

func (filter *RequestAPIKeyFilter) Name() string {
	return "request_api_key_filter"
}

func init() {
	pipeline.RegisterFilterPluginWithConfigMetadata("request_api_key_filter", NewRequestAPIKeyFilter, &RequestAPIKeyFilter{})
}

func NewRequestAPIKeyFilter(c *config.Config) (pipeline.Filter, error) {

	runner := RequestAPIKeyFilter{}
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

func (filter *RequestAPIKeyFilter) Filter(ctx *fasthttp.RequestCtx) {
	exists, apiID, _ := ctx.ParseAPIKey()
	if !exists {
		if global.Env().IsDebug {
			log.Tracef("API not exist")
		}
		return
	}

	apiIDStr := string(apiID)
	valid, hasRule := CheckExcludeStringRules(apiIDStr, filter.Exclude, ctx)
	if hasRule && !valid {
		filter.genericFilter.Filter(ctx)
		if global.Env().IsDebug {
			log.Debugf("must_not rules matched, this request has been filtered: %v", ctx.Request.PhantomURI().String())
		}
		return
	}

	valid, hasRule = CheckIncludeStringRules(apiIDStr, filter.Include, ctx)
	if hasRule && !valid {
		filter.genericFilter.Filter(ctx)
		if global.Env().IsDebug {
			log.Debugf("must_not rules matched, this request has been filtered: %v", ctx.Request.PhantomURI().String())
		}
		return
	}

}
