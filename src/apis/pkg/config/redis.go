package config

type Redis struct {
	Host     string `mapstructure:"REDIS_HOST"`
	Port     string `mapstructure:"REDIS_PORT"`
	Database int    `mapstructure:"REDIS_DATABASE"`
	Password string `mapstructure:"REDIS_PASSWORD"`
}
