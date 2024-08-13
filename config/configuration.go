package config

import (
	"os"

	"github.com/bitshifted/liftoff/log"
	"gopkg.in/yaml.v3"
)

const (
	defaultTemplateRepo = "https://github.com/bitshifted/autoconfig-templates"
)

type Configuration struct {
	TemplateRepo   string     `yaml:"template-repo,omitempty"`
	TempateVersion string     `yaml:"template-version,omitempty"`
	TemplateDir    string     `yaml:"template-dir,omitempty"`
	Terraform      *Terraform `yaml:"terraform,omitempty"`
}

func LoadConfig(path string) (*Configuration, error) {
	var config Configuration
	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to load configuration file")
		return nil, err
	}
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		log.Logger.Error().Err(err).Msgf("Failed to parse configuration file %s", path)
		return nil, err
	}
	err = config.postLoad()
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (c *Configuration) postLoad() error {
	if c.TemplateRepo == "" {
		c.TemplateRepo = defaultTemplateRepo
	}
	return c.Terraform.postLoad(c)
}
