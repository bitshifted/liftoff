// Copyright 2024 Bitshift D.O.O
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfigShouldPass(t *testing.T) {
	config, err := LoadConfig("./test_files/simple-config.yaml")
	assert.NoError(t, err)
	assert.NotNil(t, config.Terraform)
	assert.NotNil(t, config.Terraform.Backend)
	assert.Equal(t, Local, config.Terraform.Backend.Type)

	// ansible config
	assert.NotNil(t, config.Ansible)
	assert.Equal(t, "my-inventory", config.Ansible.InventoryFile)
	assert.Equal(t, "myplaybook.yaml", config.Ansible.PlaybookFile)
	// verify variables
	assert.NotNil(t, config.Variables)
	// hcloud variables
	hcloudvars, ok := config.Variables["hcloud"]
	assert.True(t, ok)
	textvar, ok := hcloudvars["textvar"]
	assert.True(t, ok)
	assert.NotNil(t, textvar)
	assert.Equal(t, "some text", textvar)
	intvar, ok := hcloudvars["intvar"]
	assert.True(t, ok)
	assert.NotNil(t, intvar)
	assert.Equal(t, 123, intvar)
	boolvar, ok := hcloudvars["boolvar"]
	assert.True(t, ok)
	assert.NotNil(t, boolvar)
	assert.Equal(t, true, boolvar)
	// hetznerdns vars
	hdns, ok := config.Variables["hetznerdns"]
	assert.True(t, ok)
	complexvar, ok := hdns["complexvar"]
	assert.True(t, ok)
	assert.NotNil(t, complexvar)
	assert.Equal(t, "string property", complexvar.(map[string]interface{})["stringprop"])
	assert.Equal(t, 3.14, complexvar.(map[string]interface{})["floatprop"])
	// digitalocean vars
	dovars, ok := config.Variables["digitalocean"]
	assert.True(t, ok)
	listvar, ok := dovars["listvar"]
	assert.True(t, ok)
	assert.NotNil(t, listvar)
	data, ok := listvar.([]interface{})
	assert.True(t, ok)
	assert.Equal(t, 2, len(data))
	item := data[0].(map[string]interface{})
	assert.NotNil(t, item)
	assert.Equal(t, "item1", item["name"])
	assert.Equal(t, "bar-item-1", item["foo"])
	// validate tags
	tags := config.Tags
	assert.NotNil(t, tags)
	assert.Equal(t, 2, len(tags))
	assert.Equal(t, "foovalue", tags["foo"])
	assert.Equal(t, "barvalue", tags["bar"])
}

func TestShouldErrorForInvalidBackendType(t *testing.T) {
	_, err := LoadConfig("./test_files/invalid-backend-config.yaml")
	assert.Error(t, err)
	assert.Equal(t, "invalid backend type: foo", err.Error())
}

func TestShouldReturnNilWhenNoTemplateConfig(t *testing.T) {
	config, err := LoadTemplateConfig("./test_files")
	assert.NoError(t, err)
	assert.Nil(t, config)
}

func TestShouldReturnAbsolutePathForTemplateRelativePaths(t *testing.T) {
	absPath, err := filepath.Abs("./test_files/tmpl-test-dir")
	assert.NoError(t, err)
	config, err := LoadTemplateConfig(absPath)
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, path.Join(absPath, "tf-extra-dir"), config.TerraformExtraDir)
	assert.Equal(t, path.Join(absPath, "../roles-dir"), config.AnsibleRolesDir)
}

func TestShouldReturnSameForAbsolutePath(t *testing.T) {
	absPath, err := filepath.Abs("./test_files/abs-tmpl-dir")
	assert.NoError(t, err)
	config, err := LoadTemplateConfig(absPath)
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "/tmp/terraform", config.TerraformExtraDir)
	assert.Equal(t, path.Join(absPath, "../roles-dir"), config.AnsibleRolesDir)
}
