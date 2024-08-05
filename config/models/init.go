package config

type EnvLogToFile struct {
	Enabled   bool   `yaml:"enabled" validate:"boolean"`
	Directory string `yaml:"directory" validate:"omitempty,required_if=Enabled true"`
}

type EnvConfig struct {
	LogLevel  string       `yaml:"logLevel"`
	LogToFile EnvLogToFile `yaml:"logToFile"`
}
