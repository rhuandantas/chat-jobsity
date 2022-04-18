package handler

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rhuandantas/chat-jobsity/internal/model"
	"github.com/rhuandantas/chat-jobsity/internal/service"
	"github.com/rhuandantas/chat-jobsity/internal/util"
	"net/http"
)

type Chat struct {
	validator      *util.CustomValidator
	messageService *service.Message
	chats          map[string]*model.Chat
}

func NewChatHandler(validator *util.CustomValidator, messageService *service.Message) *Chat {
	//TODO remove
	chats := map[string]*model.Chat{
		"aaaaaa": &model.Chat{},
	}
	return &Chat{
		validator:      validator,
		messageService: messageService,
		chats:          chats,
	}
}

func (h *Chat) RegisterChatApi(server *echo.Echo) {
	group := server.Group("/chats")
	group.Use(middleware.JWT([]byte("secret")))
	group.GET("/", h.GetChat)
	group.POST("/:id/message", h.sendMessage)
}

func (h Chat) GetChat(ctx echo.Context) error {
	roomName := ctx.QueryParam("name")

	if roomName == "" && len(roomName) < 3 {
		return echo.NewHTTPError(http.StatusBadRequest, "Please type a room name valid")
	}

	chat := h.chats[roomName]
	if chat != nil {
		return ctx.JSON(http.StatusOK, chat)
	}
	chat = &model.Chat{
		ID:          uuid.New(),
		Description: roomName,
		Messages:    model.Messages{},
	}
	// TODO save chat into database
	h.chats[roomName] = chat

	return ctx.JSON(http.StatusCreated, chat)
}

func (h Chat) sendMessage(ctx echo.Context) error {
	message := model.Message{}
	chatID := ctx.Param("id")
	if chatID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Please provide chat id you want to send messages")
	}

	// TODO move to service layer
	chat := h.chats[chatID]
	if chat == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "This chat doesn't exist")
	}

	err := ctx.Bind(&message)
	if err == nil {
		err = h.validator.Validate(message)
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	message.ChatID = chatID
	err = h.messageService.SendMessage(message)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	// FIXME this part must be cacheable and moved to another layer
	chat.Messages = append(chat.Messages, message)
	h.chats[chatID] = chat
	return ctx.JSON(http.StatusOK, "Message sent")
}
