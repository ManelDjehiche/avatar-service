package env

import (
	"github.com/spf13/viper"
)

type Config struct {
	HTTPPort     string
	Database     *DatabaseConfig
	AgentService *AgentServiceConfig
}

type DatabaseConfig struct {
	PSQLConfig
}

type PSQLConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	DbName   string
}

type AgentServiceConfig struct {
	Port string
	Host string
}

var Settings *Config

func init() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")

	viper.SetDefault("HTTP_PORT", "3004")

	viper.SetDefault("POSTGRES_HOST", "localhost")
	viper.SetDefault("POSTGRES_PORT", 5432)
	viper.SetDefault("POSTGRES_USERNAME", "postgres")
	viper.SetDefault("POSTGRES_PASSWORD", "manel_2021")
	viper.SetDefault("POSTGRES_DBNAME", "postgres")

	viper.SetDefault("AGENT_SERVICE_HOST", "agent-service")
	viper.SetDefault("AGENT_SERVICE_PORT", "9004")

	Settings = &Config{
		HTTPPort: viper.GetString("HTTP_PORT"),
		Database: &DatabaseConfig{
			PSQLConfig{
				Host:     viper.GetString("POSTGRES_HOST"),
				Port:     viper.GetInt("POSTGRES_PORT"),
				Username: viper.GetString("POSTGRES_USERNAME"),
				Password: viper.GetString("POSTGRES_PASSWORD"),
				DbName:   viper.GetString("POSTGRES_DBNAME"),
			},
		},
		AgentService: &AgentServiceConfig{
			Host: viper.GetString("AGENT_SERVICE_HOST"),
			Port: viper.GetString("AGENT_SERVICE_PORT"),
		},
	}
}
