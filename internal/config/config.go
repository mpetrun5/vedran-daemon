package config

import (
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
)

func InitMainConfig() {
	setDefaultValuesForMainConfig()
	viper.SetConfigName(getMainConfigName()) // name of config file (without extension)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	_ = viper.ReadInConfig()

	// Load env variables from .env
	gotenv.Load()

	if os.Getenv("ENV") != "test" {
		setupSentry()
	}
}

func setDefaultValuesForMainConfig() {
	viper.SetDefault("stats.interval", 3600)
	viper.SetDefault("log.level", "error")
}

// depending on ENV variable creates name for config file
func getMainConfigName() string {
	configFileName := "config"
	if env := os.Getenv("ENV"); env != "" {
		configFileName = configFileName + "-" + env
	}
	return configFileName
}

func setupSentry() {
	dsn := os.Getenv("SENTRY_DSN")
	sentry.Init(sentry.ClientOptions{
		Dsn:   dsn,
		Debug: false,
	})

	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)
}