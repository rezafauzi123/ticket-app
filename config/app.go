package config

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type AppConfig struct {
	Db           *sqlx.DB
	Log          *logrus.Logger
	RabbitMQConn struct {
		Conn *amqp.Connection
		Ch   *amqp.Channel
	}
}
