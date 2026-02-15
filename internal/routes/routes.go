package routes

import (
	"chassit-on-repeat/internal/db"
	"chassit-on-repeat/internal/model"
	"chassit-on-repeat/internal/utils"
	"chassit-on-repeat/static"
	_ "embed"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	fiberStatic "github.com/gofiber/fiber/v3/middleware/static"
)

func (r *Routes) SetupRoutes(app *fiber.App) {
	// Provide static files from embedded filesystem
	app.Get("/static/*", fiberStatic.New(
		"",
		fiberStatic.Config{
			FS: static.GetFiles(),
		},
	))

	// Provide files from a provided path on the server
	app.Get("/files/*", fiberStatic.New(
		utils.GetStringEnv("FILES_PATH", "/files"),
		fiberStatic.Config{
			Compress:  true,
			ByteRange: true,
			MaxAge:    86400, // 1 day
		}))

	limit := limiter.New(limiter.Config{
		Max:                1,
		Expiration:         1 * time.Second,
		SkipFailedRequests: true,
		KeyGenerator: func(ctx fiber.Ctx) string {
			return ctx.IP() + ctx.Path()
		},
	})

	// /api/* general endpoints
	api := app.Group("/api")
	// Root of api
	api.Get("/", r.ApiIndex)

	// V1 api
	v1 := api.Group("/v1")
	// Root of V1 api
	v1.Get("/", r.ApiIndex)
	// Returns playtime and video statistics
	v1.Get("/stats", r.ApiStats)

	// /api/v1/video/* endpoints
	video := v1.Group("/video")
	// Returns an array of all videos
	video.Get("/", r.ApiGetVideos)
	// Returns a random video
	video.Get("/random", r.ApiVideoRandom)

	// /api/v1/video/:id/* endpoints
	specificVideo := video.Group("/:id")
	// Returns a specific video specified by the id
	specificVideo.Get("/", r.ApiGetVideo)
	// Updates the repeated time of the specified video.
	specificVideo.Post("/", limit, r.ApiPostVideoTime)
	// Updates the start/end and safe status of the specified video.
	specificVideo.Put("/settings", r.ApiPostVideoSettings)
	// deprecated
	specificVideo.Post("/settings", r.ApiPostVideoSettings)

	// /api/v1/playlist/* endpoints
	playlist := v1.Group("/playlist")
	// Returns an array of all playlists
	playlist.Get("/", r.ApiGetPlaylists)

	// /api/v1/playlist/:id/* endpoints
	specificPlaylist := playlist.Group("/:id")
	// Returns the details of the playlist
	specificPlaylist.Get("/", r.ApiGetPlaylist)
	// Updates the repeated time of the specified playlist.
	specificPlaylist.Post("/", r.ApiPostPlaylistTime)
	// Returns a random video from the playlist
	specificPlaylist.Get("/random", r.ApiPlaylistRandom)

	// Routes serving html content

	// Serves the random video page
	app.Get("/random", r.ViewRandom).Name("random")

	// Redirects /random-safe to the real url
	app.Get("/random-safe", func(ctx fiber.Ctx) error {
		return ctx.Redirect().Route("random", fiber.RedirectConfig{
			Queries: map[string]string{
				"safe": "true",
			},
		})
	})

	// Redirects to a random video
	app.Get("/im-feeling-lucky", r.ViewFeelingLucky)

	// Serves a list of all playlists
	app.Get("/playlist", r.ViewPlaylists).Name("playlists")
	// Serves the specified playlist and a list of all videos in the playlist
	app.Get("/playlist/:id", r.ViewPlaylist).Name("playlist")

	// Serves a list of all videos available to repeat in top played order
	app.Get("/top", r.ViewTopVideos).Name("video-top-list")
	// Serves the specified video and a list of available videos in top played order
	app.Get("/top/:id", r.ViewTopVideos).Name("video-top")

	// Serves a list of all videos available to repeat in last played order
	app.Get("/", r.ViewLastVideos).Name("video-list")
	// Serves the specified video and a list of available videos in
	app.Get("/:id", r.ViewLastVideos).Name("video")
}

func (r *Routes) NotImplemented(ctx fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusNotImplemented)
}

func (r *Routes) GetResponse(v model.Video) *db.ResponseData {
	file, _ := r.Files.GetVideoFile(v.Id)
	res := db.CreateVideoData(v, file)
	return &res
}
