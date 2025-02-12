// Copyright 2025 Bitshift D.O.O
// SPDX-License-Identifier: MPL-2.0

package exec

import (
	"testing"

	"github.com/bitshifted/liftoff/config"
	"github.com/bitshifted/liftoff/log"
	"github.com/stretchr/testify/suite"
)

type ExecutionSetupTestSuite struct {
	suite.Suite
}

func (ts *ExecutionSetupTestSuite) SetupSuite() {
	// Initialize logger
	log.Init(true)
	log.Logger.Info().Msg("Running ExecutionSetupTestSuite")
}

func TestExecutionSetuptestSuite(t *testing.T) {
	suite.Run(t, new(ExecutionSetupTestSuite))
}

func (ts *ExecutionSetupTestSuite) TestSshConfigFile_NoBastion() {
	ec := &ExecutionConfig{
		Config: &config.Configuration{
			ProcessingVars: map[string]interface{}{
				"ansible_user":            "testuser",
				"ansible_ssh_private_key": "/path/to/private/key",
			},
		},
		ConfigFilePath: "/path/to/config.yaml",
	}
	err := ec.generateSSHConfig()
	ts.NoError(err)
	ts.Equal("/tmp/ssh_config_2f706174", ec.Config.ProcessingVars["ssh_config_file"])
}

func (ts *ExecutionSetupTestSuite) TestSshConfigFile_WithBastion() {
	ec := &ExecutionConfig{
		Config: &config.Configuration{
			ProcessingVars: map[string]interface{}{
				"ansible_user":            "testuser",
				"ansible_ssh_private_key": "/path/to/private/key",
				"bastion_address":         "bastion.example.com",
			},
		},
		ConfigFilePath: "/bastion/config.yaml",
	}
	err := ec.generateSSHConfig()
	ts.NoError(err)
	ts.Equal("/tmp/ssh_config_2f626173", ec.Config.ProcessingVars["ssh_config_file"])
}
