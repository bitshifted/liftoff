// Copyright 2024 Bitshift D.O.O
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfigShouldPass(t *testing.T) {
	config, err := LoadConfig("./test_files/simple-config.yaml")
	assert.NoError(t, err)
	assert.NotNil(t, config.Terraform)
	assert.NotNil(t, config.Terraform.Backend)
	assert.Equal(t, Local, config.Terraform.Backend.Type)
	localPath, err := calculateLocalBackendDefaultDir(config)
	assert.NoError(t, err)

	assert.Equal(t, path.Join(localPath, defaultTfStateFileName), config.Terraform.Backend.Local.Path)
	assert.Equal(t, path.Join(localPath, defaultTfWorkspaceDirName), config.Terraform.Backend.Local.Workspace)
	// ansible config
	assert.NotNil(t, config.Ansible)
	assert.Equal(t, "my-inventory", config.Ansible.InventoryFile)
	assert.Equal(t, "myplaybook.yaml", config.Ansible.PlaybookFile)
}

func TestShouldErrorForInvalidBackendType(t *testing.T) {
	_, err := LoadConfig("./test_files/invalid-backend-config.yaml")
	assert.Error(t, err)
	assert.Equal(t, "invalid backend type: foo", err.Error())
}
