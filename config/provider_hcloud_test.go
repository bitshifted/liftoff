// Copyright 2024 Bitshift D.O.O
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"testing"

	"github.com/bitshifted/liftoff/log"
	"github.com/stretchr/testify/suite"
)

type ProviderHcloudTestSuite struct {
	suite.Suite
}

func (ts *ProviderHcloudTestSuite) SetupSuite() {
	// Initialize logger
	log.Init(true)
	log.Logger.Info().Msg("Running UtilsTestSuite")
}

func TestProviderHcloudTestSuite(t *testing.T) {
	suite.Run(t, new(ProviderHcloudTestSuite))
}

func (ts *ProviderHcloudTestSuite) TestPostLoadPlainTpken() {
	hcloud := ProviderHcloud{
		Token:    "some token",
		Endpoint: "http://localhost:8080/endpoint",
	}
	hcloud.postLoad()
	ts.Equal("some token", hcloud.Token)
	ts.Equal("http://localhost:8080/endpoint", hcloud.Endpoint)
	ts.Equal("", hcloud.PollInterval)
	ts.Equal("", hcloud.PollInterval)
}

func (ts *ProviderHcloudTestSuite) TestShouldReturnErrorWhenNoToken() {
	hcloud := ProviderHcloud{
		Endpoint: "http://localhost",
	}
	err := hcloud.postLoad()
	ts.Error(err)
	ts.Equal("token is required for hcloud provider", err.Error())
}

func (ts *ProviderHcloudTestSuite) TestShouldThrowErrorWhenEnvVarNotSet() {
	hcloud := ProviderHcloud{
		Token: "$FOO",
	}
	err := hcloud.postLoad()
	ts.Error(err)
	ts.Equal("referenced environment variable FOO is not set", err.Error())
}

func (ts *ProviderHcloudTestSuite) TestShouldPassWhenEnvVarIsSet() {
	hcloud := ProviderHcloud{
		Token: "$FOO",
	}
	ts.Suite.T().Setenv("FOO", "value")
	err := hcloud.postLoad()
	ts.NoError(err)
	ts.Equal("", hcloud.Token)
}
