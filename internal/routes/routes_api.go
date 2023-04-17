package routes

import (
	"chassit-on-repeat/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"strconv"
)

func (r *Routes) ApiIndex(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).SendString("Api is running")
}

func (r *Routes) ApiStats(c *fiber.Ctx) error {
	videos, err := r.DB.GetDBVideos()

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	var totalTime int64
	for _, v := range *videos {
		totalTime += utils.Val(v.Video.Time, 0)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"video_count":          len(*videos),
		"total_time":           totalTime,
		"total_time_formatted": utils.FormatReadableTime(totalTime, true),
	})
}

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

func (r *Routes) ApiGetVideo(c *fiber.Ctx) error {
	id := c.Params("id")
	response, err := r.DB.GetVideoFromId(id)
	if err != nil {
		log.Error().Str("tag", "routes_api").Str("id", id).Err(err).Msg("Error getting video")
		return fiber.NewError(fiber.StatusNotFound, "Video not found")
	}

	return c.Status(fiber.StatusOK).JSON(r.GetResponse(*response))
}

func (r *Routes) ApiPostVideoTime(c *fiber.Ctx) error {
	id := c.Params("id")

	var req updateVideoTimeRequest
	err := c.BodyParser(&req)
	if err != nil {
		log.Error().Str("tag", "routes_api").Str("id", id).Err(err).Msg("Error parsing post video time request")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	v, err := r.DB.UpdateVideoPlaytime(id, int64(req.Time))
	if err != nil {
		log.Error().Str("tag", "routes_api").Str("id", id).Err(err).Msg("Error updating video time")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.Status(fiber.StatusOK).JSON(r.GetResponse(*v))
}

func (r *Routes) ApiPostVideoSettings(c *fiber.Ctx) error {
	id := c.Params("id")

	var req updateVideoSettingsRequest
	err := c.BodyParser(&req)
	if err != nil {
		log.Error().Str("tag", "routes_api").Str("id", id).Err(err).Msg("Error parsing post video settings request")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	v, err := r.DB.UpdateVideoSettings(id, req.Start, req.End, req.Safe)
	if err != nil {
		log.Error().Str("tag", "routes_api").Str("id", id).Err(err).Msg("Error updating video settings")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.Status(fiber.StatusOK).JSON(r.GetResponse(*v))
}

// ApiRandom Returns a random video
func (r *Routes) ApiRandom(c *fiber.Ctx) error {
	safe, err := strconv.ParseBool(c.Query("safe", "false"))
	if err != nil {
		safe = false
	}

	video, err := r.DB.GetRandomVideo(r.Files.GetVideoIds(), safe)
	if err != nil {
		log.Error().Str("tag", "routes_api").Bool("safe", safe).Err(err).Msg("Error getting random video")
		return fiber.NewError(fiber.StatusInternalServerError, "error getting random video")
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
