package usecase

import (
	"app/config"
	"app/internal/domain"
	"app/internal/repository"
	"app/pkg/rabbitMQ"
	"app/pkg/validator"
	"app/utils/constant"
	"app/utils/distance"
	"context"
	"encoding/json"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
)

type MqttUsecaseImpl interface {
	GetVehicleLocation(client MQTT.Client, msg MQTT.Message)
}

type MqttUsecase struct {
	cfg  *config.MainConfig
	repo repository.RepositoryImpl
	rmq  rabbitMQ.RabbitMQImpl
}

func NewMqttUsecase(cfg *config.MainConfig, repo repository.RepositoryImpl, rmq rabbitMQ.RabbitMQImpl) MqttUsecaseImpl {
	return &MqttUsecase{
		cfg:  cfg,
		repo: repo,
		rmq:  rmq,
	}
}

func (u *MqttUsecase) GetVehicleLocation(client MQTT.Client, msg MQTT.Message) {
	var payload domain.GetVehicleLocationPayload

	err := json.Unmarshal(msg.Payload(), &payload)
	if err != nil {
		fmt.Println("error unmarshal :", err)
	}

	if err = validator.ValidateStruct(&payload); err != nil {
		log.Println("error validator :", err)
		return
	}

	vehicleLocation := domain.VehicleLocation{
		VehicleId: payload.VehicleId,
		Longitude: payload.Longitude,
		Latitude:  payload.Latitude,
		Timestamp: payload.Timestamp,
	}

	err = u.repo.GetVehicleRepo().InsertVehicleLocation(context.Background(), &vehicleLocation)
	if err != nil {
		log.Println("error when insert vehicle location :", err)
		return
	}

	// set publish to rabbit mq for this loc -6.207557477971993,106.84506181119635 and radius <= 50m

	dst := distance.Distance{
		Latitude1:  payload.Latitude,
		Longitude1: payload.Longitude,
		Latitude2:  -6.20890615,
		Longitude2: 106.8471176,
	}

	distanceLoc := dst.GetDistanceOnMeter()
	log.Println("distance : ", distanceLoc)

	if distanceLoc <= 50 {
		bodyMap := map[string]interface{}{
			"vehicle_id": vehicleLocation.VehicleId,
			"event":      constant.GeofenceEntryEvent,
			"location": map[string]float64{
				"latitude":  payload.Latitude,
				"longitude": payload.Longitude,
			},
			"timestamp": payload.Timestamp,
		}

		err = u.rmq.PublishEvent(constant.FleetEvents, constant.GeofenceEntryEvent, bodyMap)
		if err != nil {
			log.Println("error publish geofence entry :", err)
			return
		}

		log.Println("publish message with payload :", bodyMap)
	}

}
