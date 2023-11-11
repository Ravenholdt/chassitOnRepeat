package model

import (
	"chassit-on-repeat/internal/utils"
	"time"

	"github.com/kamva/mgm/v3"
)

type Playlist struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string   `json:"name" bson:"name"`
	Time             *int64   `json:"playtime" bson:"playtime"`
	Safe             *bool    `json:"safe" bson:"safe"`
	LastPlayed       *int64   `json:"lastplayed" bson:"lastplayed"`
	List             []string `json:"list" bson:"list"`
}

func (p *Playlist) CollectionName() string {
	return "data"
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
