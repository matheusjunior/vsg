package configs

import (
	"os"

	"matheus.com/vgs/internal/logger"
)

var (
	_postgresUri = "postgres://admin:admin@localhost:5432/postgres"
	_mongoUri    = "mongodb://root:example@localhost:27017"
)

func init() {
	postgresUri, found := os.LookupEnv("POSTGRES_URI")
	if !found {
		logger.Logger().Warn("POSTGRES_URI env var missing. Using localhost connection")
	} else {
		_postgresUri = postgresUri
	}
	mongoUri, found := os.LookupEnv("MONGO_URI")
	if !found {
		logger.Logger().Warn("MONGO_URI env var missing. Using localhost connection")
	} else {
		_mongoUri = mongoUri
	}
}

func GetPostgresUri() string {
	return _postgresUri
}

func GetMongoUri() string {
	return _mongoUri
}
