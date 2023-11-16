package service

import (
	"chassit-on-repeat/internal"
	"chassit-on-repeat/internal/db"
	"chassit-on-repeat/internal/model"
	"chassit-on-repeat/internal/routes"
	"chassit-on-repeat/internal/utils"
	"context"
	"fmt"
	"github.com/kamva/mgm/v3"
	"github.com/rs/zerolog/log"
)

func NewService(handler *internal.FileHandler, overrides *internal.Overrides) *Service {
	client := setupMongo()
	app := createWebApp()

	r := routes.Routes{
		DB: &db.DB{
			Client:       client,
			VideoColl:    mgm.Coll(&model.Video{}),
			PlaylistColl: mgm.Coll(&model.Playlist{}),
		},
		Files:     handler,
		Overrides: overrides,
	}
	r.SetupRoutes(app)

	return &Service{
		fiber:  app,
		routes: &r,
	}
}

func (s *Service) Shutdown() {
	_ = s.fiber.Shutdown()
}

func (s *Service) Start() error {
	go s.routes.Files.Watch()
	go s.routes.Overrides.Watch()

	defer func() {
		if err := s.routes.DB.Client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	log.Info().Str("tag", "service").Msgf("Listening on port: %d", utils.GetIntEnv("PORT", 8080))
	return s.fiber.Listen(fmt.Sprintf(":%d", utils.GetIntEnv("PORT", 8080)))
}
