// Copyright 2024 Bitshift D.O.O
// SPDX-License-Identifier: MPL-2.0

package exec

import (
	"errors"
	"os"
	"path"

	"github.com/bitshifted/liftoff/common"
	"github.com/bitshifted/liftoff/log"
	"github.com/bitshifted/liftoff/template"
)

func (ec *ExecutionConfig) ExecuteTestTemplate() error {
	repo := ec.Config.TemplateRepo
	if repo == "" {
		log.Logger.Info().Msg("Template repository not specified")
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
	return processor.ProcessTemplates(ec.Config)
}
