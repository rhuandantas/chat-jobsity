package main

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rhuandantas/chat-jobsity/internal/api"
	"github.com/rhuandantas/chat-jobsity/internal/api/handler"
	"github.com/rhuandantas/chat-jobsity/internal/client"
	"github.com/rhuandantas/chat-jobsity/internal/config"
	"github.com/rhuandantas/chat-jobsity/internal/model"
	"github.com/rhuandantas/chat-jobsity/internal/service"
	"github.com/rhuandantas/chat-jobsity/internal/util"
)

func main() {
	e := api.CreateServer()
	//TODO remove this
	users := map[string]*model.User{
		"rhuannixon@gmail.com": &model.User{
			ID:       uuid.New(),
			Email:    "rhuannixon@gmail.com",
			Password: "123456",
		},
	}
	//Dependency injection
	cfg := config.GetConfig()
	validator := util.NewCustomValidator()
	clientStooq := client.NewStooqClient(cfg)
	clientRabbitMq := client.NewRabbitMQClient(cfg)
	serviceMessage := service.NewMessageService(clientStooq, clientRabbitMq)
	userHandler := handler.NewUserHandler(validator, users)
	chatHandler := handler.NewChatHandler(validator, serviceMessage)

	//Register
	userHandler.RegisterUserApi(e)
	chatHandler.RegisterChatApi(e)

	//Start server
	e.Use(middleware.Recover())
	e.Logger.Fatal(e.Start(":8080"))
}
