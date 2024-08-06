package config

type ProviderHcloud struct {
	Token        string `yaml:"token"`
	Endpoint     string `yaml:"endpoint,omitempty"`
	PollInterval string `yaml:"poll-interval,omitempty"`
	PollFunction string `yaml:"poll-function,omitempty"`
}
