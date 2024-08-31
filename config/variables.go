// Copyright 2024 Bitshift D.O.O
// SPDX-License-Identifier: MPL-2.0

package config

const (
	envDefault = "default"
)

type ConfigVariables map[string]map[string]interface{}

// returns variables for specific environment
// TODO this should take environment name as parameters
func (cv ConfigVariables) forEnvironment() map[string]interface{} {
	return cv[envDefault]
}
