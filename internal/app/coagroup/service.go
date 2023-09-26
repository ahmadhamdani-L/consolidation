package coagroup

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
	Repository repository.CoaGroup
	Db         *gorm.DB
}

type Service interface {
	Find(ctx *abstraction.Context, payload *dto.CoaGroupGetRequest) (*dto.CoaGroupGetResponse, error)
	FindByID(ctx *abstraction.Context, payload *dto.CoaGroupGetByIDRequest) (*dto.CoaGroupGetByIDResponse, error)
	Create(ctx *abstraction.Context, payload *dto.CoaGroupCreateRequest) (*dto.CoaGroupCreateResponse, error)
	Update(ctx *abstraction.Context, payload *dto.CoaGroupUpdateRequest) (*dto.CoaGroupUpdateResponse, error)
	Delete(ctx *abstraction.Context, payload *dto.CoaGroupDeleteRequest) (*dto.CoaGroupDeleteResponse, error)
}

func NewService(f *factory.Factory) *service {
	repository := f.CoaGroupRepository
	db := f.Db
	return &service{
		Repository: repository,
		Db:         db,
	}
}

func (s *service) Find(ctx *abstraction.Context, payload *dto.CoaGroupGetRequest) (*dto.CoaGroupGetResponse, error) {
	data, info, err := s.Repository.Find(ctx, &payload.CoaGroupFilterModel, &payload.Pagination)
	if err != nil {
		return &dto.CoaGroupGetResponse{}, helper.ErrorHandler(err)
	}
	result := &dto.CoaGroupGetResponse{
		Datas:          *data,
		PaginationInfo: *info,
	}
	return result, nil
}

func (s *service) FindByID(ctx *abstraction.Context, payload *dto.CoaGroupGetByIDRequest) (*dto.CoaGroupGetByIDResponse, error) {
	data, err := s.Repository.FindByID(ctx, &payload.ID)
	if err != nil {
		return &dto.CoaGroupGetByIDResponse{}, helper.ErrorHandler(err)
	}
	result := &dto.CoaGroupGetByIDResponse{
		CoaGroupEntityModel: *data,
	}
	return result, nil
}

func (s *service) Create(ctx *abstraction.Context, payload *dto.CoaGroupCreateRequest) (*dto.CoaGroupCreateResponse, error) {
	var data model.CoaGroupEntityModel

	if err = trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		data.Context = ctx
		data.CoaGroupEntity = payload.CoaGroupEntity
		result, err := s.Repository.Create(ctx, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.CoaGroupCreateResponse{}, err
	}
	result := &dto.CoaGroupCreateResponse{
		CoaGroupEntityModel: data,
	}
	return result, nil
}

func (s *service) Update(ctx *abstraction.Context, payload *dto.CoaGroupUpdateRequest) (*dto.CoaGroupUpdateResponse, error) {
	var data model.CoaGroupEntityModel

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		if _, err := s.Repository.FindByID(ctx, &payload.ID); err != nil {
			return helper.ErrorHandler(err)
		}
		data.Context = ctx
		data.CoaGroupEntity = payload.CoaGroupEntity
		result, err := s.Repository.Update(ctx, &payload.ID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.CoaGroupUpdateResponse{}, err
	}
	result := &dto.CoaGroupUpdateResponse{
		CoaGroupEntityModel: data,
	}
	return result, nil
}

func (s *service) Delete(ctx *abstraction.Context, payload *dto.CoaGroupDeleteRequest) (*dto.CoaGroupDeleteResponse, error) {
	var data model.CoaGroupEntityModel
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
		return &dto.CoaGroupDeleteResponse{}, err
	}
	result := &dto.CoaGroupDeleteResponse{
		// CoaGroupEntityModel: data,
	}
	return result, nil
}
