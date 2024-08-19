// Copyright 2024 Bitshift D.O.O
// SPDX-License-Identifier: MPL-2.0

package exec

import (
	"errors"
	"os"
	"os/exec"
	"path"

	"github.com/bitshifted/liftoff/common"
	"github.com/bitshifted/liftoff/gitops"
	"github.com/bitshifted/liftoff/log"
	"github.com/bitshifted/liftoff/template"
)

func (ec *ExecutionConfig) ExecuteTestTemplate() error {
	repo := ec.Config.TemplateRepo
	if repo == "" {
		log.Logger.Info().Msg("Template repository not specified")
	} else {
		tmpDir, err := os.MkdirTemp("", "template_repo")
		if err != nil {
			log.Logger.Error().Err(err).Msg("Failed to create temprorary directory for clone")
		}
		handler := gitops.GitHandler{
			URL:         repo,
			Version:     ec.Config.TempateVersion,
			Destination: tmpDir,
		}
		err = handler.Fetch()
		if err != nil {
			return err
		}
	}
	tmplDir := ec.Config.TemplateDir
	if tmplDir == "" {
		return errors.New("either template repository or template directory must be specified")
	}
	log.Logger.Info().Msgf("Template directory: %s", tmplDir)
	output, err := os.MkdirTemp("", tmpDirPrefix)
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
	// run Terraform validattion
	if ec.TerraformPath == "" {
		tfCmdPath, e := exec.LookPath(defaltTerraformCmd)
		if e != nil {
			log.Logger.Error().Err(e).Msg("Failed to lookup Terraform path")
			return e
		}
		ec.TerraformPath = tfCmdPath
	}
	cmdInit := exec.Command(ec.TerraformPath, "init") //nolint:gosec
	cmdInit.Stdout = os.Stdout
	cmdInit.Stderr = os.Stderr
	cmdInit.Dir = ec.TerraformWorkDir
	cmdInit.Env = append(cmdInit.Env, "TF_DATA_DIR="+ec.Config.Terraform.DataDir)
	log.Logger.Debug().Msgf("Terraform work directory: %s", cmdInit.Dir)
	log.Logger.Info().Msg("Running Terraform init...")
	err = cmdInit.Run()
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to run Terraform init")
		return err
	}
	cmdValidate := exec.Command(ec.TerraformPath, "validate") //nolint:gosec
	cmdValidate.Stdout = os.Stdout
	cmdValidate.Stderr = os.Stderr
	cmdValidate.Dir = ec.TerraformWorkDir
	cmdValidate.Env = append(cmdValidate.Env, "TF_DATA_DIR="+ec.Config.Terraform.DataDir)
	log.Logger.Info().Msg("Running Terraform validate...")
	err = cmdValidate.Run()
	if err != nil {
		log.Logger.Error().Err(err).Msg("Terraform validation failed")
	} else {
		log.Logger.Info().Msg("Terraform validation successful!")
	}
	return err
}
