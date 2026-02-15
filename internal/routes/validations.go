package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

func validateTimeRequest(c fiber.Ctx, id string) (updateTimeRequest, error) {
	var req updateTimeRequest
	err := c.Bind().Body(&req)
	if err != nil {
		log.Error().Str("tag", "routes_api").Str("id", id).Err(err).Msg("Error parsing post video time request")
		return updateTimeRequest{}, fiber.NewError(fiber.StatusBadRequest, "Bad body")
	}

	if req.Time < 0 {
		log.Error().Str("tag", "routes_api").Str("id", id).Msg("Input time was negative")
		return updateTimeRequest{}, fiber.NewError(fiber.StatusBadRequest, "Bad time value")
	}

	if req.Time > 90_000 {
		log.Error().Str("tag", "routes_api").Str("id", id).Msg("Input time was too large")
		return updateTimeRequest{}, fiber.NewError(fiber.StatusBadRequest, "Bad time value")
	}
	return req, err
}
