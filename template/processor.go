package template

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/bitshifted/easycloud/common"
	"github.com/bitshifted/easycloud/config"
	"github.com/bitshifted/easycloud/log"
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

func (tp *TemplateProcessor) ProcessTerraformTemplate(config *config.Configuration) error {
	if tp.TerraformDir == "" {
		tp.TerraformDir = common.DefaultTerraformDir
	}
	templateDir := path.Join(tp.BaseDir, tp.TerraformDir)
	log.Logger.Debug().Msgf("template directory: %s", templateDir)
	var outputDir string
	if tp.OutputDir == "" {
		tp.OutputDir = common.DefaultOutputDir
		outputDir = path.Join(tp.BaseDir, tp.OutputDir, common.DefaultTerraformDir)
	} else {
		if filepath.IsAbs(tp.OutputDir) {
			outputDir = path.Join(tp.OutputDir, common.DefaultTerraformDir)
		} else {
			outputDir = path.Join(tp.BaseDir, tp.OutputDir, common.DefaultTerraformDir)
		}
	}
	log.Logger.Info().Msgf("Output directory: %s", outputDir)
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to create output directory")
		return err
	}
	err = filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		log.Logger.Debug().Msgf("Walk file %s", path)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			processTemplate(path, config, outputDir)
		}
		return nil
	})
	return err
}

func processTemplate(templatePath string, config *config.Configuration, outputDir string) error {
	log.Logger.Debug().Msgf("Processing template file %s", templatePath)
	tmpl, err := template.New(filepath.Base(templatePath)).Delims("[[", "]]").ParseFiles(templatePath)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to parse template")
		return err
	}
	outName, err := extractFileNameFromPath(templatePath)
	if err != nil {
		return err
	}
	outFilePath := path.Join(outputDir, outName)
	log.Logger.Debug().Msgf("Output file path: %s", outFilePath)
	outFile, err := os.Create(outFilePath)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to create output template file")
		return err
	}
	return tmpl.Execute(outFile, config)
}

func extractFileNameFromPath(filePath string) (string, error) {
	name := filepath.Base(filePath)
	return strings.Split(name, templateSuffix)[0], nil
}
