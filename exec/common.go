// Copyright 2024 Bitshift D.O.O
// SPDX-License-Identifier: MPL-2.0

package exec

import "github.com/bitshifted/liftoff/config"

const (
	defaltTerraformCmd = "terraform"
	defaultAnsibleCmd  = "ansible-playbook"
)

type ExecutionConfig struct {
	Config              *config.Configuration
	SkipTerraform       bool
	SkipAnsible         bool
	TerraformPath       string
	AnsiblePlaybookPath string
	TerraformWorkDir    string
	AnsibleWorkDir      string
}
