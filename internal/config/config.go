package config

type Config struct {
	Server   Server   `mapstructure:"server"`
	Database Database `mapstructure:"database"`
	Logging  Logging  `mapstructure:"logging"`
}

type Logging struct {
	Level string `mapstructure:"level"`
}

type Server struct {
	Port            string `mapstructure:"port"`
	ShutdownTimeout string `mapstructure:"shutdownTimeout"`
}

type Database struct {
	Dsn string `mapstructure:"dsn"`
}
