package handler

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rhuandantas/chat-jobsity/internal/model"
	"github.com/rhuandantas/chat-jobsity/internal/util"
	"net/http"
)

type User struct {
	validator *util.CustomValidator
	users     map[string]*model.User
}

func NewUserHandler(validator *util.CustomValidator) *User {
	//TODO remove this
	users := map[string]*model.User{
		"rhuannixon@gmail.com": &model.User{
			ID:       uuid.New(),
			Email:    "rhuannixon@gmail.com",
			Password: "123456",
		},
	}
	return &User{
		validator: validator,
		users:     users,
	}
}

func (h *User) RegisterUserApi(server *echo.Echo) {
	server.POST("/users/signup", h.SignUp)
	server.POST("/users/signin", h.SignIn)
}

func (h *User) SignUp(ctx echo.Context) error {
	user := model.User{}
	err := ctx.Bind(&user)
	if err == nil {
		err = h.validator.Validate(user)
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return ctx.JSON(http.StatusCreated, "User Registered")
}

func (h *User) SignIn(ctx echo.Context) error {
	user := model.User{}
	err := ctx.Bind(&user)
	if err == nil {
		err = h.validator.Validate(user)
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	//TODO search into the database

	//TODO validate encrypted password

	//TODO create JWT token
	return ctx.JSON(http.StatusOK, "User authorized")
}
