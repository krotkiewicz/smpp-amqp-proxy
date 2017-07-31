package config

import (
	"github.com/caarlos0/env"
)

var Config config

type config struct {
	LoggingFormat string        `env:"LOGGING_FORMAT"`
}


func Init () () {
	Config = config{}
	env.Parse(&Config)
}
