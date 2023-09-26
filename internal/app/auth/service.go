package auth

import (
	"fmt"
	"mcash-finance-console-core/configs"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/dto"
	"mcash-finance-console-core/internal/factory"
	"mcash-finance-console-core/internal/model"
	"mcash-finance-console-core/internal/repository"
	"mcash-finance-console-core/pkg/util/helper"
	"mcash-finance-console-core/pkg/util/response"
	res "mcash-finance-console-core/pkg/util/response"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service interface {
	Login(ctx *abstraction.Context, payload *dto.AuthLoginRequest) (*dto.AuthLoginResponse, error)
	Register(ctx *abstraction.Context, payload *dto.AuthRegisterRequest) (*dto.AuthRegisterResponse, error)
	CheckAuth(ctx *abstraction.Context, authToken string) (bool, error)
	ChangePassword(ctx *abstraction.Context, payload *dto.ChangePasswordRequest) error
	GetNotificationToken(ctx *abstraction.Context) (*dto.GetNotificationTokenResponse, error)
}

type service struct {
	Repository               repository.User
	Db                       *gorm.DB
	RolePermissionRepository repository.RolePermission
	AccessScopeRepository    repository.AccessScope
}

func NewService(f *factory.Factory) *service {
	repository := f.UserRepository
	db := f.Db
	roleRepo := f.RolePermissionRepository
	accessScopeRepo := f.AccessScopeRepository
	return &service{repository, db, roleRepo, accessScopeRepo}
}

func (s *service) Login(ctx *abstraction.Context, payload *dto.AuthLoginRequest) (*dto.AuthLoginResponse, error) {
	var result *dto.AuthLoginResponse

	data, err := s.Repository.FindByUsername(ctx, &payload.Username)
	if data == nil {
		return nil, res.CustomErrorBuilder(http.StatusBadRequest, res.E_BAD_REQUEST, "Username is incorrect")
		// return result, res.ErrorBuilder(&res.ErrorConstant.Unauthorized, err)
	}

	if data.IsActive == nil || (data.IsActive != nil && !*data.IsActive) {
		return nil, res.ErrorBuilder(&res.ErrorConstant.Unauthorized, err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(data.PasswordHash), []byte(payload.Password)); err != nil {
		return nil, res.CustomErrorBuilder(http.StatusBadRequest, res.E_BAD_REQUEST, "Password is incorrect")
		// return nil, res.ErrorBuilder(&res.ErrorConstant.Unauthorized, err)
	}

	token, err := data.GenerateToken()

	if err != nil {
		return result, res.ErrorBuilder(&res.ErrorConstant.InternalServerError, err)
	}

	jwtKey := configs.Jwt().SecretKey()

	notificationAuthToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": fmt.Sprintf("%d", data.ID),
		"exp": time.Now().Add(time.Hour * 5).Unix(),
	})

	notificationAuthTokenString, err := notificationAuthToken.SignedString([]byte(jwtKey))
	if err != nil {
		return result, response.CustomErrorBuilder(http.StatusInternalServerError, "Failed to generate token", "Please try again")
	}

	notificationSubsToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":     fmt.Sprintf("%d", data.ID),
		"channel": fmt.Sprintf("NOTIFICATION:NOTIFICATION#%d", data.ID),
		"exp":     time.Now().Add(time.Hour * 5).Unix(),
	})

	notificationSubsTokenString, err := notificationSubsToken.SignedString([]byte(jwtKey))
	if err != nil {
		return result, response.CustomErrorBuilder(http.StatusInternalServerError, "Failed to generate token", "Please try again")
	}
	result = &dto.AuthLoginResponse{
		Token:                      token,
		NotificationAuthToken:      notificationAuthTokenString,
		NotificationSubscribeToken: notificationSubsTokenString,
		UserEntityModel:            *data,
	}

	return result, nil
}

