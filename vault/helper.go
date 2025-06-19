package vault

import (
	"context"
	"os"
	"strings"

	"github.com/anggardagasta/go-sdk/zlog"
	"github.com/spf13/viper"
)

// LoadConfiguration loads configuration from vault
func LoadConfiguration(env string, bindCfg interface{}) error {
	if env == "production" || env == "staging" {
		vaultURL := os.Getenv("VAULT_ADDR")
		vaultPath := os.Getenv("VAULT_PATH")
		vaultMode := os.Getenv("VAULT_MODE")

		if vaultMode == "banzai" {
			viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
			viper.SetConfigName("config")
			viper.SetConfigType("json")
			viper.AllowEmptyEnv(true)
			viper.AutomaticEnv()
		} else {
			if err := viper.AddRemoteProvider("vault", vaultURL, vaultPath); err != nil {
				zlog.Error(context.Background(), nil, err.Error())
				return err
			}
			viper.SetConfigType("json")
			if err := viper.ReadRemoteConfig(); err != nil {
				zlog.Error(context.Background(), nil, err.Error())
				return err
			}
		}

	} else {
		// load config from file config.yaml using viper
		viper.SetConfigName("config")
		viper.SetConfigType("json")
		// viper.AddConfigPath("/src")
		viper.AddConfigPath(".")

		err := viper.ReadInConfig()
		if err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				// Config file not found; ignore error if desired
				return err
			} else {
				// Config file was found but another error was produced
				return err
			}
		}
	}

	if err := viper.Unmarshal(bindCfg); err != nil {
		return err
	}

	return nil
}
