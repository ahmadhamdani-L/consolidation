package user

import (
	"encoding/base64"
	"fmt"
	"io"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/dto"
	"mcash-finance-console-core/internal/factory"
	"mcash-finance-console-core/pkg/util/helper"
	"mcash-finance-console-core/pkg/util/response"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type handler struct {
	service *service
}

func NewHandler(f *factory.Factory) *handler {
	return &handler{
		service: NewService(f),
	}
}

var tmpFileUploadPath = path.Join(os.Getenv("STORAGE_DIRECTORY_PATH"), "/assets/user_images")

// Get
// @Summary Get User
// @Description Get User
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @param request query dto.UserGetRequest true "request query"
// @Success 200 {object} dto.UserGetResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /user [get]
func (h *handler) Get(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.UserGetRequest)
	if err := c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}
	rolePayload, err := helper.MultiRoleFilter(c.Request().URL.Query())
	if err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}
	if len(rolePayload) > 0 {
		payload.ArrRoleID = &rolePayload
	}
	if err := c.Validate(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}
	result, err := h.service.Find(cc, payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}
	return response.CustomSuccessBuilder(http.StatusOK, result.Datas, "Get Data Success", &result.PaginationInfo).Send(c)
}

// Get By ID
// @Summary Get User by id
// @Description Get User by id
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "id path"
// @Success 200 {object} dto.UserGetByIDResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /user/{id} [get]
func (h *handler) GetByID(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.UserGetByIDRequest)

	if err := c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	result, err := h.service.FindByID(cc, payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}

	return response.SuccessResponse(result).Send(c)
}

// Create godoc
// @Summary Create User
// @Description Create User
// @Tags User
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param request body dto.UserCreateRequest true "request body"
// @Success 200 {object} dto.UserCreateResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /user [post]
func (h *handler) Create(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.UserCreateRequest)

	if err := c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	file, err := c.FormFile("image_profile")
	if err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}
	src, err := file.Open()
	if err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}
	defer src.Close()

	_, err = os.Stat(tmpFileUploadPath)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(tmpFileUploadPath, os.ModePerm); err != nil {
				return response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err).Send(c)
			}
		} else {
			return response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err).Send(c)
		}
	}

	tmpExtFile := filepath.Ext(file.Filename)
	if tmpExtFile != ".jpg" && tmpExtFile != ".jpeg" && tmpExtFile != ".png" {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}
	timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	tmpFileName := fmt.Sprintf("%s%s", base64.StdEncoding.EncodeToString([]byte(timestamp)), tmpExtFile)

	pathFile := fmt.Sprintf("%s/%s", tmpFileUploadPath, tmpFileName)
	payload.ImageProfile = tmpFileName
	dst, err := os.Create(pathFile)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	defer dst.Close()

	result, err := h.service.Create(cc, payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}

	return response.SuccessResponse(result).Send(c)
}

// Update godoc
// @Summary Update User
// @Description Update User
// @Tags User
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path int true "id path"
// @Param request body dto.UserUpdateRequest true "request body"
// @Success 200 {object} dto.UserUpdateResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /user/{id} [patch]
func (h *handler) Update(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.UserUpdateRequest)

	if err := c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	file, err := c.FormFile("image_profile")
	if err == nil {
		src, err := file.Open()
		if err != nil {
			return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
		}
		defer src.Close()

		_, err = os.Stat(tmpFileUploadPath)
		if err != nil {
			if os.IsNotExist(err) {
				if err := os.MkdirAll(tmpFileUploadPath, os.ModePerm); err != nil {
					return response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err).Send(c)
				}
			} else {
				return response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err).Send(c)
			}
		}

		tmpExtFile := filepath.Ext(file.Filename)
		timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
		tmpFileName := fmt.Sprintf("%s%s", base64.StdEncoding.EncodeToString([]byte(timestamp)), tmpExtFile)

		pathFile := fmt.Sprintf("%s/%s", tmpFileUploadPath, tmpFileName)
		payload.ImageProfile = tmpFileName
		dst, err := os.Create(pathFile)
		if err != nil {
			return err
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			return err
		}
	} else {
		if err.Error() != "http: no such file" {
			return response.ErrorResponse(err).Send(c)
		}
	}

	result, err := h.service.Update(cc, payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}

	return response.SuccessResponse(result).Send(c)
}

// Delete godoc
// @Summary Delete User
// @Description Delete User
// @Tags User
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path int true "id path"
// @Success 200 {object} dto.UserDeleteResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /user/{id} [delete]
func (h *handler) Delete(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.UserDeleteRequest)

	if err := c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	result, err := h.service.Delete(cc, payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}

	return response.SuccessResponse(result).Send(c)
}

func (h *handler) ForgotPassword(c echo.Context) error {
	cc := c.(*abstraction.Context)

	payload := new(dto.UserForgotPasswordRequest)
	if err := c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	data, err := h.service.ForgotPassword(cc, payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}

	return response.SuccessResponse(data).Send(c)
}

func (h *handler) ResetPassword(c echo.Context) error {
	cc := c.(*abstraction.Context)

	payload := new(dto.UserResetPasswordRequest)
	if err := c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	result, err := h.service.ResetPassword(cc, payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}

	return response.SuccessResponse(result).Send(c)
}

func (h *handler) UserActive(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.UserStatusRequest)
	if err := c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	result, err := h.service.ToggleIsActive(cc, payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}

	return response.SuccessResponse(result).Send(c)
}
