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
	api.Get("/", r.ApiIndex)

	// V1 api
	v1 := api.Group("/v1")
	v1.Get("/", r.ApiIndex)
	v1.Get("/stats", r.ApiStats)

	// /api/video/* endpoints
	video := v1.Group("/video")
	video.Get("/", r.ApiGetVideos)
	video.Get("/random", r.ApiRandom)

	specificVideo := video.Group("/:id")
	specificVideo.Get("/", r.ApiGetVideo)
	specificVideo.Post("/", r.ApiPostVideoTime)
	specificVideo.Post("/settings", r.ApiPostVideoSettings)

	// Index route with no
	app.Get("/random", r.ViewRandom)

	app.Get("/", r.ViewVideo)
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
