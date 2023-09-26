package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/dto"
	"mcash-finance-console-core/internal/factory"
	"mcash-finance-console-core/internal/model"
	"mcash-finance-console-core/internal/repository"
	"mcash-finance-console-core/pkg/redis"
	"mcash-finance-console-core/pkg/util/helper"
	"mcash-finance-console-core/pkg/util/response"
	"mcash-finance-console-core/pkg/util/trxmanager"
	"net/http"

	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

type ResetPassword struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Html    string `json:"text"`
	Subject string `json:"subject"`
	Company string `json:"company"`
}

type service struct {
	Repository                  repository.User
	AccessScopeRepository       repository.AccessScope
	AccessScopeDetailRepository repository.AccessScopeDetail
	Db                          *gorm.DB
}

type Service interface {
	Find(ctx *abstraction.Context, payload *dto.UserGetRequest) (*dto.UserGetResponse, error)
	FindByID(ctx *abstraction.Context, payload *dto.UserGetByIDRequest) (*dto.UserGetByIDResponse, error)
	Create(ctx *abstraction.Context, payload *dto.UserCreateRequest) (*dto.UserCreateResponse, error)
	Update(ctx *abstraction.Context, payload *dto.UserUpdateRequest) (*dto.UserUpdateResponse, error)
	Delete(ctx *abstraction.Context, payload *dto.UserDeleteRequest) (*dto.UserDeleteResponse, error)
	ToggleIsActive(ctx *abstraction.Context, payload *dto.UserStatusRequest) (*dto.UserStatusResponse, error)
}

func NewService(f *factory.Factory) *service {
	repository := f.UserRepository
	accessRepo := f.AccessScopeRepository
	accessDetailRepo := f.AccessScopeDetailRepository
	db := f.Db
	return &service{
		Repository:                  repository,
		Db:                          db,
		AccessScopeRepository:       accessRepo,
		AccessScopeDetailRepository: accessDetailRepo,
	}
}

func (s *service) Find(ctx *abstraction.Context, payload *dto.UserGetRequest) (*dto.UserGetResponse, error) {
	data, info, err := s.Repository.Find(ctx, &payload.UserFilterModel, &payload.Pagination)
	if err != nil {
		return &dto.UserGetResponse{}, helper.ErrorHandler(err)
	}

	result := &dto.UserGetResponse{
		Datas:          *data,
		PaginationInfo: *info,
	}
	return result, nil
}

func (s *service) FindByID(ctx *abstraction.Context, payload *dto.UserGetByIDRequest) (*dto.UserGetByIDResponse, error) {
	data, err := s.Repository.FindByID(ctx, &payload.ID)
	if err != nil {
		return &dto.UserGetByIDResponse{}, helper.ErrorHandler(err)
	}
	result := &dto.UserGetByIDResponse{
		UserEntityModel: *data,
	}
	return result, nil
}

func (s *service) Create(ctx *abstraction.Context, payload *dto.UserCreateRequest) (*dto.UserCreateResponse, error) {
	var data model.UserEntityModel

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		criteriaUser := model.UserFilterModel{}
		criteriaUser.Username = &payload.Username
		totalData, err := s.Repository.CountByCriteria(ctx, &criteriaUser)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		if totalData > 0 {
			return response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, fmt.Sprintf("Username %s already exist", payload.Username))
		}

		criteriaUser = model.UserFilterModel{}
		criteriaUser.Email = &payload.Email
		totalData, err = s.Repository.CountByCriteria(ctx, &criteriaUser)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		if totalData > 0 {
			return response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, fmt.Sprintf("Email %s already exist", payload.Email))
		}

		data.Context = ctx
		userData := model.UserEntity{}
		userData.Name = payload.Name
		userData.CompanyID = payload.CompanyID
		userData.Email = payload.Email
		userData.RoleID = payload.RoleID
		userData.Username = payload.Username
		userData.Password = payload.Password
		tmpActive := true
		userData.IsActive = &tmpActive
		userData.ImageProfile = payload.ImageProfile
		data.UserEntity = userData

		result, err := s.Repository.Create(ctx, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		accessScopeData := model.AccessScopeEntityModel{}
		accessScopeData.UserID = result.ID
		tmpFalse := false
		accessScopeData.AccessAll = &tmpFalse

		branch, err := s.AccessScopeRepository.Create(ctx, &accessScopeData)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		accessScopeDetailData := model.AccessScopeDetailEntityModel{}
		accessScopeDetailData.AccessScopeID = branch.ID
		accessScopeDetailData.CompanyID = result.CompanyID

		_, err = s.AccessScopeDetailRepository.Create(ctx, &accessScopeDetailData)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		err = redis.RedisClient.LPush(ctx.Request().Context(), fmt.Sprintf("access_company:user:%d", result.ID), result.CompanyID).Err()
		if err != nil {
			return err
		}

		data = *result
		return nil
	}); err != nil {
		return &dto.UserCreateResponse{}, err
	}

	result := &dto.UserCreateResponse{
		UserEntityModel: data,
	}
	return result, nil
}

