// Copyright 2024 Bitshift D.O.O
// SPDX-License-Identifier: MPL-2.0

package template

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/bitshifted/liftoff/common"
	"github.com/bitshifted/liftoff/config"
	"github.com/bitshifted/liftoff/log"
)

const (
	templateSuffix = ".tmpl"
)

type TemplateProcessor struct {
	BaseDir      string
	OutputDir    string
	TerraformDir string
	AnsibleDir   string
}

func (tp *TemplateProcessor) ProcessTemplates(conf *config.Configuration) error {
	// process Terraform templates
	if tp.TerraformDir == "" {
		tp.TerraformDir = common.DefaultTerraformDir
	}
	tfTemplateDir := path.Join(tp.BaseDir, tp.TerraformDir)
	log.Logger.Debug().Msgf("Terraform template directory: %s", tfTemplateDir)
	outputDir := tp.calculateOutputDirectory(common.DefaultTerraformDir)
	log.Logger.Info().Msgf("Terraform output directory: %s", outputDir)
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to create output directory")
		return err
	}
	err = tp.fileWalker(tfTemplateDir, conf)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to process Terraform templates")
		return err
	}
	// process Ansible templates
	if tp.AnsibleDir == "" {
		tp.AnsibleDir = common.DefaultAnsibleDir
	}
	ansibleTemplateDir := path.Join(tp.BaseDir, tp.AnsibleDir)
	log.Logger.Debug().Msgf("Ansible template directory: %s", ansibleTemplateDir)
	outputDir = tp.calculateOutputDirectory(common.DefaultAnsibleDir)
	log.Logger.Info().Msgf("Ansible output directory: %s", outputDir)
	err = os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to create output directory")
		return err
	}
	return tp.fileWalker(ansibleTemplateDir, conf)
}

func (tp *TemplateProcessor) processTemplate(templatePath string, conf *config.Configuration) error {
	log.Logger.Debug().Msgf("Processing template file %s", templatePath)
	tmpl, err := template.New(filepath.Base(templatePath)).Delims("[[", "]]").ParseFiles(templatePath)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to parse template")
		return err
	}
	outName := extractFileNameFromPath(templatePath)
	relPath, err := filepath.Rel(tp.BaseDir, outName)
	if err != nil {
		log.Logger.Error().Err(err).Msgf("Failed to find relative path for %s", outName)
		return err
	}
	outFilePath := path.Join(tp.OutputDir, relPath)
	log.Logger.Debug().Msgf("Output file path: %s", outFilePath)
	outFile, err := os.Create(outFilePath)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to create output template file")
		return err
	}
	return tmpl.Execute(outFile, conf)
}

func extractFileNameFromPath(filePath string) string {
	if strings.HasSuffix(filePath, templateSuffix) {
		return strings.Replace(filePath, templateSuffix, "", 1)
	}
	return filePath
}

func (tp *TemplateProcessor) fileWalker(templateDir string, conf *config.Configuration) error {
	return filepath.Walk(templateDir, func(fpath string, info os.FileInfo, err error) error {
		log.Logger.Debug().Msgf("Walk file %s", fpath)
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(tp.BaseDir, fpath)
		if err != nil {
			log.Logger.Error().Err(err).Msgf("Failed to find reataive path for %s", fpath)
			return err
		}
		if info.IsDir() {
			log.Logger.Debug().Msgf("Creating output directory %s", relPath)
			return os.MkdirAll(path.Join(tp.OutputDir, relPath), os.ModePerm)
		} else {
			return tp.processTemplate(path.Join(tp.BaseDir, relPath), conf)
		}
	})
}

func (tp *TemplateProcessor) calculateOutputDirectory(dirName string) string {
	var outputDir string
	if tp.OutputDir == "" {
		tp.OutputDir = common.DefaultOutputDir
		outputDir = path.Join(tp.BaseDir, tp.OutputDir, dirName)
	} else {
		if filepath.IsAbs(tp.OutputDir) {
			outputDir = path.Join(tp.OutputDir, dirName)
		} else {
			outputDir = path.Join(tp.BaseDir, tp.OutputDir, dirName)
		}
	}
	return outputDir
}
