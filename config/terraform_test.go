// Copyright 2024 Bitshift D.O.O
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"os"
	"path"
	"testing"

	"github.com/bitshifted/liftoff/common"
	"github.com/bitshifted/liftoff/log"
	"github.com/stretchr/testify/suite"
)

type TerraformTestSuite struct {
	suite.Suite
}

func (ts *TerraformTestSuite) SetupSuite() {
	// Initialize logger
	log.Init(true)
	log.Logger.Info().Msg("Running TerraformTestSuite")
}

func TestTerraformTestSuite(t *testing.T) {
	suite.Run(t, new(TerraformTestSuite))
}

func (ts *TerraformTestSuite) TestLocalBackendWhenTemplateRepoSpecified() {
	tmpDir, err := os.MkdirTemp("", "tf_test_suite")
	ts.NoError(err)
	osHomeDir = func() (string, error) {
		return tmpDir, nil
	}
	conf := Configuration{
		TemplateRepo: "https://github.com/my/repo.git",
	}
	finalPath, err := calculateTFBaseDir(&conf)
	ts.NoError(err)
	ts.Equal(path.Join(tmpDir, common.DefaultHomeDirName, "github.com", "my", "repo.git"), finalPath)
}

func (ts *TerraformTestSuite) TestLocalBackendWhenOnlyDirectorySet() {
	tmpDir, err := os.MkdirTemp("", "tf_test_suite")
	ts.NoError(err)
	osHomeDir = func() (string, error) {
		return tmpDir, nil
	}
	conf := Configuration{
		TemplateDir: "/some/dir/path",
	}
	finalPath, err := calculateTFBaseDir(&conf)
	ts.NoError(err)
	ts.Equal(path.Join(tmpDir, common.DefaultHomeDirName, "some", "dir", "path", defaultTerraformDatDirName), finalPath)
}

func (ts *TerraformTestSuite) TestBackendPathBothRepoDirSet() {
	tmpDir, err := os.MkdirTemp("", "tf_test_suite")
	ts.NoError(err)
	osHomeDir = func() (string, error) {
		return tmpDir, nil
	}
	conf := Configuration{
		TemplateRepo: "https://github.com/my/repo.git",
		TemplateDir:  "some/dir",
	}
	finalPath, err := calculateTFBaseDir(&conf)
	ts.NoError(err)
	ts.Equal(path.Join(tmpDir, common.DefaultHomeDirName, "github.com", "my", "repo.git", "some", "dir", defaultTerraformDatDirName),
		finalPath)
}

func (ts *TerraformTestSuite) HasProviderReturnsCorrectValue() {
	conf := Configuration{
		Terraform: &Terraform{
			Providers: []string{"hcloud", "hetznerdns"},
		},
	}
	hasProvider := conf.Terraform.HasProvider("hetznerdns")
	ts.True(hasProvider)
	hasProvider = conf.Terraform.HasProvider("hcloud")
	ts.False(hasProvider)
}

func (ts *TerraformTestSuite) ShouldFailOnUnsupportedProvider() {
	tf := Terraform{
		Providers: []string{providerHcloud, providerHetznerdns, "foo"},
	}
	err := tf.checkSupportedProviders()
	ts.Error(err)
}
