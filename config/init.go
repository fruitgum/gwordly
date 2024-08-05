package config

import (
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
	config "gwordly/config/models"
	"io"
	"os"
)

func EnvVars(envFile string) *config.EnvConfig {

	logger := BotLogger

	validate := validator.New(validator.WithRequiredStructEnabled())

	envFileOpen, err := os.OpenFile(envFile, os.O_RDONLY, 0755)
	if err != nil {
		logger.Fatal("CONFIG: " + err.Error())
	}

	envFileRead, err := io.ReadAll(envFileOpen)
	if err != nil {
		logger.Fatal("CONFIG: " + err.Error())
	}

	var configMap config.EnvConfig

	yamlErr := yaml.Unmarshal(envFileRead, &configMap)
	if yamlErr != nil {
		logger.Fatal("CONFIG: Parse: %v", err)
	}

	err = validate.Struct(&configMap)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			logger.Error("Config: %v is not set", err.Field())
		}
		logger.Fatal("Invalid config file")
	}

	return &configMap

}
