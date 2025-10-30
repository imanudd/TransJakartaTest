package Consumer

import (
	"app/internal/repository"
	"app/internal/usecase"
	"context"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"os"
	"os/signal"
	"syscall"
)

type MqttConsumerImpl interface {
	RegisterHandler()
	Shutdown()
}

type MqttConsumer struct {
	client  mqtt.Client
	usecase usecase.Usecase
	repo    repository.RepositoryImpl
}

func NewMqttConsumer(client *mqtt.Client, usecase usecase.Usecase, repo repository.RepositoryImpl) MqttConsumerImpl {
	return &MqttConsumer{
		client:  *client,
		usecase: usecase,
		repo:    repo,
	}
}

func (m *MqttConsumer) RegisterHandler() {
	vehicles, err := m.repo.GetVehicleRepo().Get(context.Background())
	if err != nil {
		return
	}

	for _, vehicle := range vehicles {
		m.client.Subscribe(fmt.Sprintf("/fleet/vehicle/%s/location", vehicle.VehicleId), 0, m.usecase.MqttUsecase.GetVehicleLocation)
	}

}

func (m *MqttConsumer) Shutdown() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	fmt.Println("Disconnecting...")
	m.client.Disconnect(250)
}
