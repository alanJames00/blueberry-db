// database and server wide configuraton loader and manager
package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort string;
	AofEnabled bool;
	AofFilePath string;
	LogLevel string; // info, debug, error
	Password string;
}

func LoadConfig() *Config {
	
	// load config from config.toml
	viper.SetConfigFile("config.toml")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".");

	// read config file if it exist
	err := viper.ReadInConfig();
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// create and populate the config struct
	config := &Config{
		ServerPort: viper.GetString("server.port"),	
		AofEnabled: viper.GetBool("persistence.enabled"),
		AofFilePath: viper.GetString("persistence.file_path"),
		LogLevel: viper.GetString("logging.level"),
		Password: viper.GetString("security.password"),
	}
	
	return config;
}
