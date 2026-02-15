package routes

import (
	"chassit-on-repeat/internal"
	"chassit-on-repeat/internal/db"
)

type Routes struct {
	DB        *db.DB
	Files     *internal.FileHandler
	Overrides *internal.Overrides
}

type updateTimeRequest struct {
	Time float64 `json:"time" validate:"required,numeric,gt=0,lt=90000" message:"required=Time must be between 0 and 90000"`
}

type updateVideoSettingsRequest struct {
	Start *float64 `json:"start"`
	End   *float64 `json:"end"`
	Safe  bool     `json:"safe"`
}
