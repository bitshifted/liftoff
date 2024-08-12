package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfigShouldPass(t *testing.T) {
	config, err := LoadConfig("./test_files/simple-config.yaml")
	assert.NoError(t, err)
	assert.Equal(t, defaultTemplateRepo, config.TemplateRepo)
	assert.NotNil(t, config.Terraform)
	assert.NotNil(t, config.Terraform.Backend)
	assert.Equal(t, Local, config.Terraform.Backend.Type)
	localPath, err := calculateLocalBackendPath(config)
	assert.NoError(t, err)

	assert.Equal(t, localPath, config.Terraform.Backend.Local.Path)
}

func TestShouldErrorForInvalidBackendType(t *testing.T) {
	_, err := LoadConfig("./test_files/invalid-backend-config.yaml")
	assert.Error(t, err)
	assert.Equal(t, "invalid backend type: foo", err.Error())
}