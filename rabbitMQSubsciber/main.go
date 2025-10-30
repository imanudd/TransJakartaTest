package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type GeofenceEvent struct {
	VehicleID string `json:"vehicle_id"`
	Event     string `json:"event"`
	Location  struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"location"`
	Timestamp int64 `json:"timestamp"`
}

func main() {
	conn, err := amqp.Dial("amqp://admin:admin123@localhost:5672/")
	if err != nil {
		log.Fatalf("gagal konek ke RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("gagal buka channel: %v", err)
	}
	defer ch.Close()

	queueName := "geofence_alerts"

	_, err = ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("gagal declare queue: %v", err)
	}

	msgs, err := ch.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("gagal consume: %v", err)
	}

	log.Println("Menunggu pesan dari queue:", queueName)

	ctx := context.Background()

	for msg := range msgs {
		go func(d amqp.Delivery) {
			defer func() {
				d.Ack(false)
			}()

			var event GeofenceEvent
			if err := json.Unmarshal(d.Body, &event); err != nil {
				log.Printf("gagal parse pesan: %v", err)
				return
			}

			log.Printf("Diterima event dari %s: %+v\n", event.VehicleID, event)
			processGeofenceEvent(ctx, event)
		}(msg)
	}
}

func processGeofenceEvent(ctx context.Context, event GeofenceEvent) {
	log.Printf("🚗 [%s] %s di lokasi (%.5f, %.5f) @ %s\n",
		event.VehicleID,
		event.Event,
		event.Location.Latitude,
		event.Location.Longitude,
		time.Unix(event.Timestamp, 0).Format("2006-01-02 15:04:05"),
	)
}
