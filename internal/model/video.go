package model

import (
	"chassit-on-repeat/internal/utils"
	"github.com/kamva/mgm/v3"
	"time"
)

type Video struct {
	mgm.DefaultModel `bson:",inline"`
	Id               string   `json:"id" bson:"name"`
	Start            *float64 `json:"start" bson:"start"`
	End              *float64 `json:"end" bson:"end"`
	Time             *int64   `json:"playtime" bson:"playtime"`
	Safe             *bool    `json:"safe" bson:"safe"`
	LastPlayed       *int64   `json:"lastplayed" bson:"lastplayed"`
}

func (v *Video) CollectionName() string {
	return "data"
}

func (v *Video) AddTime(t int64) {
	newTime := utils.Val(v.Time, 0) + t
	v.Time = &newTime
}

func (v *Video) UpdateLastPlayed() {
	unix := time.Now().Unix()
	v.LastPlayed = &unix
}

func (v *Video) SetSafe(safe bool) {
	v.Safe = &safe
}

func NewVideo(id string) *Video {
	v := &Video{
		Id: id,
	}
	v.SetSafe(true)
	return v
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
