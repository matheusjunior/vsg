package repo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"matheus.com/vgs/configs"
	zapWrapper "matheus.com/vgs/internal/logger"
)

func NewPostgresConnection() *gorm.DB {
	host := configs.GetPostgresUri()
	db, err := gorm.Open(postgres.Open(host), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		zapWrapper.Logger().Fatal(err)
	}
	return db
}

func NewMongoConnection() *mongo.Client {
	opts := options.Client().ApplyURI(configs.GetMongoUri())
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		zapWrapper.Logger().Fatal(err)
	}
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		zapWrapper.Logger().Fatal(err)
	}
	return client
}
