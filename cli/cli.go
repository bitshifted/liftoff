package cli

import (
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
	TerraformPath   string      `help:"Path to Terraform binary"`
	PlaybookBinPath string      `help:"Path to ansible-playbook binary"`
	ConfigFile      string      `help:"Path to configuration file"`
	EnableDebug     bool        `help:"Enable debug logging"`
	Setup           SetupCmd    `cmd:"" help:"Setup and configure infrastructure"`
	TearDown        TearDownCmd `cmd:"" name:"teardown" help:"Cleanup created infrastructure"`
}

type SetupCmd struct {
	SkipTerraform bool `help:"Do not run Terraform"`
	SkipAnsible   bool `help:"Do not run Ansible"`
}

type TearDownCmd struct {
}

func (s *SetupCmd) Run(ctx *kong.Context) error {
	log.Logger.Info().Msg("Executing setup...")
	configFile := extractArgumentValue(ctx.Args, configFileArg, 1, common.DefaultConfigFileName)
	log.Logger.Info().Msgf("Reading configuration file %s", configFile)
	conf, err := config.LoadConfig(configFile)
	if err != nil {
		return err
	}
	executionConfig := exec.ExecutionConfig{
		Config:              conf,
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

func extractArgumentValue(args []string, argument string, valueIndex int8, defaultValue string) string {
	value := defaultValue
	for i, s := range args {
		if s == argument {
			value = args[i+int(valueIndex)]
			break
		}
	}
	return value
}
