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
	"github.com/bitshifted/liftoff/config"
	"github.com/bitshifted/liftoff/log"
	"github.com/bitshifted/liftoff/template"
)

func (ec *ExecutionConfig) ExecuteSetup() error {
	tmplDir, err := ec.templateDirAbsPath()
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to get template directory path")
		return err
	}
	if tmplDir == "" {
		return errors.New("either template repository or template directory must be specified")
	}
	log.Logger.Info().Msgf("Template directory: %s", tmplDir)
	tmplConfig, err := config.LoadTemplateConfig(tmplDir)
	if err != nil {
		log.Logger.Error().Err(err).Msgf("Failed to load template configuration: %s", err.Error())
		return err
	}
	ec.Config.TemplateConfig = tmplConfig
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
	log.Logger.Info().Msg("Processing Terraform templates...")
	err = processor.ProcessTerraformTemplates(ec.Config)
	if err != nil {
		return err
	}
	var tfOutputs map[string]interface{}
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
	tfOutputs, err = ec.getTerraformOutputs()
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to get Terraform outputs")
		return err
	}
	// add TF  outputs to variables
	for k, v := range tfOutputs {
		// extract values of TF output variables
		ec.Config.ProcessingVars[k] = v.(map[string]interface{})["value"]
	}
	log.Logger.Debug().Msgf("Terraform output: %v", tfOutputs)

	if !ec.SkipAnsible {
		log.Logger.Info().Msg("Processing Ansible configuration...")
		err = processor.ProcessAnsibleTemplates(ec.Config)
		if err != nil {
			return err
		}
		return ec.executeAnsiblePlaybook()
	} else {
		log.Logger.Info().Msg("Skipping Ansible configuration")
	}

	return nil
}

func (ec *ExecutionConfig) executeTerraform() error {
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
	cmdOut.Stderr = os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to create pipe")
		return nil, err
	}
	cmdOut.Stdout = w
	err = cmdOut.Run()
	// r.Close()
	w.Close()
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to run Terraform output")
		return nil, nil
	}
	var buf bytes.Buffer
	numBytes, err := io.Copy(&buf, r)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to read TF output")
		return nil, err
	}
	var outputs map[string]interface{}
	if numBytes == 0 {
		log.Logger.Debug().Msg("output is empty")
		outputs = make(map[string]interface{})
	} else {
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
		ec.AnsiblePlaybookPath = ansibleCmdPath
	}
	if ec.Config.Ansible == nil || ec.Config.Ansible.InventoryFile == "" || ec.Config.Ansible.PlaybookFile == "" {
		log.Logger.Warn().Msg("Either Ansible inventory file or playbook were not specified. Aborting.")
		return nil
	}
	log.Logger.Debug().Msgf("Using ansible-playbook command: %s", ec.AnsiblePlaybookPath)
	log.Logger.Info().Msgf("Running ansible-playbook %s", ec.Config.Ansible.PlaybookFile)
	cmdPlaybook := osExec.Command(ec.AnsiblePlaybookPath, "-i", ec.Config.Ansible.InventoryFile, ec.Config.Ansible.PlaybookFile) //nolint:gosec
	cmdPlaybook.Stdout = os.Stdout
	cmdPlaybook.Stderr = os.Stderr
	cmdPlaybook.Dir = ec.AnsibleWorkDir
	cmdPlaybook.Env = append(cmdPlaybook.Env, os.Environ()...)
	// append custom roles dir if needed
	if ec.Config.TemplateConfig != nil && ec.Config.TemplateConfig.AnsibleRolesDir != "" {
		log.Logger.Debug().Msgf("Ansible roles directory: %s", ec.Config.TemplateConfig.AnsibleRolesDir)
		cmdPlaybook.Env = append(cmdPlaybook.Env, "ANSIBLE_ROLES_PATH="+ec.Config.TemplateConfig.AnsibleRolesDir)
	} else {
		log.Logger.Info().Msg("No custom Ansible roles directory specified")
	}
	err := cmdPlaybook.Run()
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to run ansible-playbook")
	}
	return err
}
