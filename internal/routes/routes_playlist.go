package routes

import (
	"chassit-on-repeat/internal/utils"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"sort"
)

// ViewPlaylists Renders a list of playlists
func (r *Routes) ViewPlaylists(c *fiber.Ctx) error {
	playlists, err := r.DB.GetPlaylists()
	if err != nil {
		return fiber.NewError(500, "error loading playlist")
	}

	var playlistMap []fiber.Map
	for _, playlist := range *playlists {
		url, _ := c.GetRouteURL("playlist", fiber.Map{"id": playlist.Id})
		lastPlayed := playlist.LastPlayed.Unix()
		playlistMap = append(playlistMap, fiber.Map{
			"url":            url,
			"name":           playlist.Name,
			"safe":           playlist.Safe,
			"last_played":    &lastPlayed,
			"time":           playlist.Time,
			"time_formatted": utils.FormatReadableTime(playlist.Time, false),
			"videos":         len(playlist.Videos),
		})
	}

	// Sort based on time, if not played sort based on name
	sort.Slice(playlistMap, r.sortJustPlayed(playlistMap))

	return c.Status(fiber.StatusOK).Render("playlists", fiber.Map{
		"title":     "Playlists - Chassit on Repeat",
		"playlists": playlistMap,
	}, "base")
}

// ViewPlaylist Renders a specific playlist and the list of videos
func (r *Routes) ViewPlaylist(c *fiber.Ctx) error {
	id := c.Params("id")

	playlist, err := r.DB.GetPlaylistFromId(id)
	if err != nil {
		return fiber.NewError(404, "no playlist found with id")
	}

	randomVideo, err := r.DB.GetRandomPlaylistVideo(id, r.Files.GetVideoIds())
	if err != nil {
		return fiber.NewError(404, "no random video found on playlist")
	}

	videos, err := r.DB.GetDBVideos(playlist.Videos...)
	if err != nil {
		return fiber.NewError(500, "error loading videos for playlist")
	}

	r.mergeFileVideos(videos, false)

	var playlistVideos []fiber.Map
	for _, v := range *videos {
		name := v.GetId()
		if v.GetName() != "" {
			name = v.GetName()
		}
		h := fiber.Map{
			"id":          v.GetId(),
			"name":        name,
			"safe":        v.GetSafe(),
			"last_played": v.GetLastPlayed(),
		}
		playlistVideos = append(playlistVideos, h)
	}

	// Sort based on time, if not played sort based on name
	sort.Slice(playlistVideos, r.sortJustPlayed(playlistVideos))

	data := r.GetResponse(*randomVideo)

	return c.Status(fiber.StatusOK).Render("playlist", fiber.Map{
		"title":                fmt.Sprintf("%s - %s - Chassit on Repeat", data.GetName(), playlist.Name),
		"id":                   playlist.Id,
		"name":                 playlist.Name,
		"video":                data.ToMap(),
		"videos":               playlistVideos,
		"total_time_formatted": utils.FormatReadableTime(playlist.Time, true),
	}, "base")
}