// func (s *service) Register(ctx *abstraction.Context, payload *dto.AuthRegisterRequest) (*dto.AuthRegisterResponse, error) {
// 	var result *dto.AuthRegisterResponse
// 	var data *model.UserEntityModel

// 	if err = trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
// 		data, err = s.Repository.Create(ctx, &payload)
// 		if err != nil {
// 			return res.ErrorBuilder(&res.ErrorConstant.UnprocessableEntity, err)
// 		}

// 		return nil
// 	}); err != nil {
// 		return result, err
// 	}

// 	result = &dto.AuthRegisterResponse{
// 		UserEntityModel: *data,
// 	}

// 	return result, nil
// }

func (s *service) CheckAuth(ctx *abstraction.Context, authToken string) (*dto.CheckAuthResponse, error) {
	jwtKey := configs.Jwt().SecretKey()
	splitToken := strings.Split(authToken, "Bearer ")
	token, err := jwt.Parse(splitToken[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method :%v", token.Header["alg"])
		}

		return []byte(jwtKey), nil
	})
	if !token.Valid || err != nil {
		return nil, res.ErrorBuilder(&res.ErrorConstant.Unauthorized, err)
	}
	var id int
	destructID := token.Claims.(jwt.MapClaims)["id"]
	if destructID != nil {
		id = int(destructID.(float64))
	} else {
		id = 0
	}

	data, err := s.Repository.FindByID(ctx, &id)
	if err != nil || data == nil {
		return nil, res.ErrorBuilder(&res.ErrorConstant.Unauthorized, err)
	}

	criteriaRolePermission := model.RolePermissionFilterModel{}
	criteriaRolePermission.RoleID = &data.RoleID
	rolePermission, err := s.RolePermissionRepository.FindsByCriteria(ctx, &criteriaRolePermission)
	if err != nil {
		return nil, err
	}

	criteria := model.AccessScopeFilterModel{}
	criteria.UserID = &data.ID
	dataAccessScope, err := s.AccessScopeRepository.FindByCriteria(ctx, &criteria)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}

	result := &dto.CheckAuthResponse{
		User:        *data,
		Permission:  rolePermission,
		AccessScope: dataAccessScope,
	}

	return result, nil
}

func (s *service) ChangePassword(ctx *abstraction.Context, payload *dto.ChangePasswordRequest) error {
	user, err := s.Repository.FindByID(ctx, &ctx.Auth.ID)
	if err != nil {
		return helper.ErrorHandler(err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(payload.OldPassword)); err != nil {
		return res.CustomErrorBuilder(http.StatusBadRequest, "Old password is wrong", "Please re-input old password")
	}

	modelUser := model.UserEntityModel{}
	modelUser.Context = ctx
	bytes, _ := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	modelUser.PasswordHash = string(bytes)

	user, err = s.Repository.Update(ctx, &user.ID, &modelUser)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetNotificationToken(ctx *abstraction.Context) (*dto.GetNotificationTokenResponse, error) {
	var result *dto.GetNotificationTokenResponse
	jwtKey := configs.Jwt().SecretKey()
	data, err := s.Repository.FindByID(ctx, &ctx.Auth.ID)
	if err != nil {
		return result, helper.ErrorHandler(err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": fmt.Sprintf("%d", data.ID),
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenAuthString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return result, response.CustomErrorBuilder(http.StatusInternalServerError, "Failed to generate token", "Please try again")
	}

	token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":     fmt.Sprintf("%d", data.ID),
		"channel": fmt.Sprintf("NOTIFICATION:NOTIFICATION#%d", data.ID),
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenSubscribeString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return result, response.CustomErrorBuilder(http.StatusInternalServerError, "Failed to generate token", "Please try again")
	}

	result = &dto.GetNotificationTokenResponse{
		NotificationAuthToken:      tokenAuthString,
		NotificationSubscribeToken: tokenSubscribeString,
	}

	return result, nil
}
