package db

import (
	"chassit-on-repeat/internal/db/data"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type DB struct {
	Client       *mongo.Client
	db           *mongo.Database
	VideoColl    *data.Collection
	PlaylistColl *data.Collection
}

func NewDB(client *mongo.Client) *DB {
	db := client.Database("repeat")

	return &DB{
		Client:       client,
		db:           db,
		VideoColl:    &data.Collection{Collection: db.Collection("data")},
		PlaylistColl: &data.Collection{Collection: db.Collection("playlists")},
	}
}
