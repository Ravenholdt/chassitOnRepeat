package routes

import (
	"chassit-on-repeat/internal/db"
	"chassit-on-repeat/internal/model"
	"chassit-on-repeat/internal/utils"
	"errors"
	"fmt"
	"sort"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (r *Routes) MergeFileVideos(videos *db.ResponseDataMap) {
	for _, file := range r.Files.GetVideos() {
		if v, ok := (*videos)[file.Id]; ok {
			v.File = file
			(*videos)[file.Id] = v
		} else {
			(*videos)[file.Id] = db.CreateVideoData(model.Video{Id: file.Id}, &file)
		}
	}
}

func sortJustPlayed(history []fiber.Map) func(i int, j int) bool {
	return func(i, j int) bool {
		aVal := history[i]["last_played"].(*int64)
		bVal := history[j]["last_played"].(*int64)

		// Sort alphabetically if there is no last played
		if aVal == nil && bVal == nil {
			return history[i]["name"].(string) < history[j]["name"].(string)
		}

		if aVal == nil {
			return false
		}
		if bVal == nil {
			return true
		}
		return *aVal > *bVal
	}
}

func sortTime(history []fiber.Map) func(i int, j int) bool {
	return func(i, j int) bool {
		aVal := history[i]["time"].(int64)
		bVal := history[j]["time"].(int64)

		// Sort alphabetically if there is no time
		if aVal == 0 && bVal == 0 {
			return history[i]["name"].(string) < history[j]["name"].(string)
		}
		return aVal > bVal
	}
}

// ViewFeelingLucky Redirects to a random video
func (r *Routes) ViewFeelingLucky(c *fiber.Ctx) error {
	video, err := r.DB.GetRandomVideo(r.Files.GetVideoIds(), false)
	if err != nil {
		log.Error().Str("tag", "routes_views").Err(err).Msg("Error getting random video")
		return fiber.NewError(fiber.StatusNotFound, "Error getting random video")
	}

	return c.RedirectToRoute("video", fiber.Map{"id": video.Id}, 302)
}

// ViewLastVideos Renders a list of last videos and handle video view
func (r *Routes) ViewLastVideos(c *fiber.Ctx) error {
	history, videos, totalTime, err := r.getHistory("")
	if err != nil {
		return err
	}

	// Sort based on time, if not played sort based on name
	sort.Slice(history, sortJustPlayed(history))

	return r.renderVideoView(c, videos, history, totalTime, "video-top-list")
}

// ViewTopVideos Renders a list of top videos and handle video view
func (r *Routes) ViewTopVideos(c *fiber.Ctx) error {
	history, videos, totalTime, err := r.getHistory("")
	if err != nil {
		return err
	}

	// Sort based on time, if not played sort based on name
	sort.Slice(history, sortTime(history))

	return r.renderVideoView(c, videos, history, totalTime, "video-list")
}

func (r *Routes) renderVideoView(c *fiber.Ctx, videos *db.ResponseDataMap, history []fiber.Map, totalTime int64, timeRoute string) error {
	id := c.Params("id")
	videoFile, _ := r.Files.GetVideoFile(id)

	timeUrl, err := c.GetRouteURL(timeRoute, fiber.Map{})
	if err != nil {
		timeUrl = "/"
	}

	// If no video is found or no video is selected only show history
	if videoFile == nil {
		return c.Status(fiber.StatusOK).Render("index", fiber.Map{
			"title":                "Home",
			"time_link":            timeUrl,
			"total_time":           totalTime,
			"total_time_formatted": utils.FormatReadableTime(totalTime, true),
			"history":              history,
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
		"time_link":            timeUrl,
		"start_input":          startInput,
		"end_input":            endInput,
		"total_time":           totalTime,
		"total_time_formatted": utils.FormatReadableTime(totalTime, true),
		"history":              history,
	})
}

func (r *Routes) getHistory(urlPrefix string) ([]fiber.Map, *db.ResponseDataMap, int64, error) {
	data, err := r.DB.GetDBVideos()
	if err != nil {
		log.Error().Str("tag", "routes_views").Err(err).Msg("Error getting DB data")
		return nil, nil, 0, errors.New("error getting data")
	}

	r.MergeFileVideos(data)

	playlists, err := r.DB.GetPlaylists()
	if err != nil {
		log.Error().Str("tag", "routes_views").Err(err).Msg("Error getting DB playlists")
		return nil, nil, 0, errors.New("error getting playlists")
	}
	for _, playlist := range *playlists {
		(*data)[playlist.Id] = db.CreatePlaylistData(playlist)
	}

	var totalTime int64
	var history []fiber.Map
	for _, v := range *data {
		name := v.GetId()
		if v.GetName() != "" {
			name = v.GetName()
		}
		t := v.GetTime()
		h := fiber.Map{
			"url":            fmt.Sprintf("%s%s%s", urlPrefix, v.GetPrefix(), v.GetId()),
			"name":           name,
			"tag":            v.GetTag(),
			"time":           t,
			"safe":           v.GetSafe(),
			"last_played":    v.GetLastPlayed(),
			"time_formatted": utils.FormatReadableTime(t, false),
		}
		totalTime += t
		history = append(history, h)
	}

	return history, data, totalTime, err
}

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

	return c.Status(fiber.StatusOK).Render("random", fiber.Map{
		"video":                r.GetResponse(*video).ToMap(),
		"total_time_formatted": totalTime,
		"safe":                 safe,
	})
}
