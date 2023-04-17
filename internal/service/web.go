package service

import (
	"chassit-on-repeat/internal/service/fiberlog"
	"chassit-on-repeat/internal/utils"
	"chassit-on-repeat/views"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/helmet/v2"
	"github.com/gofiber/template/html"
	"github.com/rs/zerolog/log"
	"math/rand"
	"net/http"
	"strings"
)

func createWebApp() *fiber.App {
	// Initialize standard Go html template engine
	engine := html.NewFileSystem(http.FS(views.Views), ".tmpl")

	proxy := utils.GetBoolEnv("ENABLE_PROXY", false)
	proxyHeader := ""
	if proxy {
		proxyHeader = fiber.HeaderXForwardedFor
	}
	trustedProxy := utils.GetStringEnv("TRUSTED_PROXY", "127.0.0.1")

	app := fiber.New(fiber.Config{
		Views:                   engine,
		DisableStartupMessage:   true,
		EnableTrustedProxyCheck: proxy,
		TrustedProxies:          []string{trustedProxy},
		ProxyHeader:             proxyHeader,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			msg := err.Error()
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
				msg = e.Message
			}

			if strings.HasPrefix(ctx.OriginalURL(), "/api/") {
				return ctx.Status(code).JSON(fiber.Map{
					"message": endpointErrors[rand.Intn(len(endpointErrors))],
					"path":    ctx.OriginalURL(),
					"error": fiber.Map{
						"code":    code,
						"message": msg,
					},
				})
			}

			if code == fiber.StatusNotFound {
				return ctx.Status(code).Render("views/404", fiber.Map{})
			}

			ctx.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
			return ctx.Status(code).SendString(err.Error())
		},
	})

	app.Use(requestid.New())
	app.Use(fiberlog.New(log.Logger))
	app.Use(helmet.New())
	app.Use(compress.New())
	app.Use(cors.New())

	return app
}