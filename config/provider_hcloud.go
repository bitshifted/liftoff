// Copyright 2024 Bitshift D.O.O
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"errors"
	"fmt"

	"github.com/bitshifted/liftoff/common"
	"github.com/bitshifted/liftoff/log"
)

type ProviderHcloud struct {
	Token        string `yaml:"token"`
	Endpoint     string `yaml:"endpoint,omitempty"`
	PollInterval string `yaml:"poll-interval,omitempty"`
	PollFunction string `yaml:"poll-function,omitempty"`
}

func (hc *ProviderHcloud) postLoad() error {
	if hc.Token == "" {
		return errors.New("token is required for hcloud provider")
	}
	valueType := common.ValueTypeFromString(hc.Token)
	switch valueType {
	case common.EnvVariableString:
		varName, err := common.ExtractEnvVarName(hc.Token)
		if err != nil {
			return err
		}
		if !common.IsEnvVariableSet(varName) {
			err := fmt.Errorf("referenced environment variable %s is not set", varName)
			log.Logger.Error().Err(err).Msg("Environment variable is not set")
			return err
		} else {
			hc.Token = ""
		}
	case common.FileContentString:
		// add logic here
	}
	return nil
}
