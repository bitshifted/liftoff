// Copyright 2024 Bitshift D.O.O
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"fmt"
	"path/filepath"

	"github.com/alecthomas/kong"
	"github.com/bitshifted/liftoff/common"
	"github.com/bitshifted/liftoff/config"
	"github.com/bitshifted/liftoff/exec"
	"github.com/bitshifted/liftoff/log"
)

const (
	configFileArg          = "--config-file"
	terraformPathArg       = "--terraform-path"
	ansiblePlaybookPathArg = "--playbook-bin-path"
)

type CLI struct {
	TerraformPath   string          `help:"Path to Terraform binary"`
	PlaybookBinPath string          `help:"Path to ansible-playbook binary"`
	ConfigFile      string          `help:"Path to configuration file"`
	EnableDebug     bool            `help:"Enable debug logging"`
	Setup           SetupCmd        `cmd:"" help:"Setup and configure infrastructure"`
	TearDown        TearDownCmd     `cmd:"" name:"teardown" help:"Cleanup created infrastructure"`
	Version         VersionCmd      `cmd:"" name:"version" help:"Display version information"`
	TestTemplate    TestTemplateCmd `cmd:"" name:"test-template" help:"Generate code from template and perform sanity checks"`
}

type SetupCmd struct {
	SkipTerraform bool `help:"Do not run Terraform"`
	SkipAnsible   bool `help:"Do not run Ansible"`
}

type TearDownCmd struct {
}

type VersionCmd struct {
}

type TestTemplateCmd struct {
}

func (s *SetupCmd) Run(ctx *kong.Context) error {
	log.Logger.Info().Msg("Executing setup...")
	configFile := extractArgumentValue(ctx.Args, configFileArg, 1, common.DefaultConfigFileName)
	conf, err := config.LoadConfig(configFile)
	if err != nil {
		return err
	}
	configFileAbsPath, err := filepath.Abs(configFile)
	if err != nil {
		return err
	}
	log.Logger.Info().Msgf("Reading configuration file %s", configFileAbsPath)
	executionConfig := exec.ExecutionConfig{
		Config:              conf,
		ConfigFilePath:      configFileAbsPath,
		SkipTerraform:       s.SkipTerraform,
		SkipAnsible:         s.SkipAnsible,
		TerraformPath:       extractArgumentValue(ctx.Args, terraformPathArg, 1, ""),
		AnsiblePlaybookPath: extractArgumentValue(ctx.Args, ansiblePlaybookPathArg, 1, ""),
	}
	return executionConfig.ExecuteSetup()
}

func (t *TearDownCmd) Run(ctx *kong.Context) error {
	log.Logger.Info().Msg("Executing teardown...")
	return nil
}

func (vc *VersionCmd) Run(ctx *kong.Context) error {
	fmt.Printf("Version: %s\nBuild number: %s\nCommit ID: %s\n",
		ProgramVersion.Version, ProgramVersion.BuildNumber, ProgramVersion.CommitID)
	return nil
}

func (tc *TestTemplateCmd) Run(ctx *kong.Context) error {
	log.Logger.Info().Msg("Performing template test...")
	configFile := extractArgumentValue(ctx.Args, configFileArg, 1, common.DefaultConfigFileName)
	conf, err := config.LoadConfig(configFile)
	if err != nil {
		return err
	}
	configFileAbsPath, err := filepath.Abs(configFile)
	if err != nil {
		return err
	}
	log.Logger.Info().Msgf("Reading configuration file %s", configFileAbsPath)
	executionConfig := exec.ExecutionConfig{
		Config:         conf,
		ConfigFilePath: configFileAbsPath,
	}
	return executionConfig.ExecuteTestTemplate()
}

func extractArgumentValue(args []string, argument string, valueIndex int8, defaultValue string) string { //nolint:unparam
	value := defaultValue
	for i, s := range args {
		if s == argument {
			value = args[i+int(valueIndex)]
			break
		}
	}
	return value
}
