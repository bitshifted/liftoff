// Copyright 2024 Bitshift D.O.O
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"fmt"
	"os"
	"testing"

	"github.com/bitshifted/liftoff/log"
	"github.com/stretchr/testify/suite"
)

type UtilsTestSuite struct {
	suite.Suite
}

func (ts *UtilsTestSuite) SetupSuite() {
	// Initialize logger
	log.Init(true)
	log.Logger.Info().Msg("Running UtilsTestSuite")
}

func TestUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(UtilsTestSuite))
}

func (ts *UtilsTestSuite) TestEnvVarIsSetReturnsCorrectValue() {
	osGetEnv = func(key string) string {
		if key == "FOO" {
			return "bar"
		}
		return ""
	}
	ts.True(IsEnvVariableSet("FOO"))
	ts.False(IsEnvVariableSet("BAR"))
}

func (ts *UtilsTestSuite) TestShouldReturnEnvVarTypeForString() {
	valueType := ValueTypeFromString("fromenv:FOO")
	ts.Equal(EnvVariableString, valueType)
	valueType = ValueTypeFromString("fromenv:BAR")
	ts.Equal(EnvVariableString, valueType)
	valueType = ValueTypeFromString("fromenv:my_var123")
	ts.Equal(EnvVariableString, valueType)
}

func (ts *UtilsTestSuite) TestShouldReturnFileContentType() {
	valueType := ValueTypeFromString("fromfile:/path/to/file")
	ts.Equal(FileContentString, valueType)
	valueType = ValueTypeFromString("content:blah")
	ts.Equal(PlainString, valueType)
}

func (ts *UtilsTestSuite) TestShouldReturnCorrectEnvVarName() {
	name := ExtractEnvVarName("fromenv:FOO_VAR")
	ts.Equal("FOO_VAR", name)
	name = ExtractEnvVarName("fromenv:BAR")
	ts.Equal("BAR", name)
}

func (ts *UtilsTestSuite) TestShouldReturnFileContentRelPath() {
	content, err := ProcessStringValue("fromfile:test_files/sample.txt")
	ts.NoError(err)
	ts.Equal("sample text", content)
}

func (ts *UtilsTestSuite) TestShouldReturnFileContentAbsPath() {
	file, err := os.CreateTemp("", "content")
	log.Logger.Debug().Msgf("temp file path: %s", file.Name())
	ts.NoError(err)
	file.WriteString("sample text")
	content, err := ProcessStringValue(fmt.Sprintf("fromfile:%s", file.Name()))
	ts.NoError(err)
	ts.Equal("sample text", content)
}
