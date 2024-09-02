// Copyright 2024 Bitshift D.O.O
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessVariables(t *testing.T) {
	t.Setenv("FOO", "foo")
	input := map[string]interface{}{
		"num":  1,
		"bool": true,
		"nested": map[string]interface{}{
			"plain":  "plain string",
			"envval": "fromenv:FOO",
			"nested1": map[string]interface{}{
				"filecontent": "fromfile:test_files/testfile.txt",
			},
			"list": []interface{}{
				"plain value",
				"fromfile:test_files/testfile.txt",
			},
		},
	}
	err := processVariables(input)
	assert.NoError(t, err)
	assert.Equal(t, 1, input["num"])
	assert.True(t, input["bool"].(bool))
	nested := input["nested"].(map[string]interface{})
	assert.Equal(t, "plain string", nested["plain"])
	assert.Equal(t, "foo", nested["envval"])
	nested1 := nested["nested1"].(map[string]interface{})
	assert.Equal(t, "test content", nested1["filecontent"])
	list := nested["list"].([]interface{})
	assert.Equal(t, "plain value", list[0])
	assert.Equal(t, "test content", list[1])
	assert.NotNil(t, list)
	assert.Equal(t, 2, len(list))
	t.Setenv("FOO", "")
}
