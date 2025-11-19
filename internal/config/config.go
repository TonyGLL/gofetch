package config

import (
	"strings"

	"github.com/spf13/viper"
)

// Config almacena toda la configuración de la aplicación.
// Viper lee los valores desde un archivo de configuración o variables de entorno.
type Config struct {
	MongoURI   string        `mapstructure:"mongo_uri"`
	DBName     string        `mapstructure:"db_name"`
	ServerPort int           `mapstructure:"server_port"`
	Crawler    CrawlerConfig `mapstructure:"crawler"`
	Indexer    IndexerConfig `mapstructure:"indexer"`
}

// IndexerConfig stores the configuration for the indexer.
type IndexerConfig struct {
	Path string `mapstructure:"path"`
}

// CrawlerConfig almacena la configuración para el crawler.
type CrawlerConfig struct {
	URLs     []string `mapstructure:"urls"`
	MaxDepth int      `mapstructure:"max_depth"`
}

// LoadConfig lee la configuración desde un archivo o variables de entorno.
func LoadConfig() (config Config, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
