package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"app/config"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func InitMosquitto(cfg *config.MainConfig) *MQTT.Client {
	opts := MQTT.NewClientOptions()
	opts.AddBroker(cfg.MQTTUrl)
	opts.SetClientID("go-subscriber-1")
	opts.OnConnect = func(c MQTT.Client) {
		fmt.Println("✅ Connected to MQTT broker")

		if token := c.Subscribe("test/topic", 0, nil); token.Wait() && token.Error() != nil {
			fmt.Println("Subscribe error:", token.Error())
		} else {
			fmt.Println("Subscribed to topic: test/topic")
		}
	}
	opts.OnConnectionLost = func(c MQTT.Client, err error) {
		fmt.Printf("Connection lost: %v\n", err)
	}

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Println("Connected to MQTT broker")
	return &client
}

func NewPostgres(cfg *config.MainConfig) *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresUsername,
		cfg.PostgresPassword,
		cfg.DBName,
	)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

// init postgress with gorm
func InitPostgreSQL(cfg *config.MainConfig) *gorm.DB {
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Asia/Jakarta", cfg.PostgresHost, cfg.PostgresUsername, cfg.PostgresPassword, cfg.DBName, cfg.PostgresPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		SkipDefaultTransaction: true,
	})
	if err != nil {
		log.Fatalf(err.Error())
		panic(err)
	}

	log.Printf("Successfully connected to database server")

	rdb, err := db.DB()
	if err != nil {
		log.Fatalf(err.Error())
		panic(err)
	}

	rdb.SetMaxIdleConns(cfg.MaxIdleConns)
	rdb.SetMaxOpenConns(cfg.MaxOpenConns)
	rdb.SetConnMaxLifetime(time.Duration(int(time.Minute) * cfg.ConnMaxLifetime))

	return db
}
