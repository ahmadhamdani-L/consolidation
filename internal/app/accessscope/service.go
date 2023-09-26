package accessscope

import (
	"fmt"
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

	"gorm.io/gorm"
)

type service struct {
	Repository        repository.AccessScope
	AccessScopeDetail repository.AccessScopeDetail
	UserRepository    repository.User
	CompanyRepository repository.Company
	Db                *gorm.DB
}

type Service interface {
	Find(ctx *abstraction.Context, payload *dto.AccessScopeGetRequest) (*dto.AccessScopeGetResponse, error)
	FindByID(ctx *abstraction.Context, payload *dto.AccessScopeGetByIDRequest) (*dto.AccessScopeGetByIDResponse, error)
	Create(ctx *abstraction.Context, payload *dto.AccessScopeCreateRequest) (*dto.AccessScopeCreateResponse, error)
	Update(ctx *abstraction.Context, payload *dto.AccessScopeUpdateRequest) (*dto.AccessScopeUpdateResponse, error)
	Delete(ctx *abstraction.Context, payload *dto.AccessScopeDeleteRequest) (*dto.AccessScopeDeleteResponse, error)
}

func NewService(f *factory.Factory) *service {
	repository := f.AccessScopeRepository
	accessDetail := f.AccessScopeDetailRepository
	userRepo := f.UserRepository
	companyRepo := f.CompanyRepository
	db := f.Db
	return &service{
		Repository:        repository,
		Db:                db,
		AccessScopeDetail: accessDetail,
		UserRepository:    userRepo,
		CompanyRepository: companyRepo,
	}
}

func (s *service) Find(ctx *abstraction.Context, payload *dto.AccessScopeGetRequest) (*dto.AccessScopeGetResponse, error) {
	data, info, err := s.Repository.Find(ctx, &payload.AccessScopeFilterModel, &payload.Pagination)
	if err != nil {
		return &dto.AccessScopeGetResponse{}, helper.ErrorHandler(err)
	}

	result := &dto.AccessScopeGetResponse{
		Datas:          *data,
		PaginationInfo: *info,
	}
	return result, nil
}

func (s *service) FindByID(ctx *abstraction.Context, payload *dto.AccessScopeGetByIDRequest) (*dto.AccessScopeGetByIDResponse, error) {
	data, err := s.Repository.FindByID(ctx, &payload.ID)
	if err != nil {
		return &dto.AccessScopeGetByIDResponse{}, helper.ErrorHandler(err)
	}
	result := &dto.AccessScopeGetByIDResponse{
		AccessScopeEntityModel: *data,
	}
	return result, nil
}

func (s *service) Create(ctx *abstraction.Context, payload *dto.AccessScopeCreateRequest) (*dto.AccessScopeCreateResponse, error) {
	var data model.AccessScopeEntityModel

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		data.Context = ctx
		data.AccessScopeEntity = payload.AccessScopeEntity

		result, err := s.Repository.Create(ctx, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.AccessScopeCreateResponse{}, err
	}

	result := &dto.AccessScopeCreateResponse{
		AccessScopeEntityModel: data,
	}
	return result, nil
}

