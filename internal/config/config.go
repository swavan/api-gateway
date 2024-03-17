package config

type Server struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type Service struct {
	Name        string `mapstructure:"name"`
	Endpoint    string `mapstructure:"endpoint"`
	Destination string `mapstructure:"destination"`
}

type configuration struct {
	Server   Server    `mapstructure:"server"`
	Services []Service `mapstructure:"services"`
}

var Config configuration
