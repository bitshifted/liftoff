// Copyright 2024 Bitshift D.O.O
// SPDX-License-Identifier: MPL-2.0

package exec

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	osExec "os/exec"
	"path"

	"github.com/bitshifted/liftoff/common"
	"github.com/bitshifted/liftoff/log"
	"github.com/bitshifted/liftoff/template"
)

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
	output, err := ec.calculateOutputDirectory()
	if err != nil {
		return err
	}
	ec.TerraformWorkDir = path.Join(output, common.DefaultTerraformDir)
	ec.AnsibleWorkDir = path.Join(output, common.DefaultAnsibleDir)
	processor := template.TemplateProcessor{
		BaseDir:   tmplDir,
		OutputDir: output,
	}
	log.Logger.Info().Msg("Processing templates...")
	err = processor.ProcessTemplates(ec.Config)
	if err != nil {
		return err
	}
	var tfOutput map[string]interface{}
	if ec.TerraformPath == "" {
		tfCmdPath, e := osExec.LookPath(defaltTerraformCmd)
		if e != nil {
			log.Logger.Error().Err(e).Msg("Failed to lookup Terraform path")
			return e
		}
		ec.TerraformPath = tfCmdPath
	}
	log.Logger.Debug().Msgf("Using Terraform command: %s", ec.TerraformPath)
	if !ec.SkipTerraform {
		err = ec.executeTerraform()
		if err != nil {
			return err
		}
	} else {
		log.Logger.Info().Msg("Skipping Terraform configuration")
	}
	tfOutput, err = ec.getTerraformOutputs()
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to get Terraform outputs")
		return err
	}

	if !ec.SkipAnsible {
		log.Logger.Info().Msg("Processing Ansible configuration...")
		log.Logger.Debug().Msgf("Terraform output: %v", tfOutput)
		return ec.executeAnsiblePlaybook()
	} else {
		log.Logger.Info().Msg("Skipping Ansible configuration")
	}

	return nil
}

func (ec *ExecutionConfig) executeTerraform() error {
	// TODO copy .terraform.hcl.lock if exists in workspace
	log.Logger.Info().Msg("Running Terraform init...")
	err := ec.executeTerraformCommand("init")
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to run Terraform init")
		return err
	}
	log.Logger.Info().Msg("Running Terraform apply")
	err = ec.executeTerraformCommand("apply", "-auto-approve")
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to run Terraform apply")
	}
	return err
}

func (ec *ExecutionConfig) getTerraformOutputs() (map[string]interface{}, error) {
	// // run terraform output
	log.Logger.Info().Msg("Collecting Terraform outputs")
	cmdOut := osExec.Command(ec.TerraformPath, "output", "-json") //nolint:gosec
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
	var outputs map[string]interface{}
	if buf.Len() == 0 {
		outputs = make(map[string]interface{})
	} else {
		_, err = io.Copy(&buf, r)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(buf.Bytes(), &outputs)
		if err != nil {
			log.Logger.Error().Err(err).Msg("Failed to unmarshal Terraform output")
			return nil, err
		}
	}

	return outputs, nil
}

func (ec *ExecutionConfig) executeAnsiblePlaybook() error {
	if ec.AnsiblePlaybookPath == "" {
		ansibleCmdPath, err := osExec.LookPath(defaultAnsibleCmd)
		if err != nil {
			log.Logger.Error().Err(err).Msg("Failed to lookup ansible-playbook path")
			return err
		}
		ec.TerraformPath = ansibleCmdPath
	}
	if ec.Config.Ansible == nil || ec.Config.Ansible.InventoryFile == "" || ec.Config.Ansible.PlaybookFile == "" {
		log.Logger.Warn().Msg("Either Ansible inventory file or playbook were not specified. Aborting.")
		return nil
	}
	log.Logger.Debug().Msgf("Using ansible-playbook command: %s", ec.AnsiblePlaybookPath)
	log.Logger.Info().Msg("Running ansible-playbook")
	cmdApply := osExec.Command(ec.AnsiblePlaybookPath, "-i", ec.Config.Ansible.InventoryFile, ec.Config.Ansible.PlaybookFile) //nolint:gosec
	cmdApply.Stdout = os.Stdout
	cmdApply.Stderr = os.Stderr
	cmdApply.Dir = ec.AnsibleWorkDir
	err := cmdApply.Run()
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to run ansible-playbook")
	}
	return err
}
