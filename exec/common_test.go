// Copyright 2024 Bitshift D.O.O
// SPDX-License-Identifier: MPL-2.0

package exec

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bitshifted/liftoff/config"
	"github.com/bitshifted/liftoff/log"
	"github.com/stretchr/testify/suite"
)

type ExecutionConfigTestSuite struct {
	suite.Suite
}

func (ts *ExecutionConfigTestSuite) SetupSuite() {
	// Initialize logger
	log.Init(true)
	log.Logger.Info().Msg("Running ExecutionConfigTestSuite")
}

func TestExecutionConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ExecutionConfigTestSuite))
}

func (ts *ExecutionConfigTestSuite) TestCalculateOutputDirectory() {
	tempDir := ts.T().TempDir()
	configFilePath := filepath.Join(tempDir, "test-config.yaml")
	_, err := os.Create(configFilePath)
	ts.NoError(err)

	ec := &ExecutionConfig{
		ConfigFilePath: configFilePath,
	}

	outputDir, err := ec.calculateOutputDirectory()
	ts.NoError(err)
	ts.Equal(filepath.Join(tempDir, "test-config"), outputDir)

	_, err = os.Stat(outputDir)
	ts.NoError(err)
}

func (ts *ExecutionConfigTestSuite) TestExecutionConfig_templateDirAbsPath() {
	// tempDir := ts.T().TempDir()
	config := &config.Configuration{
		TemplateRepo: "https://github.com/bitshifted/liftoff-templates.git",
		TemplateDir:  "tmpl-dir",
	}

	ec := &ExecutionConfig{
		Config: config,
	}

	tmplDir, err := ec.templateDirAbsPath()
	ts.NoError(err)
	ts.Contains(tmplDir, "template_repo")
	ts.True(strings.HasSuffix(tmplDir, "tmpl-dir"))
}

func (ts *ExecutionConfigTestSuite) TestTemplateDirAbsPath_NoRepo() {
	config := &config.Configuration{}

	ec := &ExecutionConfig{
		Config: config,
	}

	tmplDir, err := ec.templateDirAbsPath()
	ts.NoError(err)
	ts.Equal("", tmplDir)
}
