package client

import (
	"fmt"
	"github.com/rhuandantas/chat-jobsity/internal/config"
	"github.com/rhuandantas/chat-jobsity/internal/model"
)

type RabbitMQ struct {
	cfg config.Config
}

func NewRabbitMQClient(cfg config.Config) *RabbitMQ {
	return &RabbitMQ{
		cfg: cfg,
	}
}

func (c RabbitMQ) PublishMessage(message model.Message) error {
	//TODO implement this
	fmt.Println("Message published", message)
	return nil
}
