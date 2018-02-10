package conf

import (
	"strings"

	"os"

	"bufio"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	API struct {
		Host string `mapstructure:"host" json:"host"`
		Port int    `mapstructure:"port" json:"port"`
	} `mapstructure:"api" json:"api"`

	DB struct {
		Path string `mapstructure:"path" json:"path"`
	} `mapstructure:"db" json:"db"`

	LogConf struct {
		Level string `mapstructure:"level"`
		File  string `mapstructure:"file"`
	} `mapstructure:"log_conf"`
}

// Load will construct the config from the file
func Load(configFile string) (*Config, error) {
	viper.SetConfigType("json")

	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath("./") // ./config.[json | toml]
	}

	viper.SetEnvPrefix("CLOUD_INITER")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil && !os.IsNotExist(err) {
		return nil, errors.Wrap(err, "reading configuration from files")
	}

	config := new(Config)
	if err := viper.Unmarshal(config); err != nil {
		return nil, errors.Wrap(err, "unmarshaling configuration")
	}

	if err := configureLogging(config); err != nil {
		return nil, errors.Wrap(err, "configure logging")
	}

	return validateConfig(config)
}

func configureLogging(config *Config) error {
	// always use the full timestamp
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:    true,
		DisableTimestamp: false,
	})

	// use a file if you want
	if config.LogConf.File != "" {
		f, errOpen := os.OpenFile(config.LogConf.File, os.O_RDWR|os.O_APPEND, 0660)
		if errOpen != nil {
			return errOpen
		}
		logrus.SetOutput(bufio.NewWriter(f))
		logrus.Infof("Set output file to %s", config.LogConf.File)
	}

	if config.LogConf.Level != "" {
		level, err := logrus.ParseLevel(strings.ToUpper(config.LogConf.Level))
		if err != nil {
			return err
		}
		logrus.SetLevel(level)
		logrus.Debug("Set log level to: " + logrus.GetLevel().String())
	}

	return nil
}

func validateConfig(config *Config) (*Config, error) {
	if config.API.Port == 0 && os.Getenv("PORT") != "" {
		port, err := strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			return nil, errors.Wrap(err, "formatting PORT into int")
		}

		config.API.Port = port
	}

	if config.API.Port == 0 && config.API.Host == "" {
		config.API.Port = 8080
	}

	return config, nil
}
