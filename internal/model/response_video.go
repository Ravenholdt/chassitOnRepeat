package model

import (
	"chassit-on-repeat/internal"
	"chassit-on-repeat/internal/utils"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"sort"
)

type ResponseVideoMap map[string]ResponseVideo

type ResponseVideo struct {
	Video Video
	File  internal.VideoFile
}

func (v ResponseVideo) ToMap() fiber.Map {
	t := utils.Val(v.Video.Time, 0)
	return fiber.Map{
		"id":             v.Video.ID,
		"title":          v.File.Name,
		"start":          utils.Val(v.Video.Start, 0.0),
		"end":            utils.Val(v.Video.End, 100000.0),
		"safe":           utils.Val(v.Video.Safe, true),
		"time":           t,
		"time_formatted": utils.FormatReadableTime(t, true),
		"url":            v.File.Url,
	}
}

func (v ResponseVideo) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.ToMap())
}

func (m *ResponseVideoMap) MarshalJSON() ([]byte, error) {
	var list []ResponseVideo
	for _, video := range *m {
		list = append(list, video)
	}

	sort.Slice(list, func(i, j int) bool {
		a := list[i]
		b := list[j]
		at := utils.Val(a.Video.Time, 0)
		bt := utils.Val(b.Video.Time, 0)
		if at == 0 && bt == 0 {
			return a.File.Name > b.File.Name
		}
		return at > bt
	})

	return json.Marshal(list)
}
