// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package cast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToInt(t *testing.T) {
	var eight interface{} = 8
	assert.Equal(t, ToInt(8), 8)
	assert.Equal(t, ToInt(8.31), 8)
	assert.Equal(t, ToInt("8"), 8)
	assert.Equal(t, ToInt(true), 1)
	assert.Equal(t, ToInt(false), 0)
	assert.Equal(t, ToInt(eight), 8)
}

func TestToFloat64(t *testing.T) {
	var eight interface{} = 8
	assert.Equal(t, ToFloat64(8), 8.00)
	assert.Equal(t, ToFloat64(8.31), 8.31)
	assert.Equal(t, ToFloat64("8.31"), 8.31)
	assert.Equal(t, ToFloat64(eight), 8.0)
}

func TestToString(t *testing.T) {
	var foo interface{} = "one more time"
	assert.Equal(t, ToString(8), "8")
	assert.Equal(t, ToString(8.12), "8.12")
	assert.Equal(t, ToString([]byte("one time")), "one time")
	assert.Equal(t, ToString(foo), "one more time")
	assert.Equal(t, ToString(nil), "")
}

func TestMaps(t *testing.T) {
	var taxonomies = map[interface{}]interface{}{"tag": "tags", "group": "groups"}
	assert.Equal(t, ToStringMap(taxonomies), map[string]interface{}{"tag": "tags", "group": "groups"})
}
