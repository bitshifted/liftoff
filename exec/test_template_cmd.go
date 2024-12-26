// Copyright 2024 Bitshift D.O.O
// SPDX-License-Identifier: MPL-2.0

package exec

import (
	"errors"
	"os/exec"
	"path"

	"github.com/bitshifted/liftoff/common"
	"github.com/bitshifted/liftoff/config"
	"github.com/bitshifted/liftoff/log"
	"github.com/bitshifted/liftoff/template"
)

func (ec *ExecutionConfig) ExecuteTestTemplate() error {
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
	// run Terraform validattion
	if ec.TerraformPath == "" {
		tfCmdPath, e := exec.LookPath(defaltTerraformCmd)
		if e != nil {
			log.Logger.Error().Err(e).Msg("Failed to lookup Terraform path")
			return e
		}
		ec.TerraformPath = tfCmdPath
	}
	log.Logger.Info().Msg("Running Terraform init...")
	err = ec.executeTerraformCommand("init")
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to run Terraform init")
		return err
	}
	log.Logger.Info().Msg("Running Terraform validate...")
	err = ec.executeTerraformCommand("validate")
	if err != nil {
		log.Logger.Error().Err(err).Msg("Terraform validation failed")
		return err
	}
	// run Terraform plan
	log.Logger.Info().Msg("Running Terraform plan...")
	err = ec.executeTerraformCommand("plan")
	if err != nil {
		log.Logger.Error().Err(err).Msg("Terraform plan failed")
		return err
	} else {
		log.Logger.Info().Msg("Terraform validation successful!")
	}
	err = processor.ProcessAnsibleTemplates(ec.Config)
	if err != nil {
		return err
	}
	return err
}
