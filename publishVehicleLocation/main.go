package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type VehicleLocation struct {
	VehicleID string  `json:"vehicle_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp int64   `json:"timestamp"`
}

func main() {
	broker := "tcp://localhost:1883"
	clientID := "go-mqtt-publisher"
	vehicleID := "B1234XYZ"
	topic := fmt.Sprintf("/fleet/vehicle/%s/location", vehicleID)

	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	opts.OnConnect = func(c mqtt.Client) {
		fmt.Println("Connected to MQTT broker")
	}
	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		fmt.Println("Connection lost:", err)
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	lat := -6.208753924495224
	lon := 106.84753593871517

	rand.Seed(time.Now().UnixNano())

	for i := 1; i <= 100; i++ {
		lat += (rand.Float64() - 0.5) * 0.001
		lon += (rand.Float64() - 0.5) * 0.001

		data := VehicleLocation{
			VehicleID: vehicleID,
			Latitude:  lat,
			Longitude: lon,
			Timestamp: time.Now().Unix(),
		}

		payload, _ := json.Marshal(data)

		token := client.Publish(topic, 0, false, payload)
		token.Wait()

		fmt.Printf("[%02d] Published to %s: %s\n", i, topic, string(payload))
		time.Sleep(2 * time.Second)
	}

	client.Disconnect(250)
	fmt.Println("Done sending 100 mock messages.")
}
