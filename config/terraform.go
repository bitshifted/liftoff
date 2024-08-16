// Copyright 2024 Bitshift D.O.O
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/bitshifted/liftoff/common"
	"github.com/bitshifted/liftoff/config/tfprovider"
	"github.com/bitshifted/liftoff/log"
)

type BackendType string

const (
	Local                     BackendType = "local"
	Remote                    BackendType = "remote"
	TerraformMinVersion                   = "1.9.0"
	defaultTfStateFileName                = "terraform.tfstate"
	defaultTfWorkspaceDirName             = "terraform.tf.d"
)

type Terraform struct {
	Backend   *TerraformBackend              `yaml:"backend,omitempty"`
	Providers *tfprovider.TerraformProviders `yaml:"providers"`
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
		return t.Providers.PostLoad()
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
	case Remote:
		// add logic here
	}
	if tb.Local != nil {
		return tb.Local.postLoad(config)
	}
	return nil
}

type LocalBackend struct {
	Path      string `yaml:"path"`
	Workspace string `yaml:"workspace"`
}

func (lb *LocalBackend) postLoad(config *Configuration) error {
	var err error
	defaultDir, err := calculateLocalBackendDefaultDir(config)
	if err != nil {
		return err
	}
	if lb.Path == "" {
		lb.Path = path.Join(defaultDir, defaultTfStateFileName)
		// create required directories
		log.Logger.Debug().Msgf("Creating required directories for local backend: %s", filepath.Dir(lb.Path))
		err = os.MkdirAll(filepath.Dir(lb.Path), os.ModePerm)
	}
	if lb.Workspace == "" {
		lb.Workspace = path.Join(defaultDir, defaultTfWorkspaceDirName)
		// create required directories
		log.Logger.Debug().Msgf("Creating required directories for local workspace: %s", lb.Workspace)
		err = os.MkdirAll(lb.Workspace, os.ModePerm)
	}
	return err
}

func calculateLocalBackendDefaultDir(config *Configuration) (string, error) {
	homeDir, err := osHomeDir()
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to get user home directory")
		return "", err
	}
	curPath := path.Join(homeDir, common.DefaultHomeDirName)
	if config.TemplateRepo != "" {
		templateURL, err := url.Parse(config.TemplateRepo)
		if err != nil {
			log.Logger.Error().Err(err).Msgf("Failed to parse repository URL %s", config.TemplateRepo)
			return "", err
		}
		curPath = path.Join(curPath, templateURL.Host, templateURL.Path)
	}
	if config.TemplateDir != "" {
		curPath = path.Join(curPath, config.TemplateDir)
	}

	return path.Join(curPath), nil
}
