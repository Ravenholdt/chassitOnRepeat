package service

import (
	"chassit-on-repeat/internal/routes"
	"github.com/gofiber/fiber/v2"
)

type Service struct {
	fiber  *fiber.App
	routes *routes.Routes
}

var endpointErrors = []string{
	"The endpoint is in another castle",
	"There is nothing here...",
	"Maybe ask a friend?",
	"Everything happens for a reason...",
	"This is not the endpoint you are looking for",
}
