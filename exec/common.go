// Copyright 2024 Bitshift D.O.O
// SPDX-License-Identifier: MPL-2.0

package exec

import (
	"os"
	"os/exec"

	"github.com/bitshifted/liftoff/config"
	"github.com/bitshifted/liftoff/log"
)

const (
	defaltTerraformCmd = "terraform"
	defaultAnsibleCmd  = "ansible-playbook"
)

type ExecutionConfig struct {
	Config              *config.Configuration
	SkipTerraform       bool
	SkipAnsible         bool
	TerraformPath       string
	AnsiblePlaybookPath string
	TerraformWorkDir    string
	AnsibleWorkDir      string
}

func (ec *ExecutionConfig) executeTerraformCommand(cmd ...string) error {
	command := exec.Command(ec.TerraformPath, cmd...) //nolint:gosec
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Dir = ec.TerraformWorkDir
	log.Logger.Debug().Msgf("Terraform work directory: %s", command.Dir)
	command.Env = append(command.Env, os.Environ()...)
	command.Env = append(command.Env, "TF_DATA_DIR="+ec.Config.Terraform.DataDir)
	log.Logger.Debug().Msgf("Command environment: %v", command.Env)
	return command.Run()
}
