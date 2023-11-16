package routes

import (
	"chassit-on-repeat/internal/utils"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (r *Routes) ViewRandom(c *fiber.Ctx) error {
	safe := c.Request().URI().QueryArgs().Has("safe")

	video, err := r.DB.GetRandomVideo(r.Files.GetVideoIds(), safe)
	if err != nil {
		log.Error().Str("tag", "routes_views").Bool("safe", safe).Err(err).Msg("Error getting random video")
		return fiber.NewError(fiber.StatusNotFound, "error getting random video")
	}

	id := r.Overrides.GetOverride(video.Id)
	if id != video.Id {
		// Replace with overridden video if not error happens
		fromId, err := r.DB.GetVideoFromId(id)
		if err == nil {
			video = fromId
		}
	}

	randomId := "RANDOM"
	if safe {
		randomId = "RANDOM-SAFE"
	}

	totalTime := utils.FormatReadableTime(0, true)
	vid, err := r.DB.GetVideoFromId(randomId)
	if err == nil {
		t := utils.Val(vid.Time, 0)
		totalTime = utils.FormatReadableTime(t, true)
	}

	response := r.GetResponse(*video)
	return c.Status(fiber.StatusOK).Render("random", fiber.Map{
		"title":                fmt.Sprintf("(0) Chassit radio - %s", response.GetName()),
		"video":                response.ToMap(),
		"total_time_formatted": totalTime,
		"safe":                 safe,
	}, "base")
}
