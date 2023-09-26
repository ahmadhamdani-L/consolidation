package parameter

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
	Repository repository.Parameter
	Db         *gorm.DB
}

type Service interface {
	Find(ctx *abstraction.Context, payload *dto.ParameterGetRequest) (*dto.ParameterGetResponse, error)
	FindByID(ctx *abstraction.Context, payload *dto.ParameterGetByIDRequest) (*dto.ParameterGetByIDResponse, error)
	Create(ctx *abstraction.Context, payload *dto.ParameterCreateRequest) (*dto.ParameterCreateResponse, error)
	Update(ctx *abstraction.Context, payload *dto.ParameterUpdateRequest) (*dto.ParameterUpdateResponse, error)
	Delete(ctx *abstraction.Context, payload *dto.ParameterDeleteRequest) (*dto.ParameterDeleteResponse, error)
}

func NewService(f *factory.Factory) *service {
	repository := f.ParameterRepository
	db := f.Db
	return &service{
		Repository: repository,
		Db:         db,
	}
}

func (s *service) Find(ctx *abstraction.Context, payload *dto.ParameterGetRequest) (*dto.ParameterGetResponse, error) {
	data, info, err := s.Repository.Find(ctx, &payload.ParameterFilterModel, &payload.Pagination)
	if err != nil {
		return &dto.ParameterGetResponse{}, helper.ErrorHandler(err)
	}
	result := &dto.ParameterGetResponse{
		Datas:          *data,
		PaginationInfo: *info,
	}
	return result, nil
}

func (s *service) FindByID(ctx *abstraction.Context, payload *dto.ParameterGetByIDRequest) (*dto.ParameterGetByIDResponse, error) {
	data, err := s.Repository.FindByID(ctx, &payload.ID)
	if err != nil {
		return &dto.ParameterGetByIDResponse{}, helper.ErrorHandler(err)
	}
	result := &dto.ParameterGetByIDResponse{
		ParameterEntityModel: *data,
	}
	return result, nil
}

func (s *service) Create(ctx *abstraction.Context, payload *dto.ParameterCreateRequest) (*dto.ParameterCreateResponse, error) {
	var data model.ParameterEntityModel

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		data.Context = ctx
		data.ParameterEntity = payload.ParameterEntity
		result, err := s.Repository.Create(ctx, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.ParameterCreateResponse{}, err
	}
	result := &dto.ParameterCreateResponse{
		ParameterEntityModel: data,
	}
	return result, nil
}

func (s *service) Update(ctx *abstraction.Context, payload *dto.ParameterUpdateRequest) (*dto.ParameterUpdateResponse, error) {
	var data model.ParameterEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		if _, err := s.Repository.FindByID(ctx, &payload.ID); err != nil {
			return helper.ErrorHandler(err)
		}
		data.Context = ctx
		data.ParameterEntity = payload.ParameterEntity
		result, err := s.Repository.Update(ctx, &payload.ID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.ParameterUpdateResponse{}, err
	}
	result := &dto.ParameterUpdateResponse{
		ParameterEntityModel: data,
	}
	return result, nil
}

func (s *service) Delete(ctx *abstraction.Context, payload *dto.ParameterDeleteRequest) (*dto.ParameterDeleteResponse, error) {
	var data model.ParameterEntityModel
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
		return &dto.ParameterDeleteResponse{}, err
	}
	result := &dto.ParameterDeleteResponse{
		// ParameterEntityModel: data,
	}
	return result, nil
}
