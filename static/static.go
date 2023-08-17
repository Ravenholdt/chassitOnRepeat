package static

import (
	"chassit-on-repeat/internal/utils"
	"embed"
	"io/fs"
	"os"
)

//go:embed *
var files embed.FS

func GetFiles() fs.FS {
	// Use embedded files if not running with DEBUG
	if !utils.GetBoolEnv("DEBUG", false) {
		return files
	}

	return os.DirFS("./static")
}
