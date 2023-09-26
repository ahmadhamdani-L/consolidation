package coatype

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
	Repository repository.CoaType
	Db         *gorm.DB
}

type Service interface {
	Find(ctx *abstraction.Context, payload *dto.CoaTypeGetRequest) (*dto.CoaTypeGetResponse, error)
	FindByID(ctx *abstraction.Context, payload *dto.CoaTypeGetByIDRequest) (*dto.CoaTypeGetByIDResponse, error)
	Create(ctx *abstraction.Context, payload *dto.CoaTypeCreateRequest) (*dto.CoaTypeCreateResponse, error)
	Update(ctx *abstraction.Context, payload *dto.CoaTypeUpdateRequest) (*dto.CoaTypeUpdateResponse, error)
	Delete(ctx *abstraction.Context, payload *dto.CoaTypeDeleteRequest) (*dto.CoaTypeDeleteResponse, error)
}

func NewService(f *factory.Factory) *service {
	repository := f.CoaTypeRepository
	db := f.Db
	return &service{
		Repository: repository,
		Db:         db,
	}
}

func (s *service) Find(ctx *abstraction.Context, payload *dto.CoaTypeGetRequest) (*dto.CoaTypeGetResponse, error) {
	data, info, err := s.Repository.Find(ctx, &payload.CoaTypeFilterModel, &payload.Pagination)
	if err != nil {
		return &dto.CoaTypeGetResponse{}, helper.ErrorHandler(err)
	}
	result := &dto.CoaTypeGetResponse{
		Datas:          *data,
		PaginationInfo: *info,
	}
	return result, nil
}

func (s *service) FindByID(ctx *abstraction.Context, payload *dto.CoaTypeGetByIDRequest) (*dto.CoaTypeGetByIDResponse, error) {
	data, err := s.Repository.FindByID(ctx, &payload.ID)
	if err != nil {
		return &dto.CoaTypeGetByIDResponse{}, helper.ErrorHandler(err)
	}
	result := &dto.CoaTypeGetByIDResponse{
		CoaTypeEntityModel: *data,
	}
	return result, nil
}

func (s *service) Create(ctx *abstraction.Context, payload *dto.CoaTypeCreateRequest) (*dto.CoaTypeCreateResponse, error) {
	var data model.CoaTypeEntityModel

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		data.Context = ctx
		data.CoaTypeEntity = payload.CoaTypeEntity
		result, err := s.Repository.Create(ctx, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.CoaTypeCreateResponse{}, err
	}
	result := &dto.CoaTypeCreateResponse{
		CoaTypeEntityModel: data,
	}
	return result, nil
}

func (s *service) Update(ctx *abstraction.Context, payload *dto.CoaTypeUpdateRequest) (*dto.CoaTypeUpdateResponse, error) {
	var data model.CoaTypeEntityModel

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		if _, err := s.Repository.FindByID(ctx, &payload.ID); err != nil {
			return helper.ErrorHandler(err)
		}
		data.Context = ctx
		data.CoaTypeEntity = payload.CoaTypeEntity
		result, err := s.Repository.Update(ctx, &payload.ID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.CoaTypeUpdateResponse{}, err
	}
	result := &dto.CoaTypeUpdateResponse{
		CoaTypeEntityModel: data,
	}
	return result, nil
}

func (s *service) Delete(ctx *abstraction.Context, payload *dto.CoaTypeDeleteRequest) (*dto.CoaTypeDeleteResponse, error) {
	var data model.CoaTypeEntityModel
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
		return &dto.CoaTypeDeleteResponse{}, err
	}
	result := &dto.CoaTypeDeleteResponse{
		// CoaTypeEntityModel: data,
	}
	return result, nil
}
