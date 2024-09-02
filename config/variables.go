// Copyright 2024 Bitshift D.O.O
// SPDX-License-Identifier: MPL-2.0

package config

import "github.com/bitshifted/liftoff/common"

const (
	envDefault = "default"
)

type ConfigVariables map[string]map[string]interface{}

// returns variables for specific environment
// TODO this should take environment name as parameters
func (cv ConfigVariables) forEnvironment() map[string]interface{} {
	return cv[envDefault]
}

func processVariables(vars map[string]interface{}) error {
	for key, val := range vars {
		success, err := setStringValue(vars, key, val)
		if err != nil {
			return err
		}
		if !success {
			success, err = setListValues(vars, key, val)
			if err != nil {
				return err
			}
		}
		if !success {
			mval, ok := val.(map[string]interface{})
			if ok {
				err := processVariables(mval)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func setStringValue(vars map[string]interface{}, key string, value interface{}) (bool, error) {
	sval, ok := value.(string)
	if ok {
		processed, err := common.ProcessStringValue(sval)
		if err != nil {
			return false, err
		}
		vars[key] = processed
		return true, nil
	}
	return false, nil
}

func setListValues(vars map[string]interface{}, key string, value interface{}) (bool, error) {
	lval, ok := value.([]interface{})
	var out []interface{}
	if ok {
		for _, item := range lval {
			sval, ok := item.(string)
			if ok {
				processed, err := common.ProcessStringValue(sval)
				if err != nil {
					return false, err
				}
				out = append(out, processed)
			} else {
				out = append(out, item)
			}
		}
		vars[key] = out
		return true, nil
	}
	return false, nil
}
