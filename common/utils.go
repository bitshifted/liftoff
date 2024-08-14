package common

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/bitshifted/liftoff/log"
)

type ValueType int8

const (
	envVarRegex                 = "\\$(?:(\\w+)|\\{(\\w+)\\})"
	contentPrefix               = "filecontent:"
	EnvVariableString ValueType = iota
	FileContentString
	PlainString
)

func ValueTypeFromString(input string) ValueType {
	match, err := regexp.MatchString(envVarRegex, input)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to match regex")
		return PlainString
	}
	if match {
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

func ExtractEnvVarName(input string) (string, error) {
	match, err := regexp.MatchString(envVarRegex, input)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Failed to match regex")
		return "", err
	}
	if !match {
		err := errors.New("invalid environment variable reference")
		log.Logger.Error().Err(err).Msgf("Invalid environment variable reference: %s", input)
		return "", err
	}
	if strings.HasPrefix(input, "${") && strings.HasSuffix(input, "}") {
		return input[2 : len(input)-1], nil
	} else if strings.HasPrefix(input, "$") {
		return input[1:], nil
	}
	return "", errors.New("failed to extract variable name")
}

func ProcessStringValue(input string) (string, error) {
	valueType := ValueTypeFromString(input)
	switch valueType {
	case EnvVariableString:
		// extract variable name and check if it is set
		varName, err := ExtractEnvVarName(input)
		if err != nil {
			return "", err
		}
		if !IsEnvVariableSet(varName) {
			err := fmt.Errorf("referenced environment variable %s is not set", varName)
			log.Logger.Error().Err(err).Msg("Environment variable is not set")
			return "", err
		} else {
			return "", nil
		}
	case FileContentString:
		// strip prefix
		path := input[len(contentPrefix):]
		// read file content into string
		data, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}
	return input, nil
}
