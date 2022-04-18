package handler

import (
	"crypto/sha256"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rhuandantas/chat-jobsity/internal/api/middleware"
	"github.com/rhuandantas/chat-jobsity/internal/model"
	"github.com/rhuandantas/chat-jobsity/internal/util"
	"net/http"
	"reflect"
)

type User struct {
	validator *util.CustomValidator
	users     map[string]*model.User
}

func NewUserHandler(validator *util.CustomValidator, users map[string]*model.User) *User {
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

	_, found := h.users[user.Email]
	if found {
		return ctx.JSON(http.StatusBadRequest, "User already exists")
	}

	err = h.createUser(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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
	userFound := h.users[user.Email]
	if userFound == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "User not found")
	}

	if !validatePassword(user.Password, userFound.Password) {
		return echo.NewHTTPError(http.StatusUnauthorized, "User or Pass is not valid")
	}

	token, err := middleware.GenerateAccessToken(userFound.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	//TODO create JWT session
	return ctx.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}

func (h *User) createUser(user model.User) error {
	user.ID = uuid.New()
	user.Password = encryptPassword(user.Password)
	h.users[user.Email] = &user
	return nil
}

func encryptPassword(password string) string {
	sum := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", sum)
}

func validatePassword(password, encrypted string) bool {
	return reflect.DeepEqual(encryptPassword(password), encrypted)
}
