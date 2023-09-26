package rolepermission

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/dto"
	"mcash-finance-console-core/internal/factory"
	"mcash-finance-console-core/internal/model"
	"mcash-finance-console-core/internal/repository"
	"mcash-finance-console-core/pkg/util/helper"
	"mcash-finance-console-core/pkg/util/trxmanager"

	"gorm.io/gorm"
)

type service struct {
	Repository repository.RolePermission
	Db         *gorm.DB
}

type Service interface {
	Find(ctx *abstraction.Context, payload *dto.RolePermissionGetRequest) (*dto.RolePermissionGetResponse, error)
	FindByID(ctx *abstraction.Context, payload *dto.RolePermissionGetByIDRequest) (*dto.RolePermissionGetByIDResponse, error)
	Create(ctx *abstraction.Context, payload *dto.RolePermissionCreateRequest) (*dto.RolePermissionCreateResponse, error)
	Update(ctx *abstraction.Context, payload *dto.RolePermissionUpdateRequest) (*dto.RolePermissionUpdateResponse, error)
	Delete(ctx *abstraction.Context, payload *dto.RolePermissionDeleteRequest) (*dto.RolePermissionDeleteResponse, error)
}

func NewService(f *factory.Factory) *service {
	repository := f.RolePermissionRepository
	db := f.Db
	return &service{
		Repository: repository,
		Db:         db,
	}
}

func (s *service) Find(ctx *abstraction.Context, payload *dto.RolePermissionGetRequest) (*dto.RolePermissionGetResponse, error) {
	data, info, err := s.Repository.Find(ctx, &payload.RolePermissionFilterModel, &payload.Pagination)
	if err != nil {
		return &dto.RolePermissionGetResponse{}, helper.ErrorHandler(err)
	}

	result := &dto.RolePermissionGetResponse{
		Datas:          *data,
		PaginationInfo: *info,
	}
	return result, nil
}

func (s *service) FindByID(ctx *abstraction.Context, payload *dto.RolePermissionGetByIDRequest) (*dto.RolePermissionGetByIDResponse, error) {
	data, err := s.Repository.FindByID(ctx, &payload.ID)
	if err != nil {
		return &dto.RolePermissionGetByIDResponse{}, helper.ErrorHandler(err)
	}
	result := &dto.RolePermissionGetByIDResponse{
		RolePermissionEntityModel: *data,
	}
	return result, nil
}

func (s *service) Create(ctx *abstraction.Context, payload *dto.RolePermissionCreateRequest) (*dto.RolePermissionCreateResponse, error) {
	var data model.RolePermissionEntityModel

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		data.Context = ctx
		data.RolePermissionEntity = payload.RolePermissionEntity

		result, err := s.Repository.Create(ctx, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.RolePermissionCreateResponse{}, err
	}

	result := &dto.RolePermissionCreateResponse{
		RolePermissionEntityModel: data,
	}
	return result, nil
}

func (s *service) Update(ctx *abstraction.Context, payload *dto.RolePermissionUpdateRequest) (*dto.RolePermissionUpdateResponse, error) {
	var data model.RolePermissionEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		if _, err := s.Repository.FindByID(ctx, &payload.ID); err != nil {
			return helper.ErrorHandler(err)
		}
		data.Context = ctx
		data.RolePermissionEntity = payload.RolePermissionEntity

		result, err := s.Repository.Update(ctx, &payload.ID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.RolePermissionUpdateResponse{}, err
	}
	result := &dto.RolePermissionUpdateResponse{
		RolePermissionEntityModel: data,
	}
	return result, nil
}

func (s *service) Delete(ctx *abstraction.Context, payload *dto.RolePermissionDeleteRequest) (*dto.RolePermissionDeleteResponse, error) {
	var data model.RolePermissionEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		if _, err := s.Repository.FindByID(ctx, &payload.ID); err != nil {
			return helper.ErrorHandler(err)
		}
		data.Context = ctx
		result, err := s.Repository.Delete(ctx, &payload.ID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.RolePermissionDeleteResponse{}, err
	}
	result := &dto.RolePermissionDeleteResponse{
		RolePermissionEntityModel: data,
	}
	return result, nil
}
