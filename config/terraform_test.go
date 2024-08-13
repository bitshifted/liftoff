package config

import (
	"os"
	"path"
	"testing"

	"github.com/bitshifted/easycloud/common"
	"github.com/bitshifted/easycloud/log"
	"github.com/stretchr/testify/suite"
)

type TerraformTestSuite struct {
	suite.Suite
}

func (ts *TerraformTestSuite) SetupSuite() {
	// Initialize logger
	log.Init()
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
	finalPath, err := calculateLocalBackendPath(&conf)
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
	finalPath, err := calculateLocalBackendPath(&conf)
	ts.NoError(err)
	ts.Equal(path.Join(tmpDir, common.DefaultHomeDirName, "some", "dir", "path"), finalPath)
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
	finalPath, err := calculateLocalBackendPath(&conf)
	ts.NoError(err)
	ts.Equal(path.Join(tmpDir, common.DefaultHomeDirName, "github.com", "my", "repo.git", "some", "dir"), finalPath)
}
