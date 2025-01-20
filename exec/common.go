// Copyright 2024 Bitshift D.O.O
// SPDX-License-Identifier: MPL-2.0

package exec

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/bitshifted/liftoff/config"
	"github.com/bitshifted/liftoff/gitops"
	"github.com/bitshifted/liftoff/log"
)

const (
	defaltTerraformCmd = "terraform"
	defaultAnsibleCmd  = "ansible-playbook"
	liftoffHomeDirName = ".liftoff"
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
	tfDataDir := ec.calculateTerraformDataDir()
	if tfDataDir != "" {
		command.Env = append(command.Env, fmt.Sprintf("TF_DATA_DIR=%s", tfDataDir))
	}
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

func (ec *ExecutionConfig) templateDirAbsPath() (string, error) {
	repo := ec.Config.TemplateRepo
	tmplDirAbsPath := ""
	if repo == "" {
		log.Logger.Info().Msg("Template repository not specified")
	} else {
		tmpDir, err := os.MkdirTemp("", "template_repo")
		if err != nil {
			log.Logger.Error().Err(err).Msg("Failed to create temprorary directory for clone")
			return "", err
		}
		tmplDirAbsPath = tmpDir
		log.Logger.Info().Msgf("Cloning template repository %s to %s", repo, tmplDirAbsPath)
		handler := gitops.GitHandler{
			URL:         repo,
			Version:     ec.Config.TempateVersion,
			Destination: tmpDir,
		}
		err = handler.Fetch()
		if err != nil {
			return "", err
		}
	}
	if ec.Config.TemplateDir != "" {
		tmplDirAbsPath = path.Join(tmplDirAbsPath, ec.Config.TemplateDir)
	}
	return tmplDirAbsPath, nil
}

func (ec *ExecutionConfig) calculateTerraformDataDir() string {
	configFileName := filepath.Base(ec.ConfigFilePath)
	configFileExt := filepath.Ext(ec.ConfigFilePath)
	// strip extension
	strippedFileName := strings.Replace(configFileName, configFileExt, "", 1)
	hash := sha256.New().Sum([]byte(ec.ConfigFilePath))
	resultFileName := fmt.Sprintf("%s-%s", strippedFileName, hex.EncodeToString(hash)[0:8])
	homeDirPath, err := os.UserHomeDir()
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to get user home directory")
		return ""
	}
	tfDataDirPath := path.Join(homeDirPath, liftoffHomeDirName, resultFileName)
	log.Logger.Debug().Msgf("Terraform data directory: %s", tfDataDirPath)
	return tfDataDirPath
}
