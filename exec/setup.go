package exec

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	osExec "os/exec"
	"path"

	"github.com/bitshifted/easycloud/common"
	"github.com/bitshifted/easycloud/config"
	"github.com/bitshifted/easycloud/log"
	"github.com/bitshifted/easycloud/template"
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
}

func (ec *ExecutionConfig) ExecuteSetup() error {
	repo := ec.Config.TemplateRepo
	if repo == "" {
		log.Logger.Info().Msg("Template repository not specified")
	}
	tmplDir := ec.Config.TemplateDir
	if tmplDir == "" {
		return errors.New("either template repository or template directory must be specified")
	}
	log.Logger.Info().Msgf("Template directory: %s", tmplDir)
	output, err := os.MkdirTemp("", "easycloud_setup")
	if err != nil {
		return err
	}
	ec.TerraformWorkDir = path.Join(output, common.DefaultTerraformDir)
	var tfOutput map[string]interface{}
	if !ec.SkipTerraform {
		processor := template.TemplateProcessor{
			BaseDir:   tmplDir,
			OutputDir: output,
		}
		log.Logger.Info().Msg("Processing Terraform configuration...")
		processor.ProcessTerraformTemplate(ec.Config)
		tfOutput, err = ec.executeTerraform()
		if err != nil {
			return err
		}
	} else {
		log.Logger.Info().Msg("Skipping Terraform configuration")
	}

	if !ec.SkipAnsible {
		log.Logger.Info().Msg("Processing Ansible configuration...")
		log.Logger.Debug().Msgf("Terraform output: %v", tfOutput)
	} else {
		log.Logger.Info().Msg("Skipping Ansible configuration")
	}

	return nil
}

func (ec *ExecutionConfig) executeTerraform() (map[string]interface{}, error) {
	if ec.TerraformPath == "" {
		path, err := osExec.LookPath(defaltTerraformCmd)
		if err != nil {
			log.Logger.Error().Err(err).Msg("Failed to lookup Terraform path")
			return nil, err
		}
		ec.TerraformPath = path
	}
	log.Logger.Debug().Msgf("Using Terraform command: %s", ec.TerraformPath)

	cmdInit := osExec.Command(ec.TerraformPath, "init")
	cmdInit.Stdout = os.Stdout
	cmdInit.Stderr = os.Stderr
	cmdInit.Dir = ec.TerraformWorkDir
	log.Logger.Debug().Msgf("Terraform work directory: %s", cmdInit.Dir)
	log.Logger.Info().Msg("Running Terraform init...")
	err := cmdInit.Run()
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to run Terraform init")
		return nil, err
	}
	log.Logger.Info().Msg("Running Terraform apply")
	cmdApply := osExec.Command(ec.TerraformPath, "apply", "-auto-approve")
	cmdApply.Stdout = os.Stdout
	cmdApply.Stderr = os.Stderr
	cmdApply.Dir = ec.TerraformWorkDir
	err = cmdApply.Run()
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to run Terraform apply")
		return nil, err
	}
	// run terraform output
	cmdOut := osExec.Command(ec.TerraformPath, "outpit", "-json")
	cmdOut.Dir = ec.TerraformWorkDir
	r, w, err := os.Pipe()
	defer w.Close()
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to create pipe")
		return nil, err
	}
	cmdOut.Stdout = r
	err = cmdOut.Run()
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to run Terraform output")
		return nil, nil
	}
	var buf bytes.Buffer
	io.Copy(&buf, r)
	var outputs map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &outputs)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to unmarshal Terraform output")
		return nil, err
	}
	return outputs, nil
}

func (ec *ExecutionConfig) executeAnsiblePlaybook() error {
	return nil
}
