// Copyright 2024 Bitshift D.O.O
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"errors"
	"fmt"
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
	providerCloudflare   = "cloudflare"
)

var supportedProviders = []string{providerHcloud, providerHetznerdns, providerDigitalOcean, providerCloudflare}

type Terraform struct {
	Backend   *TerraformBackend `yaml:"backend,omitempty"`
	Providers []string          `yaml:"providers"`
}

func (t *Terraform) HasProvider(name string) bool {
	for _, prov := range t.Providers {
		if name == prov {
			return true
		}
	}
	return false
}

func (t *Terraform) postLoad() error {
	// post process configuration
	if t.Backend != nil {
		switch t.Backend.Type {
		case Local, Remote:
		default:
			return fmt.Errorf("invalid backend type: %s", t.Backend.Type)
		}
		err := t.Backend.postLoad()
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

func (tb *TerraformBackend) postLoad() error {
	switch tb.Type {
	case Local:
		if tb.Local == nil {
			tb.Local = &LocalBackend{}
		}
	case Remote:
		// add logic here
	}
	return nil
}

type LocalBackend struct {
	Path      string `yaml:"path"`
	Workspace string `yaml:"workspace"`
}
