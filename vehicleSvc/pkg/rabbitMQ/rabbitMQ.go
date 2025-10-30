package rabbitMQ

import (
	"app/config"
	"app/utils/constant"
	"encoding/json"
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

var mRouting = map[string]route{
	constant.GeofenceEntryEvent: {
		eventName:    constant.GeofenceEntryEvent,
		exchangeName: constant.FleetEvents,
		queueName:    constant.GeofenceAlerts,
	},
}

type route struct {
	eventName    string
	exchangeName string
	queueName    string
}

type RabbitMQImpl interface {
	Init()
	Close()
	ExchangeDeclare(exchangeName string, exchangeType string) error
	QueueDeclare(queueName string) error
	QueueBind(queueName, routingKey, exchangeName string) error
	Configure(eventName string) error
	PublishEvent(exchange, eventName string, req any) error
}

type RabbitMQ struct {
	cfg     *config.MainConfig
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQImpl(cfg *config.MainConfig) RabbitMQImpl {
	return &RabbitMQ{
		cfg: cfg,
	}
}

func (r *RabbitMQ) Configure(eventName string) error {
	if _, ok := mRouting[eventName]; !ok {
		return errors.New("eventName not found in routing config")
	}

	routes := mRouting[eventName]

	if err := r.ExchangeDeclare(routes.exchangeName, constant.ExchangeTypeTopic); err != nil {
		return err
	}

	if err := r.QueueDeclare(routes.queueName); err != nil {
		return err
	}

	if err := r.QueueBind(routes.queueName, routes.eventName, routes.exchangeName); err != nil {
		return err
	}

	return nil
}

func (r *RabbitMQ) ExchangeDeclare(exchangeName string, exchangeType string) error {
	return r.channel.ExchangeDeclare(exchangeName, exchangeType, true, false, false, false, nil)
}

func (r *RabbitMQ) QueueDeclare(queueName string) error {
	_, err := r.channel.QueueDeclare(queueName, true, false, false, false, nil)
	return err
}

func (r *RabbitMQ) QueueBind(queueName, routingKey, exchangeName string) error {
	return r.channel.QueueBind(queueName, routingKey, exchangeName, false, nil)
}

func (r *RabbitMQ) Init() {
	for i := 0; i < 10; i++ {
		conn, err := amqp.DialConfig(r.cfg.RabbitMQUrl, amqp.Config{
			Heartbeat: 10 * time.Second,
			Locale:    "en_US",
		})
		if err == nil {
			r.conn = conn
			fmt.Println(r.cfg.RabbitMQUrl)
			break
		}

		fmt.Println("❌ Failed to connect to RabbitMQ, retrying in 2s...")
		time.Sleep(2 * time.Second)
	}

	ch, err := r.conn.Channel()
	if err != nil {
		log.Fatalf("gagal buka channel: %v", err)
	}

	r.channel = ch
	return

}

func (r *RabbitMQ) PublishEvent(exchange, eventName string, req any) error {
	bodyJson, err := json.Marshal(req)
	if err != nil {
		return err
	}

	err = r.channel.Publish(exchange, eventName, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        bodyJson,
	})
	if err != nil {
		log.Println(fmt.Sprintf("gagal publish ke RabbitMQ: %v", err))
		return err
	}

	log.Println("Published event: ", eventName)
	return nil

}

func (r *RabbitMQ) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}
