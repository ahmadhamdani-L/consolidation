package permissiondef

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
	Repository repository.PermissionDef
	Db         *gorm.DB
}

type Service interface {
	Find(ctx *abstraction.Context, payload *dto.PermissionDefGetRequest) (*dto.PermissionDefGetResponse, error)
	FindByID(ctx *abstraction.Context, payload *dto.PermissionDefGetByIDRequest) (*dto.PermissionDefGetByIDResponse, error)
	Create(ctx *abstraction.Context, payload *dto.PermissionDefCreateRequest) (*dto.PermissionDefCreateResponse, error)
	Update(ctx *abstraction.Context, payload *dto.PermissionDefUpdateRequest) (*dto.PermissionDefUpdateResponse, error)
	Delete(ctx *abstraction.Context, payload *dto.PermissionDefDeleteRequest) (*dto.PermissionDefDeleteResponse, error)
}

func NewService(f *factory.Factory) *service {
	repository := f.PermissionDefRepository
	db := f.Db
	return &service{
		Repository: repository,
		Db:         db,
	}
}

func (s *service) Find(ctx *abstraction.Context, payload *dto.PermissionDefGetRequest) (*dto.PermissionDefGetResponse, error) {
	data, info, err := s.Repository.Find(ctx, &payload.PermissionDefFilterModel, &payload.Pagination)
	if err != nil {
		return &dto.PermissionDefGetResponse{}, helper.ErrorHandler(err)
	}

	result := &dto.PermissionDefGetResponse{
		Datas:          *data,
		PaginationInfo: *info,
	}
	return result, nil
}

func (s *service) FindByID(ctx *abstraction.Context, payload *dto.PermissionDefGetByIDRequest) (*dto.PermissionDefGetByIDResponse, error) {
	data, err := s.Repository.FindByID(ctx, &payload.ID)
	if err != nil {
		return &dto.PermissionDefGetByIDResponse{}, helper.ErrorHandler(err)
	}
	result := &dto.PermissionDefGetByIDResponse{
		PermissionDefEntityModel: *data,
	}
	return result, nil
}

func (s *service) Create(ctx *abstraction.Context, payload *dto.PermissionDefCreateRequest) (*dto.PermissionDefCreateResponse, error) {
	var data model.PermissionDefEntityModel

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		data.Context = ctx
		data.PermissionDefEntity = payload.PermissionDefEntity

		result, err := s.Repository.Create(ctx, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.PermissionDefCreateResponse{}, err
	}

	result := &dto.PermissionDefCreateResponse{
		PermissionDefEntityModel: data,
	}
	return result, nil
}

func (s *service) Update(ctx *abstraction.Context, payload *dto.PermissionDefUpdateRequest) (*dto.PermissionDefUpdateResponse, error) {
	var data model.PermissionDefEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		if _, err := s.Repository.FindByID(ctx, &payload.ID); err != nil {
			return helper.ErrorHandler(err)
		}
		data.Context = ctx
		data.PermissionDefEntity = payload.PermissionDefEntity

		result, err := s.Repository.Update(ctx, &payload.ID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.PermissionDefUpdateResponse{}, err
	}
	result := &dto.PermissionDefUpdateResponse{
		PermissionDefEntityModel: data,
	}
	return result, nil
}

func (s *service) Delete(ctx *abstraction.Context, payload *dto.PermissionDefDeleteRequest) (*dto.PermissionDefDeleteResponse, error) {
	var data model.PermissionDefEntityModel
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
		return &dto.PermissionDefDeleteResponse{}, err
	}
	result := &dto.PermissionDefDeleteResponse{
		PermissionDefEntityModel: data,
	}
	return result, nil
}