func (s *service) Update(ctx *abstraction.Context, payload *dto.AccessScopeUpdateRequest) (*dto.AccessScopeUpdateResponse, error) {
	var data model.AccessScopeEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		accessScope, err := s.Repository.FindByID(ctx, &payload.ID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		userData, err := s.UserRepository.FindByID(ctx, &accessScope.UserID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		data.Context = ctx
		dataAccess := model.AccessScopeEntity{}
		dataAccess.AccessAll = payload.AccessAll
		data.AccessScopeEntity = dataAccess

		result, err := s.Repository.Update(ctx, &payload.ID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		nilai := 0
		if result.AccessAll != nil && *result.AccessAll {
			nilai = 1
		}
		err = redis.RedisClient.Set(ctx.Request().Context(), fmt.Sprintf("access_all_company:user:%d", accessScope.User.ID), nilai, 0).Err()
		if err != nil {
			return err
		}

		err = redis.RedisClient.Del(ctx.Request().Context(), fmt.Sprintf("access_company:user:%d", accessScope.User.ID)).Err()
		if err != nil {
			return err
		}
		criteriaASD := model.AccessScopeDetailEntityModel{}
		criteriaASD.CompanyID = accessScope.User.CompanyID
		err = s.AccessScopeDetail.DeleteByParent(ctx, &accessScope.ID, &criteriaASD)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		err = redis.RedisClient.LPush(ctx.Request().Context(), fmt.Sprintf("access_company:user:%d", accessScope.User.ID), accessScope.User.CompanyID).Err()
		if err != nil {
			return err
		}
		if result.AccessAll != nil && !*result.AccessAll {
			for _, v := range payload.CompanyID {
				if v == userData.CompanyID {
					continue
				}
				_, err := s.CompanyRepository.FindByID(ctx, &v)
				if err != nil {
					if err == gorm.ErrRecordNotFound {
						continue
					}
					return response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Perusahaan tidak terdaftar")
				}
				accessDetailData := model.AccessScopeDetailEntityModel{}
				accessDetailData.AccessScopeID = accessScope.ID
				accessDetailData.CompanyID = v
				_, err = s.AccessScopeDetail.Create(ctx, &accessDetailData)
				if err != nil {
					return helper.ErrorHandler(err)
				}
				err = redis.RedisClient.LPush(ctx.Request().Context(), fmt.Sprintf("access_company:user:%d", accessScope.User.ID), v).Err()
				if err != nil {
					return err
				}
			}
		}

		data = *result
		return nil
	}); err != nil {
		return &dto.AccessScopeUpdateResponse{}, err
	}
	result := &dto.AccessScopeUpdateResponse{
		AccessScopeEntityModel: data,
	}
	return result, nil
}

func (s *service) Reset(ctx *abstraction.Context, payload *dto.AccessScopeDeleteRequest) (*dto.AccessScopeUpdateResponse, error) {
	var data model.AccessScopeEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		accessScope, err := s.Repository.FindByID(ctx, &payload.ID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		user, err := s.UserRepository.FindByID(ctx, &accessScope.UserID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		data.Context = ctx
		dataAccess := model.AccessScopeEntity{}
		tmpFalse := false
		dataAccess.AccessAll = &tmpFalse
		data.AccessScopeEntity = dataAccess

		result, err := s.Repository.Update(ctx, &payload.ID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		err = redis.RedisClient.Set(ctx.Request().Context(), fmt.Sprintf("access_all_company:user:%d", accessScope.User.ID), 0, 0).Err()
		if err != nil {
			return err
		}

		criteriaASD := model.AccessScopeDetailEntityModel{}
		criteriaASD.CompanyID = user.CompanyID
		err = s.AccessScopeDetail.DeleteByParent(ctx, &accessScope.ID, &criteriaASD)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		err = redis.RedisClient.Del(ctx.Request().Context(), fmt.Sprintf("access_company:user:%d", accessScope.User.ID)).Err()
		if err != nil {
			return err
		}

		err = redis.RedisClient.LPush(ctx.Request().Context(), fmt.Sprintf("access_company:user:%d", accessScope.User.ID), user.CompanyID).Err()
		if err != nil {
			return err
		}

		data = *result
		return nil
	}); err != nil {
		return &dto.AccessScopeUpdateResponse{}, err
	}
	result := &dto.AccessScopeUpdateResponse{
		AccessScopeEntityModel: data,
	}
	return result, nil
}

func (s *service) Delete(ctx *abstraction.Context, payload *dto.AccessScopeDeleteRequest) (*dto.AccessScopeDeleteResponse, error) {
	var data model.AccessScopeEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		accessScope, err := s.Repository.FindByID(ctx, &payload.ID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		data.Context = ctx
		result, err := s.Repository.Delete(ctx, &payload.ID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		err = redis.RedisClient.Del(ctx.Request().Context(), fmt.Sprintf("access_all_company:user:%d", accessScope.User.ID)).Err()
		if err != nil {
			return err
		}

		data = *result
		return nil
	}); err != nil {
		return &dto.AccessScopeDeleteResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	result := &dto.AccessScopeDeleteResponse{
		AccessScopeEntityModel: data,
	}
	return result, nil
}
