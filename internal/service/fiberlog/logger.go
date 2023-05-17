package fiberlog

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"runtime/debug"
	"sync"
	"time"
)

type logFields struct {
	ID         string
	RemoteIP   string
	Host       string
	Method     string
	Path       string
	Protocol   string
	StatusCode int
	Latency    float64
	Error      error
	Stack      []byte
}

func (lf *logFields) MarshalZerologObject(e *zerolog.Event) {
	e.
		Str("id", lf.ID).
		Str("remote_ip", lf.RemoteIP).
		Str("host", lf.Host).
		Str("method", lf.Method).
		Str("path", lf.Path).
		Str("protocol", lf.Protocol).
		Int("status_code", lf.StatusCode).
		Float64("latency", lf.Latency).
		Str("tag", "request")

	if lf.Error != nil {
		e.Err(lf.Error)
	}

	if lf.Stack != nil {
		e.Bytes("stack", lf.Stack)
	}
}

func New(log zerolog.Logger) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		rid := c.Get(fiber.HeaderXRequestID)
		if rid == "" {
			rid = uuid.New().String()
			c.Set(fiber.HeaderXRequestID, rid)
		}

		fields := &logFields{
			ID:       rid,
			RemoteIP: c.IP(),
			Method:   c.Method(),
			Host:     c.Hostname(),
			Path:     c.Path(),
			Protocol: c.Protocol(),
		}

		// Set variables
		var (
			once       sync.Once
			errHandler fiber.ErrorHandler
		)

		// Set error handler once
		once.Do(func() {
			// override error handler
			errHandler = c.App().ErrorHandler
		})

		// Panic recover + logging
		defer func() {
			rvr := recover()
			if rvr != nil {
				err, ok := rvr.(error)
				if !ok {
					err = fmt.Errorf("%v", rvr)
				}

				fields.Error = err
				fields.Stack = debug.Stack()

				if err = errHandler(c, err); err != nil {
					_ = c.SendStatus(fiber.StatusInternalServerError)
				}

				fields.StatusCode = c.Context().Response.StatusCode()
				fields.Latency = time.Since(start).Seconds()

				log.Error().EmbedObject(fields).Msg("panic recover")
			}
		}()

		err := c.Next()
		if err != nil {
			if err := errHandler(c, err); err != nil {
				_ = c.SendStatus(fiber.StatusInternalServerError)
			}
		}

		fields.StatusCode = c.Context().Response.StatusCode()
		fields.Latency = time.Since(start).Seconds()

		switch {
		case fields.StatusCode >= 500:
			log.Error().EmbedObject(fields).Msg("server error")
		case fields.StatusCode >= 400:
			log.Error().EmbedObject(fields).Msg("client error")
		case fields.StatusCode >= 300:
			log.Info().EmbedObject(fields).Msg("redirect")
		case fields.StatusCode >= 200:
			log.Info().EmbedObject(fields).Msg("success")
		case fields.StatusCode >= 100:
			log.Info().EmbedObject(fields).Msg("informative")
		default:
			log.Warn().EmbedObject(fields).Msg("unknown status")
		}

		return err
	}
}
