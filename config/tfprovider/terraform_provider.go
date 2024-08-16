package tfprovider

type TerraformProviders struct {
	Hcloud     *ProviderHcloud     `yaml:"hcloud,omitempty"`
	HetznerDNS *ProviderHetznerDNS `yaml:"hetznerdns,omitempty"`
}

func (tp *TerraformProviders) PostLoad() error {
	if tp.Hcloud != nil {
		err := tp.Hcloud.postLoad()
		if err != nil {
			return err
		}
	}
	if tp.HetznerDNS != nil {
		err := tp.HetznerDNS.postLoad()
		if err != nil {
			return err
		}
	}
	return nil
}
