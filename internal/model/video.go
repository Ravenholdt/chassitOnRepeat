package model

import (
	"chassit-on-repeat/internal/db/data"
	"chassit-on-repeat/internal/utils"
	"time"
)

type Video struct {
	data.DefaultModel `bson:",inline"`
	Id                string   `json:"id" bson:"name"`
	Start             *float64 `json:"start" bson:"start"`
	End               *float64 `json:"end" bson:"end"`
	Time              *int64   `json:"playtime" bson:"playtime"`
	Safe              *bool    `json:"safe" bson:"safe"`
	LastPlayed        *int64   `json:"lastplayed" bson:"lastplayed"`
}

func (v *Video) AddTime(t int64) {
	v.Time = new(utils.Val(v.Time, 0) + t)
}

func (v *Video) UpdateLastPlayed() {
	v.LastPlayed = new(time.Now().Unix())
}

func (v *Video) SetSafe(safe bool) {
	v.Safe = &safe
}

func NewVideo(id string) *Video {
	return &Video{
		Id:   id,
		Safe: new(true),
	}
}

func NewVideoWithTime(id string, t int64) *Video {
	v := NewVideo(id)
	v.AddTime(t)
	v.UpdateLastPlayed()
	return v
}

func NewVideoWithLoop(id string, start *float64, end *float64) *Video {
	v := NewVideo(id)
	v.Start = start
	v.End = end
	return v
}
