package service

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/rhuandantas/chat-jobsity/internal/client"
	"github.com/rhuandantas/chat-jobsity/internal/model"
	"regexp"
	"strings"
	"time"
)

type Message struct {
	stockClient    *client.Stock
	rabbitMQClient *client.RabbitMQ
}

func NewMessageService(stockClient *client.Stock, rabbitMQClient *client.RabbitMQ) *Message {
	return &Message{
		stockClient:    stockClient,
		rabbitMQClient: rabbitMQClient,
	}
}

func (s *Message) SendStockMessage(author, stockCode string) error {
	stock, err := s.stockClient.GetStock(stockCode)
	if err != nil {
		return err
	}
	message := model.Message{
		ID:     uuid.New(),
		Text:   fmt.Sprintf("%s quote is $%.2f per share", strings.ToUpper(stockCode), stock.Open),
		Author: author,
		Time:   time.Now().String(),
	}

	return s.rabbitMQClient.PublishMessage(message)
}

func (s *Message) SendMessage(message model.Message) error {
	if s.isCommand(message.Text) {
		return s.doCommand(message.Text, message.Author)
	}
	return s.rabbitMQClient.PublishMessage(message)
}

func (s *Message) doCommand(text, author string) error {
	command := strings.Split(text, "=")
	if len(command) < 2 {
		return errors.New("invalid command")
	}

	commandCode := strings.ToLower(command[0])
	commandValue := strings.ToLower(command[1])

	switch commandCode {
	case StockCommand:
		return s.SendStockMessage(author, commandValue)
	default:
		return errors.New(fmt.Sprintf("command %s doesn't supported", commandCode))
	}
}

func (s *Message) isCommand(text string) bool {
	match, _ := regexp.MatchString("\\/\\w[a-zA-Z]{3,}=[a-z][a-zA-Z.-]{3,}$", text)
	startWithSlash := text[:1] == "/"
	return startWithSlash || match
}

const (
	StockCommand string = "/stock"
	ListCommand  string = "/list"
)
