package main

import (
	"ticket-app/config"
	"ticket-app/internal/router"
	dbPkg "ticket-app/pkg/db"
	"ticket-app/pkg/log"
	"ticket-app/pkg/rabbitmq"
)

func main() {
	log.InitLogger()
	logger := log.GetLogger()

	cfg := config.LoadDBConfig()
	db, err := dbPkg.ConnectDB(cfg)
	if err != nil {
		logger.Error(err)
		panic(err)
	}

	dbPkg.InitRedis()

	conn, err := rabbitmq.ConnectRabbitMQ()
	if err != nil {
		logger.Fatalf("Could not connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// App Init
	appConfig := config.AppConfig{
		Db:           db,
		Log:          logger,
		RabbitMQConn: *conn,
	}

	router.Router(appConfig)
}
