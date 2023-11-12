package model

import (
	"chassit-on-repeat/internal/utils"
	"time"

	"github.com/kamva/mgm/v3"
)

type Playlist struct {
	mgm.DefaultModel `bson:",inline"`
	Id               string   `json:"id" bson:"id"`
	Name             string   `json:"name" bson:"name"`
	Time             *int64   `json:"playtime" bson:"playtime"`
	Safe             *bool    `json:"safe" bson:"safe"`
	LastPlayed       *int64   `json:"last_played" bson:"last_played"`
	Videos           []string `json:"videos" bson:"videos"`
}

func (p *Playlist) CollectionName() string {
	return "playlists"
}

func (p *Playlist) AddTime(t int64) {
	newTime := utils.Val(p.Time, 0) + t
	p.Time = &newTime
}

func (p *Playlist) UpdateLastPlayed() {
	unix := time.Now().Unix()
	p.LastPlayed = &unix
}

func (p *Playlist) SetSafe(safe bool) {
	p.Safe = &safe
}
