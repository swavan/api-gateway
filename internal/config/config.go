package config

type Server struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type Resource struct {
	Name          string `mapstructure:"name"`
	Endpoint      string `mapstructure:"endpoint"`
	Authenticated bool   `mapstructure:"authenticated"`
	Destination   string `mapstructure:"destination"`
	Active        bool   `mapstructure:"active"`
}

type configuration struct {
	Server    Server     `mapstructure:"server"`
	Resources []Resource `mapstructure:"resources"`
}

var Config configuration