func (s *service) Update(ctx *abstraction.Context, payload *dto.UserUpdateRequest) (*dto.UserUpdateResponse, error) {
	var data model.UserEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		userDataExisting, err := s.Repository.FindByID(ctx, &payload.ID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		if payload.Username != userDataExisting.Username {
			criteriaUser := model.UserFilterModel{}
			criteriaUser.Username = &payload.Username
			totalData, err := s.Repository.CountByCriteria(ctx, &criteriaUser)
			if err != nil {
				return helper.ErrorHandler(err)
			}

			if totalData > 0 {
				return response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, fmt.Sprintf("Username %s already exist", payload.Username))
			}
		}

		if payload.Email != userDataExisting.Email {
			criteriaUser := model.UserFilterModel{}
			criteriaUser.Email = &payload.Email
			totalData, err := s.Repository.CountByCriteria(ctx, &criteriaUser)
			if err != nil {
				return helper.ErrorHandler(err)
			}

			if totalData > 0 {
				return response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, fmt.Sprintf("Email %s already exist", payload.Email))
			}
		}

		data.Context = ctx
		userData := model.UserEntity{}
		userData.Name = payload.Name
		userData.CompanyID = payload.CompanyID
		userData.Email = payload.Email
		userData.RoleID = payload.RoleID
		userData.Username = payload.Username
		userData.Password = payload.Password
		userData.IsActive = payload.IsActive
		data.UserEntity = userData

		result, err := s.Repository.Update(ctx, &payload.ID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		accessScopeData := model.AccessScopeFilterModel{}
		accessScopeData.UserID = &userDataExisting.ID

		branch, _, err := s.AccessScopeRepository.Find(ctx, &accessScopeData, &abstraction.Pagination{})
		if err != nil {
			return helper.ErrorHandler(err)
		}

		accessScopeDetailData := model.AccessScopeDetailEntityModel{}
		accessScopeDetailData.CompanyID = result.CompanyID
		for _, vBranch := range *branch {
			accessScopeDetailData.AccessScopeID = vBranch.ID
			break
		}

		err = s.AccessScopeDetailRepository.DeleteByParent(ctx, &accessScopeDetailData.AccessScopeID, &model.AccessScopeDetailEntityModel{})
		if err != nil {
			return helper.ErrorHandler(err)
		}

		_, err = s.AccessScopeDetailRepository.Create(ctx, &accessScopeDetailData)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		if userDataExisting.CompanyID != payload.CompanyID {
			err = redis.RedisClient.Del(ctx.Request().Context(), fmt.Sprintf("access_company:user:%d", ctx.Auth.ID)).Err()
			if err != nil {
				return err
			}

			err = redis.RedisClient.LPush(ctx.Request().Context(), fmt.Sprintf("access_company:user:%d", ctx.Auth.ID), result.CompanyID).Err()
			if err != nil {
				return err
			}
		}

		data = *result
		return nil
	}); err != nil {
		return &dto.UserUpdateResponse{}, err
	}
	result := &dto.UserUpdateResponse{
		UserEntityModel: data,
	}
	return result, nil
}

