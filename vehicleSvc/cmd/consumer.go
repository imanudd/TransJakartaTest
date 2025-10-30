package cmd

import (
	"app/config"
	"app/internal/delivery/Consumer"
	"app/internal/repository"
	"app/internal/usecase"
	"app/pkg/rabbitMQ"
	"app/utils/constant"
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

var consumerCmd = &cobra.Command{
	Use:   "consumer",
	Short: "Consumer commands",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Get()
		mqttClient := InitMosquitto(cfg)

		pgDB := NewPostgres(cfg)

		rabbitMq := rabbitMQ.NewRabbitMQImpl(cfg)
		rabbitMq.Init()
		defer rabbitMq.Close()

		err := rabbitMq.Configure(constant.GeofenceEntryEvent)
		if err != nil {
			log.Fatal(err)
		}

		repo := repository.NewRepository(pgDB)
		useCase := usecase.NewUsecase(cfg, repo, rabbitMq)
		mqttConsumer := Consumer.NewMqttConsumer(mqttClient, useCase, repo)
		mqttConsumer.RegisterHandler()

		fmt.Println("MQTT consumer started. Waiting for messages...")
		mqttConsumer.Shutdown()

	},
}
