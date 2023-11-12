package db

import (
	"chassit-on-repeat/internal/model"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math/rand"
)

func (d *DB) GetPlaylists() (*[]model.Playlist, error) {
	var result []model.Playlist
	err := d.PlaylistColl.SimpleFind(&result, bson.M{}, options.Find().SetSort(bson.D{{"playtime", -1}}))
	if err != nil {
		return nil, errors.New("error data loading from database: " + err.Error())
	}
	return &result, nil
}

func (d *DB) GetPlaylistFromId(id string) (*model.Playlist, error) {
	result := &model.Playlist{}
	err := d.VideoColl.First(bson.D{{"id", id}}, result)
	if err != nil {
		return nil, errors.New("error loading data from database: " + err.Error())
	}

	return result, nil
}

func (d *DB) GetRandomPlaylistVideo(id string, videoIds []string) (*model.Video, error) {
	if len(videoIds) <= 0 {
		return nil, errors.New("no videos to randomize")
	}
	playlist, err := d.GetPlaylistFromId(id)
	if err != nil {
		return nil, errors.New("no playlist found: " + err.Error())
	}

	filter := bson.M{"$and": []bson.M{
		{"name": bson.M{"$in": videoIds}},
		{"name": bson.M{"$in": playlist.Videos}},
	}}

	var dbVideos []model.Video
	err = d.VideoColl.SimpleFind(&dbVideos, filter)
	if err != nil {
		return nil, errors.New("no random video found: " + err.Error())
	}
	randVid := dbVideos[rand.Intn(len(dbVideos))]
	return &randVid, nil
}
