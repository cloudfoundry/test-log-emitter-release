package config

type Config struct {
	ListenPort  int               `yaml:"listen_port"`
	Loggregator LoggregatorConfig `yaml:"loggregator"`
}

type LoggregatorConfig struct {
	CA   string `yaml:"ca"`
	Cert string `yaml:"cert"`
	Key  string `yaml:"key"`
	Port int    `yaml:"port"`
}
