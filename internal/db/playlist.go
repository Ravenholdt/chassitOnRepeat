package db

import (
	"chassit-on-repeat/internal/model"
	"chassit-on-repeat/internal/utils"
	"errors"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math/rand"
)

func (d *DB) updatePlaylistJsonData(playlist *model.Playlist) (*model.Playlist, error) {
	safe, err := d.GetPlaylistSafe(*playlist)
	if err != nil {
		return nil, errors.New("error checking if safe: " + err.Error())
	}
	playlist.Safe = safe
	playlist.TimeFormatted = utils.FormatReadableTime(playlist.Time, true)
	return playlist, nil
}

func (d *DB) GetPlaylists() (*[]model.Playlist, error) {
	var result []model.Playlist
	err := d.PlaylistColl.SimpleFind(&result, bson.M{}, options.Find().SetSort(bson.D{{"playtime", -1}}))
	if err != nil {
		return nil, errors.New("error data loading from database: " + err.Error())
	}
	var playlists []model.Playlist

	for _, playlist := range result {
		p, err := d.updatePlaylistJsonData(&playlist)
		if err != nil {
			return nil, errors.New("error checking if safe: " + err.Error())
		}
		playlists = append(playlists, *p)
	}
	return &playlists, nil
}

func (d *DB) GetPlaylistFromId(id string) (*model.Playlist, error) {
	result := &model.Playlist{}
	err := d.PlaylistColl.First(bson.D{{"id", id}}, result)
	if err != nil {
		return nil, errors.New("error loading data from database: " + err.Error())
	}
	return d.updatePlaylistJsonData(result)
}

func (d *DB) UpdatePlaylistPlaytime(id string, t int64) (*model.Playlist, error) {
	playlist, err := d.GetPlaylistFromId(id)

	// If no playlist is found create one
	if err != nil {
		return nil, errors.New("error no playlist found: " + err.Error())
	}

	playlist.AddTime(t)
	playlist.UpdateLastPlayed()
	err = d.PlaylistColl.Update(playlist)
	if err != nil {
		return nil, errors.New("error updating playlist time: " + err.Error())
	}

	return playlist, nil
}

func (d *DB) GetPlaylistSafe(playlist model.Playlist) (bool, error) {
	count, err := d.VideoColl.CountDocuments(mgm.Ctx(), bson.M{
		"safe": false,
		"name": bson.M{"$in": playlist.Videos},
	})
	if err != nil {
		return false, errors.New("error loading data from database: " + err.Error())
	}
	return count == 0, nil
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
		{"name": bson.M{"$in": playlist.Videos}},
		{"name": bson.M{"$in": videoIds}},
	}}

	var dbVideos []model.Video
	err = d.VideoColl.SimpleFind(&dbVideos, filter)
	if err != nil {
		return nil, errors.New("no random video found: " + err.Error())
	}
	randVid := dbVideos[rand.Intn(len(dbVideos))]
	return &randVid, nil
}
