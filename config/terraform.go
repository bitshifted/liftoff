package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"

	"github.com/bitshifted/easycloud/common"
	"github.com/bitshifted/easycloud/log"
)

type BackendType string

const (
	Local               BackendType = "local"
	Remote              BackendType = "remote"
	TerraformMinVersion             = "1.9.0"
)

type Terraform struct {
	Backend   *TerraformBackend  `yaml:"backend,omitempty"`
	Providers *TerraformProvider `yaml:"providers"`
}

func (t *Terraform) postLoad(config *Configuration) error {
	// post process configuration
	if t.Backend != nil {
		switch t.Backend.Type {
		case Local, Remote:
		default:
			return fmt.Errorf("invalid backend type: %s", t.Backend.Type)
		}
		err := t.Backend.postLoad(config)
		if err != nil {
			return err
		}
	}
	if t.Providers == nil {
		return errors.New("at least one Terraform provider is required")
	} else {
		return t.Providers.postLoad()
	}

}

type TerraformBackend struct {
	Type  BackendType   `yaml:"type"`
	Local *LocalBackend `yaml:"local,omitempty"`
}

func (tb *TerraformBackend) postLoad(config *Configuration) error {
	switch tb.Type {
	case Local:
		if tb.Local == nil {
			tb.Local = &LocalBackend{}
		}
	}
	if tb.Local != nil {
		return tb.Local.postLoad(config)
	}
	return nil
}

type LocalBackend struct {
	Path string `yaml:"path"`
}

func (lb *LocalBackend) postLoad(config *Configuration) error {
	if lb.Path == "" {
		finalPath, err := calculateLocalBackendPath(config)
		if err != nil {
			return err
		}
		lb.Path = finalPath
		return os.MkdirAll(finalPath, os.ModePerm)
	}
	return nil
}

func calculateLocalBackendPath(config *Configuration) (string, error) {
	// set default path to be in user home dir
	homeDir, err := osHomeDir()
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to get user home directory")
		return "", err
	}
	curPath := path.Join(homeDir, common.DefaultHomeDirName)
	if config.TemplateRepo != "" {
		templateUrl, err := url.Parse(config.TemplateRepo)
		if err != nil {
			log.Logger.Error().Err(err).Msgf("Failed to parse repository URL %s", config.TemplateRepo)
			return "", err
		}
		curPath = path.Join(curPath, templateUrl.Host, templateUrl.Path)
	}
	if config.TemplateDir != "" {
		curPath = path.Join(curPath, config.TemplateDir)
	}
	return curPath, nil
}

type TerraformProvider struct {
	Hcloud *ProviderHcloud `yaml:"hcloud,omitempty"`
}

func (tp *TerraformProvider) postLoad() error {
	if tp.Hcloud != nil {
		err := tp.Hcloud.postLoad()
		if err != nil {
			return err
		}
	}
	return nil
}
