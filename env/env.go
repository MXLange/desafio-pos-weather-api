package env

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

type Env struct {
	ApiPort        string `mapstructure:"API_PORT"`
	WeatherAPIKey  string `mapstructure:"WEATHER_API_KEY"`
	WeatherBaseURL string `mapstructure:"WEATHER_BASE_URL"`
	APICEPBaseURL  string `mapstructure:"API_CEP_BASE_URL"`
}

func New(ctx context.Context) (*Env, error) {

	v := viper.New()

	if err := mergeConfigFile(v, ".env"); err != nil {
		return nil, err
	}

	v.AutomaticEnv()
	if !v.IsSet("API_PORT") && v.IsSet("PORT") {
		v.Set("API_PORT", v.GetString("PORT"))
	}

	if err := validateRequiredKeys(v); err != nil {
		return nil, err
	}

	var cfg Env
	if err := v.Unmarshal(&cfg, func(dc *mapstructure.DecoderConfig) {
		dc.TagName = "mapstructure"
	}); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func mergeConfigFile(v *viper.Viper, path string) error {
	v.SetConfigFile(path)
	v.SetConfigType("env")

	if err := v.MergeInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return err
		}
	}

	return nil
}

func validateRequiredKeys(v *viper.Viper) error {
	envType := reflect.TypeOf(Env{})
	missing := make([]string, 0, envType.NumField())

	for i := range envType.NumField() {
		field := envType.Field(i)
		key := field.Tag.Get("mapstructure")
		if key == "" {
			continue
		}

		if !v.IsSet(key) {
			missing = append(missing, key)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required env vars: %s", strings.Join(missing, ", "))
	}

	return nil
}