func (s *service) ToggleIsActive(ctx *abstraction.Context, payload *dto.UserStatusRequest) (*dto.UserStatusResponse, error) {
	var data model.UserEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		userData, err := s.Repository.FindByID(ctx, &payload.UserID)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		tmp := !*userData.IsActive
		data.Context = ctx
		data.IsActive = &tmp
		result, err := s.Repository.Update(ctx, &payload.UserID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		data = *result
		return nil
	}); err != nil {
		return nil, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}

	result := &dto.UserStatusResponse{
		UserEntityModel: data,
	}
	return result, nil
}

func (s *service) Delete(ctx *abstraction.Context, payload *dto.UserDeleteRequest) (*dto.UserDeleteResponse, error) {
	var data model.UserEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		if _, err := s.Repository.FindByID(ctx, &payload.ID); err != nil {
			return helper.ErrorHandler(err)
		}
		data.Context = ctx
		result, err := s.Repository.Delete(ctx, &payload.ID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		err = redis.RedisClient.Del(ctx.Request().Context(), fmt.Sprintf("access_company:user:%d", ctx.Auth.ID)).Err()
		if err != nil {
			return err
		}

		data = *result
		return nil
	}); err != nil {
		return &dto.UserDeleteResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	result := &dto.UserDeleteResponse{
		UserEntityModel: data,
	}
	return result, nil
}

func (s *service) ForgotPassword(ctx *abstraction.Context, payload *dto.UserForgotPasswordRequest) (*model.UserEntityModel, error) {

	var tmpData ResetPassword

	user, err := s.Repository.FindByEmail(ctx, &payload.Email)

	if err != nil {
		return nil, response.CustomErrorBuilder(http.StatusBadRequest, payload.Email+" tidak terdaftar sebagai user", payload.Email+" tidak terdaftar sebagai user")
	}

	resetToken, err := user.GenerateTokenResetPassword()
	// passwordResetToken := utils.Encode(resetToken)

	if err != nil {
		err = errors.New("Invalid User .. User doesn't exists")
		return nil, err
	}

	content := fmt.Sprintf("https:://konsolidasi.codeoffice.net/reset-password/%s", resetToken)

	tmpData.To = user.Email
	tmpData.Html = content
	tmpData.Subject = "RESET PASSWORD"
	tmpData.Company = user.Company.Name
	tmpData.From = "danilelouch@gmail.com"

	jsonData, err := json.Marshal(tmpData)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	byteData := []byte(jsonData)
	// template := "templates/resetpassword.html"

	sendResetPasswordEmail(byteData)
	fmt.Println("silahkan cek email")
	return nil, err
}

func ParseTemplate(templateFileName string, payload []byte) (string, error) {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		fmt.Println("erorr parse template")
		return "erorr parse template", err
	}

	var details ResetPassword
	err = json.Unmarshal(payload, &details)

	buf := new(bytes.Buffer)
	if err = t.Execute(buf, details); err != nil {
		fmt.Println("erorr parse template")
		return "erorr parse template", err
	}
	return buf.String(), nil
}

func sendResetPasswordEmail(payload []byte) error {

	var details ResetPassword
	json.Unmarshal(payload, &details)
	// result, _ := ParseTemplate(templateFile, payload)

	// Render template ke buffer
	m := gomail.NewMessage()

	m.SetHeader("From", details.From)
	m.SetHeader("To", details.To)
	m.SetHeader("Subject", details.Subject)
	m.SetBody("text/html", details.Html)

	d := gomail.NewDialer("smtp.gmail.com", 465, "ahmadhamdani040995@gmail.com", "tuagldbhmqnggxae")

	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("email sudah terkirim")
	return nil
}

func (s *service) ResetPassword(ctx *abstraction.Context, payload *dto.UserResetPasswordRequest) (*dto.UserResetPasswordResponse, error) {
	var result *dto.UserResetPasswordResponse

	resetToken := ctx.Param("resetToken")

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		_, err := s.Repository.FindByID(ctx, &ctx.Auth.ID)
		if err != nil {
			return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err)
		}

		payload.ResetToken = model.Encode(resetToken)
		hashPassword, _ := model.HashPassword(payload.Password)

		err = s.Repository.UpdateUserPassword(ctx, &payload.ResetToken, hashPassword)
		if err != nil {
			return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err)
		}
		return nil
	}); err != nil {
		return result, err
	}

	return result, nil
}