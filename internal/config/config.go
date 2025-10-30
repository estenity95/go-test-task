package config

import (
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type DBConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	Debug    bool   `mapstructure:"debug"`
}

type AppConfig struct {
	Server ServerConfig `mapstructure:"server"`
	DB     DBConfig     `mapstructure:"db"`
}

func Load() (*AppConfig, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	_ = godotenv.Load(".env")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.BindEnv("db.host", "POSTGRES_HOST")
	viper.BindEnv("db.port", "POSTGRES_PORT")
	viper.BindEnv("db.username", "POSTGRES_USER")
	viper.BindEnv("db.password", "POSTGRES_PASSWORD")
	viper.BindEnv("db.dbname", "POSTGRES_DB")
	viper.BindEnv("db.debug", "POSTGRES_DEBUG_LOG")

	var cfg AppConfig
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
