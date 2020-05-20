package config

import (
	"log"

	"github.com/dhyaniarun1993/foody-common/datastore/redis"
	"github.com/dhyaniarun1993/foody-common/datastore/sql"
	"github.com/dhyaniarun1993/foody-common/logger"
	"github.com/dhyaniarun1993/foody-common/tracer"
	customerHttpClient "github.com/dhyaniarun1993/foody-customer-service/client/http"
	"github.com/kelseyhightower/envconfig"
)

// Configuration provides application configuration
type Configuration struct {
	Port              int    `required:"true" split_words:"true"`
	AccessTokenSecret string `required:"true" split_words:"true"`
	AccessTokenIssuer string `required:"true" split_words:"true"`
	Redis             redis.Configuration
	SQL               sql.Configuration
	Log               logger.Configuration
	Jaeger            tracer.Configuration
	Customer          customerHttpClient.Configuration
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
