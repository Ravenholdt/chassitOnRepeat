package views

import (
	"chassit-on-repeat/internal/utils"
	"embed"
	"net/http"
)

//go:embed *
var views embed.FS

func GetViews() http.FileSystem {
	// Use embedded files if not running with DEBUG
	if !utils.GetBoolEnv("DEBUG", false) {
		return http.FS(views)
	}

	return http.Dir("./views")
}
