package template

import (
	"os"
	"path"
	"testing"

	"github.com/bitshifted/easycloud/config"
	"github.com/bitshifted/easycloud/log"
	"github.com/stretchr/testify/assert"
)

func TestProcessTerraformFiles(t *testing.T) {
	log.Init()
	tmpDir, err := os.MkdirTemp("", "template-test")
	assert.NoError(t, err)
	processor := TemplateProcessor{
		BaseDir:   "test_files",
		OutputDir: tmpDir,
	}
	config, err := config.LoadConfig("test_files/sample-config.yaml")
	assert.NoError(t, err)
	err = processor.ProcessTerraformTemplate(config)
	assert.NoError(t, err)
	_, err = os.Stat(path.Join(tmpDir, defaultTerraformDir, "terraform.tf"))
	assert.NoError(t, err)
}
