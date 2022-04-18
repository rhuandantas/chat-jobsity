package config

import (
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
	"strings"
	"sync"
)

type Config struct {
	AppName        string
	ServerHost     string
	MetaHost       string
	StooqClient    StooqClient
	RabbitMQClient RabbitMQClient
}

type StooqClient struct {
	UrlTempate string
}

type RabbitMQClient struct {
}

var (
	runOnce sync.Once
	config  Config
)

func GetConfig() Config {
	runOnce.Do(func() {
		cfg := initConfig()

		config = Config{
			AppName:    cfg.GetString("app-name"),
			ServerHost: cfg.GetString("server.host"),
			MetaHost:   cfg.GetString("meta.host"),
			StooqClient: StooqClient{
				UrlTempate: cfg.GetString("client.stooq.url-template"),
			},
		}
	})

	return config
}

func initConfig() viper.Viper {
	log.Info("Initializing Cardstack configuration")

	cfg := viper.New()
	cfg.SetConfigFile("json")
	cfg.SetConfigName("config")
	cfg.AddConfigPath(".")
	cfg.AddConfigPath("./config")
	cfg.AddConfigPath("./internal/config")

	if err := cfg.ReadInConfig(); err != nil {
		log.Warn("Failed to load configuration file.", err.Error())
	}

	initDefaults(cfg)

	return *cfg
}

func initDefaults(config *viper.Viper) {
	config.SetDefault("app-name", "chat-jobsity")

	config.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))

	config.AutomaticEnv()
}
