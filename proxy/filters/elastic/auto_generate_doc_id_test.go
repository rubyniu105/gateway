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

/* Copyright © INFINI LTD. All rights reserved.
 * Web: https://infinilabs.com
 * Email: hello#infini.ltd */

package elastic

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAutoGenerateDocID(t *testing.T) {
	//PUT /<target>/_doc/<_id>
	//POST /<target>/_doc/
	//POST /<target>/_create/<_id>
	//PUT twitter/tweet/1/_create
	//POST twitter/tweet/

	path:="/index/doc/"
	valid,urlLevelIndex, urlLevelType,urlLevelID := ParseURLMeta(path)
	assert.Equal(t,true,valid)
	assert.Equal(t,"index",urlLevelIndex)
	assert.Equal(t,"doc",urlLevelType)
	assert.Equal(t,"",urlLevelID)


	path="/index/doc"
	valid,urlLevelIndex, urlLevelType,urlLevelID = ParseURLMeta(path)
	assert.Equal(t,true,valid)
	assert.Equal(t,"index",urlLevelIndex)
	assert.Equal(t,"doc",urlLevelType)
	assert.Equal(t,"",urlLevelID)

	path="/index/_doc"
	valid,urlLevelIndex, urlLevelType,urlLevelID = ParseURLMeta(path)
	assert.Equal(t,true,valid)
	assert.Equal(t,"index",urlLevelIndex)
	assert.Equal(t,"_doc",urlLevelType)
	assert.Equal(t,"",urlLevelID)

	path="/index/_doc/"
	valid,urlLevelIndex, urlLevelType,urlLevelID = ParseURLMeta(path)
	assert.Equal(t,true,valid)
	assert.Equal(t,"index",urlLevelIndex)
	assert.Equal(t,"_doc",urlLevelType)
	assert.Equal(t,"",urlLevelID)

	path="/index/doc/1"
	valid,urlLevelIndex, urlLevelType,urlLevelID = ParseURLMeta(path)
	assert.Equal(t,true,valid)
	assert.Equal(t,"index",urlLevelIndex)
	assert.Equal(t,"doc",urlLevelType)
	assert.Equal(t,"1",urlLevelID)

	path="/index/_doc/1"
	valid,urlLevelIndex, urlLevelType,urlLevelID = ParseURLMeta(path)
	assert.Equal(t,true,valid)
	assert.Equal(t,"index",urlLevelIndex)
	assert.Equal(t,"_doc",urlLevelType)
	assert.Equal(t,"1",urlLevelID)

	path="/index/_create/1"
	valid,urlLevelIndex, urlLevelType,urlLevelID = ParseURLMeta(path)
	assert.Equal(t,true,valid)
	assert.Equal(t,"index",urlLevelIndex)
	assert.Equal(t,"_create",urlLevelType)
	assert.Equal(t,"1",urlLevelID)

	path="/index/_doc/1"
	valid,urlLevelIndex, urlLevelType,urlLevelID = ParseURLMeta(path)
	assert.Equal(t,true,valid)
	assert.Equal(t,"index",urlLevelIndex)
	assert.Equal(t,"_doc",urlLevelType)
	assert.Equal(t,"1",urlLevelID)

	path="/index/_doc/_bulk"
	valid,urlLevelIndex, urlLevelType,urlLevelID = ParseURLMeta(path)
	assert.Equal(t,false,valid)

	path="/index/_doc/_search"
	valid,urlLevelIndex, urlLevelType,urlLevelID = ParseURLMeta(path)
	assert.Equal(t,false,valid)



}
