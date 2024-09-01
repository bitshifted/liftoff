// Copyright 2024 Bitshift D.O.O
// SPDX-License-Identifier: MPL-2.0

package template

import (
	"bufio"
	"errors"
	"fmt"
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
	templateSuffix     = ".tmpl"
	generatedFilesName = ".genfiles"
	fileMode           = 0644
)

type TemplateProcessor struct {
	BaseDir        string
	OutputDir      string
	TerraformDir   string
	AnsibleDir     string
	generatedFiles []string
}

func (tp *TemplateProcessor) ProcessTerraformTemplates(conf *config.Configuration) error {
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
	err = tp.writeGeneratedFilePaths(outputDir)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to write generated file paths")
	}
	return err
}

func (tp *TemplateProcessor) ProcessAnsibleTemplates(conf *config.Configuration) error {
	if tp.AnsibleDir == "" {
		tp.AnsibleDir = common.DefaultAnsibleDir
	}
	ansibleTemplateDir := path.Join(tp.BaseDir, tp.AnsibleDir)
	log.Logger.Debug().Msgf("Ansible template directory: %s", ansibleTemplateDir)
	outputDir := tp.calculateOutputDirectory(common.DefaultAnsibleDir)
	log.Logger.Info().Msgf("Ansible output directory: %s", outputDir)
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to create output directory")
		return err
	}
	_, err = os.Stat(ansibleTemplateDir)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		log.Logger.Warn().Msgf("Ansible template directory %s does not exist. Skipping", ansibleTemplateDir)
		return nil
	}
	tp.generatedFiles = []string{}
	err = tp.fileWalker(ansibleTemplateDir, conf)
	if err != nil {
		return err
	}
	err = tp.writeGeneratedFilePaths(outputDir)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to write generated file paths")
	}
	return err
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
	tp.generatedFiles = append(tp.generatedFiles, outFilePath)
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

func (tp *TemplateProcessor) cleanupGeneratedFiles(baseDir string) error {
	log.Logger.Debug().Msgf("Cleaning up files in %s", baseDir)
	inFile, err := os.OpenFile(path.Join(baseDir, generatedFilesName), os.O_RDONLY, fileMode)
	if err != nil {
		log.Logger.Warn().Err(err).Msg("Failed to open .genfiles for reading")
		return err
	}
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	for scanner.Scan() {
		toDelete := true
		scanned := scanner.Text()
		for _, curFile := range tp.generatedFiles {
			if curFile == scanned {
				toDelete = false
				break
			}
		}
		if toDelete {
			log.Logger.Info().Msgf("Deleting redundant file %s", scanned)
			derr := os.Remove(scanned)
			if err != nil {
				log.Logger.Error().Err(derr).Msgf("Failed to delete file %s", scanned)
			}
		}
	}
	return nil
}

func (tp *TemplateProcessor) writeGeneratedFilePaths(baseDir string) error {
	err := tp.cleanupGeneratedFiles(baseDir)
	if err != nil {
		log.Logger.Warn().Err(err).Msg("Error cleaning generated files")
	}
	outFile, err := os.OpenFile(path.Join(baseDir, generatedFilesName), os.O_RDWR|os.O_CREATE|os.O_TRUNC, fileMode)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to open file for writing")
		return err
	}
	defer outFile.Close()
	writer := bufio.NewWriter(outFile)
	for _, fpath := range tp.generatedFiles {
		txt := fmt.Sprintf("%s\n", fpath)
		_, err = writer.WriteString(txt)
		if err != nil {
			log.Logger.Warn().Msgf("Could not write generated file path: %s", fpath)
		}
	}
	return writer.Flush()
}
