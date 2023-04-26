package routes

import (
	"chassit-on-repeat/internal"
	"chassit-on-repeat/internal/model"
	"chassit-on-repeat/internal/utils"
	"chassit-on-repeat/static"
	_ "embed"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"net/http"
)

func (r *Routes) SetupRoutes(app *fiber.App) {
	// Provide static files from embedded filesystem
	app.Use("/static", filesystem.New(filesystem.Config{
		Root: http.FS(static.Files),
	}))
	app.Get("/favicon.ico", func(ctx *fiber.Ctx) error {
		icon, _ := static.Files.ReadFile("favicon.ico")
		return ctx.Status(200).Send(icon)
	})

	// Provide files from provided path on the server
	app.Static("/files", utils.GetStringEnv("FILES_PATH", "/files"), fiber.Static{
		Compress: true,
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

	// /api/video/* endpoints
	video := v1.Group("/video")
	// Returns an array of all videos
	video.Get("/", r.ApiGetVideos)
	// Returns a random video
	video.Get("/random", r.ApiRandom)

	// /api/video/:id/* endpoints
	specificVideo := video.Group("/:id")
	// Returns a specific video specified by the id
	specificVideo.Get("/", r.ApiGetVideo)
	// Updates the repeated time of the specified video.
	specificVideo.Post("/", r.ApiPostVideoTime)
	// Updates the start/end and safe status of the specified video.
	specificVideo.Post("/settings", r.ApiPostVideoSettings)

	// Routes serving html content

	// Serves the random video page
	app.Get("/random", r.ViewRandom)
	// Serves a list of all videos available to repeat
	app.Get("/", r.ViewVideo)
	// Serves the specified video and a list of available videos
	app.Get("/:id", r.ViewVideo)
}

func (r *Routes) GetResponse(v model.Video) *model.ResponseVideo {
	file, _ := r.Files.GetVideoFile(v.ID)
	res := &model.ResponseVideo{
		Video: v,
		File:  utils.Val(file, internal.VideoFile{}),
	}
	return res
}
