package dto

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
)

// Get
type JpmDetailGetRequest struct {
	abstraction.Pagination
	model.JpmDetailFilterModel
}
type JpmDetailGetResponse struct {
	Datas          []model.JpmDetailEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type JpmDetailGetResponseDoc struct {
	Body struct {
		Meta res.Meta                       `json:"meta"`
		Data []model.JpmDetailEntityModel `json:"data"`
	} `json:"body"`
}

// GetByID
type JpmDetailGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type JpmDetailGetByIDResponse struct {
	model.JpmDetailEntityModel
}
type JpmDetailGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta                   `json:"meta"`
		Data JpmDetailGetByIDResponse `json:"data"`
	} `json:"body"`
}

// Create
type JpmDetailCreateRequest struct {
	ID int `param:"id" validate:"required,numeric"`
	JpmDetail []model.JpmDetailEntity
}
type JpmDetailCreateResponse struct {
	JpmDetail []model.JpmDetailEntity
}
type JpmDetailCreateResponseDoc struct {
	Body struct {
		Meta res.Meta                       `json:"meta"`
		Data JcteDetailCreateResponse `json:"data"`
	} `json:"body"`
}

// Update

// Delete
type JpmDetailDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type JpmDetailDeleteResponse struct {
	model.JpmDetailEntityModel
}
type JpmDetailDeleteResponseDoc struct {
	Body struct {
		Meta res.Meta                  `json:"meta"`
		Data JpmDetailDeleteResponse `json:"data"`
	} `json:"body"`
}

// export
type JpmDetailExportRequest struct {
	model.JpmDetailFilterModel
}
type JpmDetailExportResponse struct {
	File string `json:"file"`
}
type JpmDetailExportResponseDoc struct {
	Body struct {
		Meta res.Meta                  `json:"meta"`
		Data JpmDetailExportResponse `json:"data"`
	} `json:"body"`
}

// Import
type JpmDetailImportRequest struct {
	Datas []model.JpmDetailEntity
}
type JpmDetailImportResponse struct {
	Datas          []model.JpmDetailEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type JpmDetailImportResponseDoc struct {
	Body struct {
		Meta res.Meta                       `json:"meta"`
		Data []model.JpmDetailEntityModel `json:"data"`
	} `json:"body"`
}

type JpmDetailUpdateRequest struct {
	ID int `param:"id" validate:"required,numeric"`
	model.JpmDetailEntity
}
type JpmDetailUpdateResponse struct {
	model.JpmDetailEntityModel
}
type JpmDetailUpdateResponseDoc struct {
	Body struct {
		Meta res.Meta          `json:"meta"`
		Data JpmUpdateResponse `json:"data"`
	} `json:"body"`
}
