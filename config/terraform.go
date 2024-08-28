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
	"github.com/bitshifted/liftoff/log"
)

type BackendType string

const (
	Local                      BackendType = "local"
	Remote                     BackendType = "remote"
	TerraformMinVersion                    = "1.9.0"
	defaultTfStateFileName                 = "terraform.tfstate"
	defaultTfWorkspaceDirName              = "terraform.tf.d"
	defaultTerraformDatDirName             = ".terraform"
	// providers
	providerHcloud       = "hcloud"
	providerHetznerdns   = "hetznerdns"
	providerDigitalOcean = "digitalocean"
)

var supportedProviders = []string{providerHcloud, providerHetznerdns, providerDigitalOcean}

type Terraform struct {
	Backend   *TerraformBackend `yaml:"backend,omitempty"`
	Providers []string          `yaml:"providers"`
	DataDir   string
}

func (t *Terraform) HasProvider(name string) bool {
	for _, prov := range t.Providers {
		if name == prov {
			return true
		}
	}
	return false
}

func (t *Terraform) postLoad(config *Configuration) error {
	dataDir, err := calculateTFBaseDir(config)
	log.Logger.Debug().Msgf("Terraform data directory: %s", dataDir)
	if err != nil {
		return err
	}
	t.DataDir = dataDir
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
	if len(t.Providers) == 0 {
		return errors.New("at least one Terraform provider is required")
	}
	return t.checkSupportedProviders()
}

func (t *Terraform) checkSupportedProviders() error {
	for _, prov := range t.Providers {
		if !t.isSupportedProvider(prov) {
			return fmt.Errorf("provider '%s' is not supported", prov)
		}
	}
	return nil
}

func (t *Terraform) isSupportedProvider(provider string) bool {
	for _, prov := range supportedProviders {
		if provider == prov {
			return true
		}
	}
	return false
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
	defaultDir, err := calculateTFBaseDir(config)
	if err != nil {
		return err
	}
	if lb.Path == "" {
		lb.Path = path.Join(filepath.Dir(defaultDir), defaultTfStateFileName)
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

func calculateTFBaseDir(config *Configuration) (string, error) {
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
		curPath = path.Join(curPath, config.TemplateDir, defaultTerraformDatDirName)
	}

	return path.Join(curPath), nil
}
