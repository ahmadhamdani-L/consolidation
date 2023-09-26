package dto

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
)

// Get
type ImportedWorksheetDetailGetRequest struct {
	abstraction.Pagination
	model.ImportedWorksheetDetailFilterModel
}
type ImportedWorksheetDetailGetResponse struct {
	Datas          []model.ImportedWorksheetDetailEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type ImportedWorksheetDetailGetResponseDoc struct {
	Body struct {
		Meta res.Meta                                   `json:"meta"`
		Data []model.ImportedWorksheetDetailEntityModel `json:"data"`
	} `json:"body"`
}

// GetByID
type ImportedWorksheetDetailGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type ImportedWorksheetDetailGetByIDResponse struct {
	model.ImportedWorksheetDetailEntityModel
}
type ImportedWorksheetDetailGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta                               `json:"meta"`
		Data ImportedWorksheetDetailGetByIDResponse `json:"data"`
	} `json:"body"`
}

type GetTemplateRequest struct {
	Template string `query:"template"`
}
