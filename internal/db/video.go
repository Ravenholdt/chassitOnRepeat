package db

import (
	"chassit-on-repeat/internal/model"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math/rand"
)

func (d *DB) GetDBVideos(videoIdFilter ...string) (*ResponseDataMap, error) {
	var result []model.Video
	filter := bson.M{}
	if len(videoIdFilter) > 0 {
		filter = bson.M{"name": bson.M{"$in": videoIdFilter}}
	}

	err := d.VideoColl.SimpleFind(&result, filter, options.Find().SetSort(bson.D{{"playtime", -1}}))
	if err != nil {
		return nil, errors.New("error data loading from database: " + err.Error())
	}

	videos := ResponseDataMap{}
	for _, v := range result {
		videos[v.Id] = CreateVideoData(v, nil)
	}
	return &videos, nil
}

func (d *DB) GetVideoFromId(id string) (*model.Video, error) {
	result := &model.Video{}
	err := d.VideoColl.First(bson.D{{"name", id}}, result)
	if err != nil {
		return nil, errors.New("error loading data from database: " + err.Error())
	}

	return result, nil
}

func (d *DB) UpdateVideoPlaytime(id string, t int64) (*model.Video, error) {
	video, err := d.GetVideoFromId(id)

	// If no video is found create one
	if err != nil {
		video = model.NewVideoWithTime(id, t)
		return video, d.VideoColl.Create(video)
	}

	video.AddTime(t)
	video.UpdateLastPlayed()
	err = d.VideoColl.Update(video)
	if err != nil {
		return nil, errors.New("error updating video time: " + err.Error())
	}

	return video, nil
}

func (d *DB) UpdateVideoSettings(id string, start *float64, end *float64, safe bool) (*model.Video, error) {
	video, err := d.GetVideoFromId(id)

	// If no video is found create one
	if err != nil {
		video = model.NewVideoWithLoop(id, start, end)
		return video, d.VideoColl.Create(video)
	}

	video.Start = start
	video.End = end
	video.SetSafe(safe)

	err = d.VideoColl.Update(video)
	if err != nil {
		return nil, errors.New("error updating video time: " + err.Error())
	}

	return video, nil
}

func (d *DB) GetRandomVideo(ids []string, safe bool) (*model.Video, error) {
	if len(ids) <= 0 {
		return nil, errors.New("no videos to randomize")
	}
	filter := bson.M{"name": bson.M{"$in": ids}}
	if safe {
		filter["safe"] = bson.M{"$ne": false}
	}

	var dbVideos []model.Video
	err := d.VideoColl.SimpleFind(&dbVideos, filter)
	if err != nil {
		return nil, errors.New("no random video found: " + err.Error())
	}
	randVid := dbVideos[rand.Intn(len(dbVideos))]
	return &randVid, nil
}
