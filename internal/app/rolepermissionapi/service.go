package rolepermissionapi

import (
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/dto"
	"mcash-finance-console-core/internal/factory"
	"mcash-finance-console-core/internal/model"
	"mcash-finance-console-core/internal/repository"
	"mcash-finance-console-core/pkg/redis"
	"mcash-finance-console-core/pkg/util/helper"
	"mcash-finance-console-core/pkg/util/trxmanager"

	"gorm.io/gorm"
)

type service struct {
	Repository            repository.RolePermissionApi
	AccessScopeRepository repository.AccessScope
	Db                    *gorm.DB
}

type Service interface {
	Find(ctx *abstraction.Context, payload *dto.RolePermissionApiGetRequest) (*dto.RolePermissionApiGetResponse, error)
	FindByID(ctx *abstraction.Context, payload *dto.RolePermissionApiGetByIDRequest) (*dto.RolePermissionApiGetByIDResponse, error)
	Create(ctx *abstraction.Context, payload *dto.RolePermissionApiCreateRequest) (*dto.RolePermissionApiCreateResponse, error)
	Update(ctx *abstraction.Context, payload *dto.RolePermissionApiUpdateRequest) (*dto.RolePermissionApiUpdateResponse, error)
	Delete(ctx *abstraction.Context, payload *dto.RolePermissionApiDeleteRequest) (*dto.RolePermissionApiDeleteResponse, error)
	RoleCacheInit()
}

func NewService(f *factory.Factory) *service {
	repository := f.RolePermissionApiRepository
	accessScopeRepo := f.AccessScopeRepository
	db := f.Db
	return &service{
		Repository:            repository,
		Db:                    db,
		AccessScopeRepository: accessScopeRepo,
	}
}

func (s *service) Find(ctx *abstraction.Context, payload *dto.RolePermissionApiGetRequest) (*dto.RolePermissionApiGetResponse, error) {
	data, info, err := s.Repository.Find(ctx, &payload.RolePermissionApiFilterModel, &payload.Pagination)
	if err != nil {
		return &dto.RolePermissionApiGetResponse{}, helper.ErrorHandler(err)
	}

	result := &dto.RolePermissionApiGetResponse{
		Datas:          *data,
		PaginationInfo: *info,
	}
	return result, nil
}

func (s *service) FindByID(ctx *abstraction.Context, payload *dto.RolePermissionApiGetByIDRequest) (*dto.RolePermissionApiGetByIDResponse, error) {
	data, err := s.Repository.FindByID(ctx, &payload.ID)
	if err != nil {
		return &dto.RolePermissionApiGetByIDResponse{}, helper.ErrorHandler(err)
	}
	result := &dto.RolePermissionApiGetByIDResponse{
		RolePermissionApiEntityModel: *data,
	}
	return result, nil
}

func (s *service) Create(ctx *abstraction.Context, payload *dto.RolePermissionApiCreateRequest) (*dto.RolePermissionApiCreateResponse, error) {
	var data model.RolePermissionApiEntityModel

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		data.Context = ctx
		data.RolePermissionApiEntity = payload.RolePermissionApiEntity

		result, err := s.Repository.Create(ctx, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.RolePermissionApiCreateResponse{}, err
	}

	err := redis.RedisClient.Set(ctx.Request().Context(), fmt.Sprintf("role=%d:allow=%s_%s", data.RoleID, data.ApiMethod, data.ApiPath), true, 0).Err()
	if err != nil {
		return nil, err
	}

	result := &dto.RolePermissionApiCreateResponse{
		RolePermissionApiEntityModel: data,
	}
	return result, nil
}

func (s *service) Update(ctx *abstraction.Context, payload *dto.RolePermissionApiUpdateRequest) (*dto.RolePermissionApiUpdateResponse, error) {
	var data model.RolePermissionApiEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		if _, err := s.Repository.FindByID(ctx, &payload.ID); err != nil {
			return helper.ErrorHandler(err)
		}
		data.Context = ctx
		data.RolePermissionApiEntity = payload.RolePermissionApiEntity

		result, err := s.Repository.Update(ctx, &payload.ID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.RolePermissionApiUpdateResponse{}, err
	}

	err := redis.RedisClient.Set(ctx.Request().Context(), fmt.Sprintf("role=%d:allow=%s_%s", data.RoleID, data.ApiMethod, data.ApiPath), true, 0).Err()
	if err != nil {
		return nil, err
	}

	result := &dto.RolePermissionApiUpdateResponse{
		RolePermissionApiEntityModel: data,
	}
	return result, nil
}

func (s *service) Delete(ctx *abstraction.Context, payload *dto.RolePermissionApiDeleteRequest) (*dto.RolePermissionApiDeleteResponse, error) {
	var data model.RolePermissionApiEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		rolePADelete, err := s.Repository.FindByID(ctx, &payload.ID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		err = redis.RedisClient.Del(ctx.Request().Context(), fmt.Sprintf("role=%d:allow=%s_%s", rolePADelete.RoleID, rolePADelete.ApiMethod, rolePADelete.ApiPath)).Err()
		if err != nil {
			return err
		}

		data.Context = ctx
		result, err := s.Repository.Delete(ctx, &payload.ID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.RolePermissionApiDeleteResponse{}, helper.ErrorHandler(err)
	}

	result := &dto.RolePermissionApiDeleteResponse{
		RolePermissionApiEntityModel: data,
	}
	return result, nil
}

func (s *service) RoleCacheInit(ctx *abstraction.Context) error {
	nolimit := 100000
	criteriaPaginationRolePermission := abstraction.Pagination{}
	criteriaPaginationRolePermission.PageSize = &nolimit
	rolePermissionData, _, err := s.Repository.Find(ctx, &model.RolePermissionApiFilterModel{}, &criteriaPaginationRolePermission)
	if err != nil {
		return err
	}

	redis.RedisClient.FlushAll(ctx.Request().Context())
	for _, v := range *rolePermissionData {
		// err := redis.RedisClient.Set(ctx.Request().Context(), fmt.Sprintf("role=%d:functional_id=%s:allow=%s_%s", v.RoleID, v.FunctionalID, v.ApiMethod, v.ApiPath), true, 0).Err()
		err := redis.RedisClient.Set(ctx.Request().Context(), fmt.Sprintf("role=%d:allow=%s_%s", v.RoleID, v.ApiMethod, v.ApiPath), true, 0).Err()
		if err != nil {
			return err
		}
	}

	criteriaPaginationAccessPermission := abstraction.Pagination{}
	criteriaPaginationAccessPermission.PageSize = &nolimit
	accessScopeData, _, err := s.AccessScopeRepository.Find(ctx, &model.AccessScopeFilterModel{}, &criteriaPaginationAccessPermission)
	if err != nil {
		return err
	}

	for _, v := range *accessScopeData {
		err := redis.RedisClient.Set(ctx.Request().Context(), fmt.Sprintf("access_all_company:user:%d", v.UserID), *v.AccessAll, 0).Err()
		if err != nil {
			return err
		}
		for _, vDetail := range v.AccessScopeDetail {
			err = redis.RedisClient.LPush(ctx.Request().Context(), fmt.Sprintf("access_company:user:%d", v.UserID), vDetail.CompanyID).Err()
			if err != nil {
				return err
			}
		}
	}

	return nil
}
