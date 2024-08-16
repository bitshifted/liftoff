package tfprovider

import (
	"errors"

	"github.com/bitshifted/liftoff/common"
)

type ProviderHetznerDNS struct {
	APIToken string `yaml:"api-token,omitempty"`
}

func (phd *ProviderHetznerDNS) postLoad() error {
	if phd.APIToken == "" {
		return errors.New("api token is required for hetznerdns provider")
	}
	// replace values if needed
	replacement, err := common.ProcessStringValue(phd.APIToken)
	if err != nil {
		return err
	}
	phd.APIToken = replacement
	return nil
}
