package common

import (
	"testing"

	"github.com/bitshifted/easycloud/log"
	"github.com/stretchr/testify/suite"
)

type UtilsTestSuite struct {
	suite.Suite
}

func (ts *UtilsTestSuite) SetupSuite() {
	// Initialize logger
	log.Init()
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
	valueType := ValueTypeFromString("$FOO")
	ts.Equal(EnvVariableString, valueType)
	valueType = ValueTypeFromString("${BAR}")
	ts.Equal(EnvVariableString, valueType)
	valueType = ValueTypeFromString("$my_var123")
	ts.Equal(EnvVariableString, valueType)
	valueType = ValueTypeFromString("${INVALID_VAR")
	ts.Equal(PlainString, valueType)
}

func (ts *UtilsTestSuite) TestShouldReturnFileContentType() {
	valueType := ValueTypeFromString("filecontent:/path/to/file")
	ts.Equal(FileContentString, valueType)
	valueType = ValueTypeFromString("content:blah")
	ts.Equal(PlainString, valueType)
}

func (ts *UtilsTestSuite) TestShouldReturnCorrectEnvVarName() {
	name, err := ExtractEnvVarName("${FOO_VAR}")
	ts.NoError(err)
	ts.Equal("FOO_VAR", name)
	name, err = ExtractEnvVarName("$BAR")
	ts.NoError(err)
	ts.Equal("BAR", name)
	_, err = ExtractEnvVarName("${INVALID")
	ts.Error(err)
}
