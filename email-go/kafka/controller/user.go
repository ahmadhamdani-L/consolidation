package controller

import (
	"encoding/json"
	"mail-go/model"
	"mail-go/kafka/config"
	"mail-go/kafka/producer"
	"mail-go/util"
	"net/http"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

type UserController struct {
	producer producer.Producer
}

func NewUserController(producer producer.Producer) *UserController {
	return &UserController{producer: producer}
}

// SaveUser godoc
// @Summary Create a user
// @Description Create a new user item
// @Tags users
// @Accept json,xml
// @Produce json
// @Param mediaType query string false "mediaType" Enums(json, xml)
// @Param user body model.UserInput true "New User"
// @Success 200 {object} model.User
// @Failure 500 {object} handler.APIError
// @Router /signup [post]
func (userController *UserController) SaveUser(c echo.Context) error {
	payload := new(model.UserInput)
	if err := util.BindAndValidate(c, payload); err != nil {
		return err
	}

	user := &model.User1{UserInput: payload}
	user.ID = uuid.NewV4().String()

	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	userController.producer.Produce(config.UserNotificationTopic, string(data))

	return util.Negotiate(c, http.StatusOK, user)
}
