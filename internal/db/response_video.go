package db

import (
	"chassit-on-repeat/internal"
	"chassit-on-repeat/internal/model"
	"chassit-on-repeat/internal/utils"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"sort"
)

type ResponseDataMap map[string]ResponseData

type ResponseData struct {
	Video    *model.Video
	Playlist *model.Playlist
	File     internal.VideoFile
}

func CreateVideoData(v model.Video, file *internal.VideoFile) ResponseData {
	return ResponseData{
		Video: &v,
		File:  utils.Val(file, internal.VideoFile{}),
	}
}

func CreatePlaylistData(p model.Playlist) ResponseData {
	return ResponseData{
		Playlist: &p,
	}
}

func (r *ResponseData) GetId() string {
	if r.Video != nil {
		return r.Video.Id
	} else {
		return r.Playlist.Id
	}
}

func (r *ResponseData) GetName() string {
	if r.Video != nil {
		return r.File.Name
	} else {
		return r.Playlist.Name
	}
}

func (r *ResponseData) GetTime() int64 {
	if r.Video != nil {
		return utils.Val(r.Video.Time, 0)
	} else {
		return r.Playlist.Time
	}
}

func (r *ResponseData) GetSafe() bool {
	if r.Video != nil {
		return utils.Val(r.Video.Safe, true)
	} else {
		return r.Playlist.Safe
	}
}

func (r *ResponseData) GetLastPlayed() *int64 {
	if r.Video != nil {
		return r.Video.LastPlayed
	} else {
		t := r.Playlist.LastPlayed.Unix()
		return &t
	}
}

func (r *ResponseData) GetPrefix() string {
	if r.Video != nil {
		return ""
	} else {
		return "playlist/"
	}
}

func (r *ResponseData) ToMap() fiber.Map {
	t := r.GetTime()
	result := fiber.Map{
		"id":             r.GetId(),
		"title":          r.GetName(),
		"start":          0.0,
		"end":            100000.0,
		"safe":           r.GetSafe(),
		"time":           t,
		"time_formatted": utils.FormatReadableTime(t, true),
		"url":            "",
	}

	if r.Video != nil {
		result["start"] = utils.Val(r.Video.Start, 0.0)
		result["end"] = utils.Val(r.Video.Start, 100000.0)
		result["url"] = r.File.Url
	}

	return result
}

func (r *ResponseData) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.ToMap())
}

func (m *ResponseDataMap) MarshalJSON() ([]byte, error) {
	var list []ResponseData
	for _, data := range *m {
		list = append(list, data)
	}

	sort.Slice(list, func(i, j int) bool {
		a := list[i]
		b := list[j]
		at := a.GetTime()
		bt := b.GetTime()
		if at == 0 && bt == 0 {
			return a.GetName() > b.GetName()
		}
		return at > bt
	})

	return json.Marshal(list)
}
