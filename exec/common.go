// Copyright 2024 Bitshift D.O.O
// SPDX-License-Identifier: MPL-2.0

package exec

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/bitshifted/liftoff/config"
	"github.com/bitshifted/liftoff/log"
)

const (
	defaltTerraformCmd = "terraform"
	defaultAnsibleCmd  = "ansible-playbook"
)

type ExecutionConfig struct {
	Config              *config.Configuration
	ConfigFilePath      string
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
	return command.Run()
}

// calculates outpur directory name based on configuration file name
func (ec *ExecutionConfig) calculateOutputDirectory() (string, error) {
	configFileName := filepath.Base(ec.ConfigFilePath)
	configFileExt := filepath.Ext(ec.ConfigFilePath)
	configDir := filepath.Dir(ec.ConfigFilePath)
	// strip extension
	genDirName := strings.Replace(configFileName, configFileExt, "", 1)
	log.Logger.Debug().Msgf("Directory for generated files: %s", genDirName)
	// create directory
	genDirPath := path.Join(configDir, genDirName)
	err := os.MkdirAll(genDirPath, os.ModePerm)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to create directory for generated files")
		return "", err
	}
	return genDirPath, nil
}
