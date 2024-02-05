package config

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

func NewDatabase(config *viper.Viper, log *logrus.Logger) *mongo.Database {
	username := config.GetString("database.username")
	password := config.GetString("database.password")
	host := config.GetString("database.host")
	port := config.GetInt("database.port")
	database := config.GetString("database.name")
	databaseType := config.GetString("database.type")

	if databaseType != "nosql" {
		log.Fatalf("Database is not NoSQL: %v", databaseType)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	credentials := options.Credential{
		Username: username,
		Password: password,
	}

	client, err := mongo.Connect(ctx, options.Client().SetAuth(credentials).ApplyURI(fmt.Sprintf("mongodb://%s:%d", host, port)))
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	log.Debugf("Connected to database: %s", host)

	return client.Database(database)
}
