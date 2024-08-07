package config

import (
	"fmt"

	"github.com/bitshifted/easycloud/common"
	"github.com/bitshifted/easycloud/log"
)

type ProviderHcloud struct {
	Token        string `yaml:"token"`
	Endpoint     string `yaml:"endpoint,omitempty"`
	PollInterval string `yaml:"poll-interval,omitempty"`
	PollFunction string `yaml:"poll-function,omitempty"`
}

func (hc *ProviderHcloud) postLoad() error {
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
	}
	return nil
}
