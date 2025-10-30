package cmd

import (
	"app/pkg/rabbitMQ"
	"app/utils/constant"
	"log"

	"app/config"
	rest "app/internal/delivery/http"
	"app/internal/repository"
	"app/internal/usecase"
	"github.com/spf13/cobra"
)

var restCommand = &cobra.Command{
	Use: "rest",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Get()

		pgDB := NewPostgres(cfg)

		rabbitMq := rabbitMQ.NewRabbitMQImpl(cfg)
		rabbitMq.Init()
		defer rabbitMq.Close()

		err := rabbitMq.Configure(constant.GeofenceEntryEvent)
		if err != nil {
			log.Fatal(err)
		}

		app := rest.NewRest(cfg)
		repo := repository.NewRepository(pgDB)
		useCase := usecase.NewUsecase(cfg, repo, rabbitMq)

		route := &rest.Route{
			Config:     cfg,
			App:        app,
			UseCase:    useCase,
			Repository: repo,
		}

		route.RegisterRoutes()

		if err := rest.Serve(app, cfg); err != nil {
			log.Fatalf("Failed to start server: %v\n", err)
		}

	},
}
