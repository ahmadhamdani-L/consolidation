package dto

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
)

// Get
type JelimDetailGetRequest struct {
	abstraction.Pagination
	model.JelimDetailFilterModel
}
type JelimDetailGetResponse struct {
	Datas          []model.JelimDetailEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type JelimDetailGetResponseDoc struct {
	Body struct {
		Meta res.Meta                     `json:"meta"`
		Data []model.JelimDetailEntityModel `json:"data"`
	} `json:"body"`
}

// GetByID
type JelimDetailGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type JelimDetailGetByIDResponse struct {
	model.JelimDetailEntityModel
}
type JelimDetailGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta                 `json:"meta"`
		Data JelimDetailGetByIDResponse `json:"data"`
	} `json:"body"`
}

// Create
type JelimDetailCreateRequest struct {
	ID int `param:"id" validate:"required,numeric"`
	JelimDetail []model.JelimDetailEntity
}
type JelimDetailCreateResponse struct {
	JelimDetail []model.JelimDetailEntity
}
type JelimDetailCreateResponseDoc struct {
	Body struct {
		Meta res.Meta                       `json:"meta"`
		Data JcteDetailCreateResponse `json:"data"`
	} `json:"body"`
}

// Update

// Delete
type JelimDetailDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type JelimDetailDeleteResponse struct {
	model.JelimDetailEntityModel
}
type JelimDetailDeleteResponseDoc struct {
	Body struct {
		Meta res.Meta                `json:"meta"`
		Data JelimDetailDeleteResponse `json:"data"`
	} `json:"body"`
}

// export
type JelimDetailExportRequest struct {
	model.JelimDetailFilterModel
}
type JelimDetailExportResponse struct {
	File string `json:"file"`
}
type JelimDetailExportResponseDoc struct {
	Body struct {
		Meta res.Meta                `json:"meta"`
		Data JelimDetailExportResponse `json:"data"`
	} `json:"body"`
}

// Import
type JelimDetailImportRequest struct {
	Datas []model.JelimDetailEntity
}
type JelimDetailImportResponse struct {
	Datas          []model.JelimDetailEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type JelimDetailImportResponseDoc struct {
	Body struct {
		Meta res.Meta                     `json:"meta"`
		Data []model.JelimDetailEntityModel `json:"data"`
	} `json:"body"`
}

type JelimDetailUpdateRequest struct {
	ID int `param:"id" validate:"required,numeric"`
	model.JelimDetailEntity
}
type JelimDetailUpdateResponse struct {
	model.JelimDetailEntityModel
}
type JelimDetailUpdateResponseDoc struct {
	Body struct {
		Meta res.Meta            `json:"meta"`
		Data JelimUpdateResponse `json:"data"`
	} `json:"body"`
}
