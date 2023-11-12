package db

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo"
)

type DB struct {
	Client       *mongo.Client
	VideoColl    *mgm.Collection
	PlaylistColl *mgm.Collection
}
