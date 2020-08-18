package config

// Config represents goship config
type Config struct {
	UseEC2Connect     bool   `mapstructure:"use_ec2_connect" yaml:"use_ec2_connect"`
	EC2ConnectKeyPath string `mapstructure:"ec2_connect_key_path" yaml:"ec2_connect_key_path"`
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
