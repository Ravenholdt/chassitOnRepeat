package routes

import (
	"chassit-on-repeat/internal/model"
	"chassit-on-repeat/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"sort"
	"strconv"
)

func (r *Routes) MergeFileVideos(videos *model.ResponseVideoMap) {
	for _, file := range r.Files.GetVideos() {
		if v, ok := (*videos)[file.Id]; ok {
			v.File = file
			(*videos)[file.Id] = v
		} else {
			(*videos)[file.Id] = model.ResponseVideo{
				Video: model.Video{
					ID: file.Id,
				},
				File: file,
			}
		}
	}
}

func (r *Routes) ViewVideo(c *fiber.Ctx) error {
	id := c.Params("id")
	videoFile, _ := r.Files.GetVideoFile(id)

	videos, err := r.DB.GetDBVideos()
	if err != nil {
		log.Error().Str("tag", "routes_views").Str("id", id).Err(err).Msg("Error getting DB videos")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	r.MergeFileVideos(videos)

	var totalTime int64
	var history []fiber.Map
	for _, v := range *videos {
		name := v.Video.ID
		if v.File.Name != "" {
			name = v.File.Name
		}
		t := utils.Val(v.Video.Time, 0)
		h := fiber.Map{
			"url":            v.Video.ID,
			"name":           name,
			"time":           t,
			"time_formatted": utils.FormatReadableTime(t, false),
		}
		totalTime += t
		history = append(history, h)
	}

	// Sort based on time, if not played sort based on name
	sort.Slice(history, func(i, j int) bool {
		a := history[i]
		b := history[j]
		if a["time"].(int64) == 0 && b["time"].(int64) == 0 {
			return a["name"].(string) > b["name"].(string)
		}
		return a["time"].(int64) > b["time"].(int64)
	})

	// If no videoFile is found or no id provided only show list
	if videoFile == nil {
		return c.Status(fiber.StatusOK).Render("index", fiber.Map{
			"title":   "Home",
			"history": history,
		})
	}

	video := (*videos)[videoFile.Id]
	startInput := ""
	endInput := ""
	start := utils.Val(video.Video.Start, 0)
	if start > 0 {
		startInput = strconv.FormatFloat(start, 'f', -1, 64)
	}
	end := utils.Val(video.Video.End, 100000.0)
	if end < 70000 {
		endInput = strconv.FormatFloat(end, 'f', -1, 64)
	}

	return c.Status(fiber.StatusOK).Render("index", fiber.Map{
		"video":                video.ToMap(),
		"start_input":          startInput,
		"end_input":            endInput,
		"total_time":           totalTime,
		"total_time_formatted": utils.FormatReadableTime(totalTime, true),
		"history":              history,
	})
}

func (r *Routes) ViewRandom(c *fiber.Ctx) error {
	safe := c.Request().URI().QueryArgs().Has("safe")

	video, err := r.DB.GetRandomVideo(r.Files.GetVideoIds(), safe)
	if err != nil {
		log.Error().Str("tag", "routes_views").Bool("safe", safe).Err(err).Msg("Error getting random video")
		return fiber.NewError(fiber.StatusNotFound, "error getting random video")
	}

	id := r.Overrides.GetOverride(video.ID)
	if id != video.ID {
		// Replace with overridden video if not error happens
		fromId, err := r.DB.GetVideoFromId(id)
		if err == nil {
			video = fromId
		}
	}

	return c.Status(fiber.StatusOK).Render("random", fiber.Map{
		"video": r.GetResponse(*video).ToMap(),
		"safe":  safe,
	})
}
