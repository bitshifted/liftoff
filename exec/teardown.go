// Copyright 2024 Bitshift D.O.O
// SPDX-License-Identifier: MPL-2.0

package exec

import (
	osExec "os/exec"
	"path"

	"github.com/bitshifted/liftoff/common"
	"github.com/bitshifted/liftoff/log"
)

func (ec *ExecutionConfig) ExecuteTeardown() error {
	output, err := ec.calculateOutputDirectory()
	if err != nil {
		return err
	}
	ec.TerraformWorkDir = path.Join(output, common.DefaultTerraformDir)
	if ec.TerraformPath == "" {
		tfCmdPath, e := osExec.LookPath(defaltTerraformCmd)
		if e != nil {
			log.Logger.Error().Err(e).Msg("Failed to lookup Terraform path")
			return e
		}
		ec.TerraformPath = tfCmdPath
	}
	err = ec.executeTerraformCommand("apply", "-destroy", "-auto-approve")
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to run Terraform destroy")
	}
	return err
}
