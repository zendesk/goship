package config

// Config represents goship config
type Config struct {
	LoginUsername     string `mapstructure:"username" yaml:"username"`
	UsePrivateNetwork bool   `mapstructure:"use_private_network" yaml:"use_private_network"`
	UseDNS            bool   `mapstructure:"use_dns" yaml:"use_dns"`
	CacheDirectory    string `mapstructure:"cache_directory" yaml:"cache_directory"`
	CacheFilePrefix   string `mapstructure:"cache_file_prefix" yaml:"cache_file_prefix"`
	CacheValidity     uint64 `mapstructure:"cache_validity" yaml:"cache_validity"`
	Verbose           bool   `mapstructure:"verbose" yaml:"verbose"`

	SSHBinary      string   `mapstructure:"ssh_binary" yaml:"ssh_binary"`
	SSHExtraParams []string `mapstructure:"ssh_extra_params" yaml:"ssh_extra_params"`
	ScpBinary      string   `mapstructure:"scp_binary" yaml:"scp_binary"`
	ScpExtraParams []string `mapstructure:"scp_extra_params" yaml:"scp_extra_params"`

	Providers map[string]interface{} `yaml:"providers,omitempty"`
}

// GlobalConfig holds globally accessible config var
var GlobalConfig = Config{
	SSHBinary: "/usr/bin/ssh",
	ScpBinary: "/usr/bin/scp",
}
