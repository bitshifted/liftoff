// Copyright 2024 Bitshift D.O.O
// SPDX-License-Identifier: MPL-2.0

package tfprovider

import (
	"errors"

	"github.com/bitshifted/liftoff/common"
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
	// replace values if needed
	replacement, err := common.ProcessStringValue(hc.Token)
	if err != nil {
		return err
	}
	hc.Token = replacement
	replacement, err = common.ProcessStringValue(hc.Endpoint)
	if err != nil {
		return err
	}
	hc.Endpoint = replacement
	replacement, err = common.ProcessStringValue(hc.PollInterval)
	if err != nil {
		return err
	}
	hc.PollInterval = replacement
	replacement, err = common.ProcessStringValue(hc.PollFunction)
	if err != nil {
		return err
	}
	hc.PollFunction = replacement
	return nil
}
