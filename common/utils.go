// Copyright 2024 Bitshift D.O.O
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"fmt"
	"os"
	"strings"

	"github.com/bitshifted/liftoff/log"
)

type ValueType int8

const (
	envPrefix                   = "fromenv:"
	contentPrefix               = "fromfile:"
	EnvVariableString ValueType = iota
	FileContentString
	PlainString
)

func ValueTypeFromString(input string) ValueType {
	if strings.HasPrefix(input, envPrefix) {
		log.Logger.Debug().Msgf("Found environment variable setting: %s", input)
		return EnvVariableString
	}
	// check if it's file content reference
	if strings.HasPrefix(input, contentPrefix) {
		log.Logger.Debug().Msgf("Found file content setting: %s", input)
		return FileContentString
	}
	return PlainString
}

func IsEnvVariableSet(envVarName string) bool {
	value := osGetEnv(envVarName)
	return value != ""
}

func ExtractEnvVarName(input string) string {
	return strings.Replace(input, envPrefix, "", 1)
}

func ProcessStringValue(input string) (string, error) {
	valueType := ValueTypeFromString(input)
	switch valueType {
	case EnvVariableString:
		// extract variable name and check if it is set
		varName := ExtractEnvVarName(input)
		if !IsEnvVariableSet(varName) {
			err := fmt.Errorf("referenced environment variable %s is not set", varName)
			log.Logger.Error().Err(err).Msg("Environment variable is not set")
			return "", err
		} else {
			return osGetEnv(varName), nil
		}
	case FileContentString:
		// strip prefix
		path := input[len(contentPrefix):]
		// adjust path for home directory
		if strings.HasPrefix(path, "~") {
			homeDirPath, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			log.Logger.Debug().Msgf("Home directory: %s", homeDirPath)
			path = strings.Replace(path, "~", homeDirPath, 1)
		}
		// read file content into string
		data, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(data)), nil
	}
	return input, nil
}
