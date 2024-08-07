package common

import (
	"errors"
	"regexp"
	"strings"

	"github.com/bitshifted/easycloud/log"
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
	// check if it's file content refrence
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
