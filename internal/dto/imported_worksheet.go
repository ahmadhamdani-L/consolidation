package dto

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
)

type ImportedWorksheetCreateRequest struct {
	model.ImportedWorksheetEntityModel
}
type ImportedWorksheetCreateResponse struct {
	model.ImportedWorksheetEntityModel
}
type ImportedWorksheetCreateResponseDoc struct {
	Body struct {
		Meta res.Meta                        `json:"meta"`
		Data ImportedWorksheetCreateResponse `json:"data"`
	} `json:"body"`
}

// Get
type ImportedWorksheetGetRequest struct {
	abstraction.Pagination
	model.ImportedWorksheetFilterModel
}
type ImportedWorksheetGetResponse struct {
	Datas          []model.ImportedWorksheetEntityModel
	PaginationInfo abstraction.PaginationInfo
	Description    int
}
type ImportedWorksheetGetResponseDoc struct {
	Body struct {
		Meta res.Meta                             `json:"meta"`
		Data []model.ImportedWorksheetEntityModel `json:"data"`
	} `json:"body"`
}

// GetByID
type ImportedWorksheetGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type ImportedWorksheetGetByIDResponse struct {
	model.ImportedWorksheetEntityModel
}
type ImportedWorksheetGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta                         `json:"meta"`
		Data ImportedWorksheetGetByIDResponse `json:"data"`
	} `json:"body"`
}

type ImportedWorksheetGetByIDDownloadAllResponse struct {
	Datas []model.ImportedWorksheetDetailEntityModel
	FileName []string
}

type CompanyGetByIDBulkRequest struct {
	ID int `param:"id" validate:"required,numeric"`
	model.CompanyEntityModel
}
type CompanyGetByIDBulkResponse struct {
	model.CompanyEntityModel
}
type CompanyGetByIDBulkResponseDoc struct {
	Body struct {
		Meta res.Meta               `json:"meta"`
		Data CompanyGetByIDResponse `json:"data"`
	} `json:"body"`
}

type ImportJurnal struct {
	TbID int `json:"tb_id" validate:"required,numeric"`
	DataJurnal []model.AdjustmentDetailEntity
}

// Delete
type ImportedWorksheetDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type ImportedWorksheetDeleteResponse struct {
	// model.EmployeeBenefitEntityModel
}