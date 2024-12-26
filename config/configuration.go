// Copyright 2024 Bitshift D.O.O
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"os"
	"path"
	"path/filepath"

	"github.com/bitshifted/liftoff/log"
	"gopkg.in/yaml.v3"
)

const (
	templateConfigFileName = "template-cfg.yaml"
)

type Configuration struct {
	TemplateRepo   string            `yaml:"template-repo,omitempty"`
	TempateVersion string            `yaml:"template-version,omitempty"`
	TemplateDir    string            `yaml:"template-dir,omitempty"`
	Terraform      *Terraform        `yaml:"terraform,omitempty"`
	Ansible        *AnsibleConfig    `yaml:"ansible,omitempty"`
	Variables      ConfigVariables   `yaml:"variables"`
	Tags           map[string]string `yaml:"tags"`
	ProcessingVars map[string]interface{}
	TemplateConfig *TemplateConfig
}

type TemplateConfig struct {
	TerraformExtraDir string `yaml:"terraform-extra-dir,omitempty"`
	AnsibleRolesDir   string `yaml:"ansible-roles-dir,omitempty"`
}

func LoadConfig(configPath string) (*Configuration, error) {
	var config Configuration
	bytes, err := os.ReadFile(configPath)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to load configuration file")
		return nil, err
	}
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		log.Logger.Error().Err(err).Msgf("Failed to parse configuration file %s", configPath)
		return nil, err
	}
	err = config.postLoad()
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func LoadTemplateConfig(templateDir string) (*TemplateConfig, error) {
	tmplConfigPath := path.Join(templateDir, templateConfigFileName)
	if _, err := os.Stat(tmplConfigPath); os.IsNotExist(err) {
		log.Logger.Info().Msgf("Template config file does not exist: %s", tmplConfigPath)
		return nil, nil
	}
	bytes, err := os.ReadFile(tmplConfigPath)
	if err != nil {
		log.Logger.Error().Err(err).Msgf("Failed to read template config file %s", tmplConfigPath)
		return nil, err
	}
	var tmplConfig TemplateConfig
	err = yaml.Unmarshal(bytes, &tmplConfig)
	if err != nil {
		log.Logger.Error().Err(err).Msgf("Failed to parse template config file %s", tmplConfigPath)
		return nil, err
	}
	// convert paths to absolute paths
	tfExtraDir := tmplConfig.TerraformExtraDir
	if tfExtraDir != "" && !filepath.IsAbs(tfExtraDir) {
		tmplConfig.TerraformExtraDir = path.Join(templateDir, tfExtraDir)
	}
	ansibleRolesDir := tmplConfig.AnsibleRolesDir
	if ansibleRolesDir != "" && !filepath.IsAbs(ansibleRolesDir) {
		tmplConfig.AnsibleRolesDir = path.Join(templateDir, ansibleRolesDir)
	}
	return &tmplConfig, nil
}

func (c *Configuration) postLoad() error {
	c.ProcessingVars = c.Variables.forEnvironment()
	err := processVariables(c.ProcessingVars)
	if err != nil {
		return err
	}
	return c.Terraform.postLoad()
}
