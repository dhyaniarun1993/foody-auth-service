package config

import (
	"log"

	"github.com/dhyaniarun1993/foody-common/datastore/sql"
	"github.com/dhyaniarun1993/foody-common/logger"
	"github.com/dhyaniarun1993/foody-common/tracer"
	"github.com/kelseyhightower/envconfig"
)

// Configuration provides application configuration
type Configuration struct {
	Port   int `required:"true" split_words:"true"`
	SQL    sql.Configuration
	Log    logger.Configuration
	Jaeger tracer.Configuration
}

// InitConfiguration initialize the configuration
func InitConfiguration() Configuration {
	var config Configuration
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatalln(err)
	}
	return config
}
