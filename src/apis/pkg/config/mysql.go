package config

type MySQL struct {
	Host     string `mapstructure:"MYSQL_HOST"`
	Port     string `mapstructure:"MYSQL_PORT"`
	Database string `mapstructure:"MYSQL_DATABASE"`
	Username string `mapstructure:"MYSQL_USERNAME"`
	Password string `mapstructure:"MYSQL_PASSWORD"`
}
