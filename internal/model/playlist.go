package model

import (
	"time"

	"github.com/kamva/mgm/v3"
)

type Playlist struct {
	mgm.DefaultModel `bson:",inline"`
	Id               string    `json:"id" bson:"id"`
	Name             string    `json:"name" bson:"name"`
	Time             int64     `json:"playtime" bson:"playtime"`
	LastPlayed       time.Time `json:"last_played" bson:"last_played"`
	Videos           []string  `json:"videos" bson:"videos"`
	Safe             bool      `json:"safe" bson:"-"`
	TimeFormatted    string    `json:"time_formatted" bson:"-"`
}

func (p *Playlist) CollectionName() string {
	return "playlists"
}

func (p *Playlist) AddTime(t int64) {
	p.Time = p.Time + t
}

func (p *Playlist) UpdateLastPlayed() {
	p.LastPlayed = time.Now().UTC()
}
