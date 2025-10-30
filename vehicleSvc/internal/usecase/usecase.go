package usecase

import (
	"app/config"
	"app/internal/repository"
	"app/pkg/rabbitMQ"
)

type Usecase struct {
	MqttUsecase    MqttUsecaseImpl
	VehicleUsecase VehicleUsecaseImpl
}

func NewUsecase(cfg *config.MainConfig, repository repository.RepositoryImpl, rabbitMq rabbitMQ.RabbitMQImpl) Usecase {
	return Usecase{
		MqttUsecase:    NewMqttUsecase(cfg, repository, rabbitMq),
		VehicleUsecase: NewVehicleUsecase(cfg, repository, rabbitMq),
	}
}

func (u *Usecase) GetMqttUseCase() MqttUsecaseImpl       { return u.MqttUsecase }
func (u *Usecase) GetVehicleUseCase() VehicleUsecaseImpl { return u.VehicleUsecase }
