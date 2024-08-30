// Copyright 2024 Bitshift D.O.O
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"testing"

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
