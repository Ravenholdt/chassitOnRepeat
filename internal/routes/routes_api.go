package routes

import (
	"chassit-on-repeat/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// ApiIndex Returns that the api is running.
//
// GET /api
// GET /api/v1
func (r *Routes) ApiIndex(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).SendString("Api is running")
}

// ApiStats Returns playtime and video statistics.
//
// GET /api/v1/stats
func (r *Routes) ApiStats(c *fiber.Ctx) error {
	dbVideos, err := r.DB.GetDBVideos()
	videoCount := len(r.Files.GetVideos())

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	var totalTime int64
	for _, v := range *dbVideos {
		totalTime += utils.Val(v.Video.Time, 0)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"repeat_entries":       len(*dbVideos),
		"video_count":          videoCount,
		"total_time":           totalTime,
		"total_time_formatted": utils.FormatReadableTime(totalTime, true),
	})
}

// ApiGetVideos Returns an array of all videos.
//
// GET /api/v1/video
func (r *Routes) ApiGetVideos(c *fiber.Ctx) error {
	videos, err := r.DB.GetDBVideos()
	if err != nil {
		return err
	}

	for i, video := range *videos {
		file, e := r.Files.GetVideoFile(video.Video.ID)
		if e == nil {
			video.File = *file
			(*videos)[i] = video
		}
	}
	return c.Status(fiber.StatusOK).JSON(videos)
}

// ApiGetVideo Returns a specific video specified by the id.
//
// GET /api/v1/video/:id
func (r *Routes) ApiGetVideo(c *fiber.Ctx) error {
	id := c.Params("id")
	response, err := r.DB.GetVideoFromId(id)
	if err != nil {
		log.Error().Str("tag", "routes_api").Str("id", id).Err(err).Msg("Error getting video")
		return fiber.NewError(fiber.StatusNotFound, "Video not found")
	}

	return c.Status(fiber.StatusOK).JSON(r.GetResponse(*response))
}

// ApiPostVideoTime Updates the repeated time of the specified video.
// Post json data:
//
//	"time": integer
//
// POST /api/v1/video/:id
func (r *Routes) ApiPostVideoTime(c *fiber.Ctx) error {
	id := c.Params("id")

	var req updateVideoTimeRequest
	err := c.BodyParser(&req)
	if err != nil {
		log.Error().Str("tag", "routes_api").Str("id", id).Err(err).Msg("Error parsing post video time request")
		return fiber.NewError(fiber.StatusBadRequest, "Bad body")
	}

	v, err := r.DB.UpdateVideoPlaytime(id, int64(req.Time))
	if err != nil {
		log.Error().Str("tag", "routes_api").Str("id", id).Err(err).Msg("Error updating video time")
		return fiber.NewError(fiber.StatusBadRequest, "Error updating time")
	}

	return c.Status(fiber.StatusOK).JSON(r.GetResponse(*v))
}

// ApiPostVideoSettings Updates the start/end and safe status of the specified video.
// Post json data:
//
//	"start": integer|nil
//	"end": integer|nil
//	"safe": boolean
//
// POST /api/v1/video/:id/settings
func (r *Routes) ApiPostVideoSettings(c *fiber.Ctx) error {
	id := c.Params("id")

	var req updateVideoSettingsRequest
	err := c.BodyParser(&req)
	if err != nil {
		log.Error().Str("tag", "routes_api").Str("id", id).Err(err).Msg("Error parsing post video settings request")
		return fiber.NewError(fiber.StatusBadRequest, "Bad body")
	}

	v, err := r.DB.UpdateVideoSettings(id, req.Start, req.End, req.Safe)
	if err != nil {
		log.Error().Str("tag", "routes_api").Str("id", id).Err(err).Msg("Error updating video settings")
		return fiber.NewError(fiber.StatusBadRequest, "Error updating settings")
	}

	return c.Status(fiber.StatusOK).JSON(r.GetResponse(*v))
}

// ApiRandom Returns a random video.
// Query parms:
//
//	"safe": boolean
//
// /api/v1/video/random
func (r *Routes) ApiRandom(c *fiber.Ctx) error {
	safe := c.Request().URI().QueryArgs().Has("safe")

	video, err := r.DB.GetRandomVideo(r.Files.GetVideoIds(), safe)
	if err != nil {
		log.Error().Str("tag", "routes_api").Bool("safe", safe).Err(err).Msg("Error getting random video")
		return fiber.NewError(fiber.StatusNotFound, "Error getting random video")
	}

	id := r.Overrides.GetOverride(video.ID)
	if id != video.ID {
		// Replace with overridden video if not error happens
		fromId, err := r.DB.GetVideoFromId(id)
		if err == nil {
			video = fromId
		}
	}

	return c.Status(fiber.StatusOK).JSON(r.GetResponse(*video))
}
