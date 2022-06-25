package config

import "github.com/spf13/viper"

type Config struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBPort     string `mapstructure:"DB_PORT"`

	GoogleClientID       string `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret   string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	FacebookClientID     string `mapstructure:"FACEBOOK_CLIENT_ID"`
	FacebookClientSecret string `mapstructure:"FACEBOOK_CLIENT_SECRET"`
	TwitterClientID      string `mapstructure:"TWITTER_CLIENT_ID"`
	TwitterClientSecret  string `mapstructure:"TWITTER_CLIENT_SECRET"`
}

var Cfg *Config

func LoadConfig(path string) (err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&Cfg)
	return
}
