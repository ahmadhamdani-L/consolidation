package importedworksheetdetail

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/dto"
	"mcash-finance-console-core/internal/factory"
	"mcash-finance-console-core/internal/repository"
	"mcash-finance-console-core/pkg/util/response"

	"gorm.io/gorm"
)

type service struct {
	Repository                  repository.ImportedWorksheetDetail
	ImportedWorksheetRepository repository.ImportedWorksheet
	Db                          *gorm.DB
}

type Service interface {
	Find(ctx *abstraction.Context, payload *dto.ImportedWorksheetDetailGetRequest) (*dto.ImportedWorksheetDetailGetResponse, error)
	FindByID(ctx *abstraction.Context, payload *dto.ImportedWorksheetDetailGetByIDRequest) (*dto.ImportedWorksheetDetailGetByIDResponse, error)
}

func NewService(f *factory.Factory) *service {
	repository := f.ImportedWorksheetDetailRepository
	ImportedWorksheetRepository := f.ImportedWorksheetRepository
	db := f.Db
	return &service{
		Repository:                  repository,
		ImportedWorksheetRepository: ImportedWorksheetRepository,
		Db:                          db,
	}
}

func (s *service) Find(ctx *abstraction.Context, payload *dto.ImportedWorksheetDetailGetRequest) (*dto.ImportedWorksheetDetailGetResponse, error) {
	data, info, err := s.Repository.Find(ctx, &payload.ImportedWorksheetDetailFilterModel, &payload.Pagination)
	if err != nil {
		return &dto.ImportedWorksheetDetailGetResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}

	result := &dto.ImportedWorksheetDetailGetResponse{
		Datas:          *data,
		PaginationInfo: *info,
	}
	return result, nil
}

func (s *service) FindByID(ctx *abstraction.Context, payload *dto.ImportedWorksheetDetailGetByIDRequest) (*dto.ImportedWorksheetDetailGetByIDResponse, error) {
	data, err := s.Repository.FindByID(ctx, &payload.ID)
	if err != nil {
		return &dto.ImportedWorksheetDetailGetByIDResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	result := &dto.ImportedWorksheetDetailGetByIDResponse{
		ImportedWorksheetDetailEntityModel: *data,
	}
	return result, nil
}
