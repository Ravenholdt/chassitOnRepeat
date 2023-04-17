package service

import (
	"context"
	"github.com/kamva/mgm/v3"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

func setupMongo() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(os.Getenv("MONGODB_URI")).
		SetServerAPIOptions(serverAPIOptions)

	err := mgm.SetDefaultConfig(nil, "repeat", clientOptions)
	if err != nil {
		panic(err)
	}

	_, client, _, _ := mgm.DefaultConfigs()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal().Str("tag", "mongo").Err(err).Msg("Can't connect to MongoDB server")
	}
	return client
}
