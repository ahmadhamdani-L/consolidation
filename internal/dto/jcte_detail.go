package dto

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
)

// Get
type JcteDetailGetRequest struct {
	abstraction.Pagination
	model.JcteDetailFilterModel
}
type JcteDetailGetResponse struct {
	Datas          []model.JcteDetailEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type JcteDetailGetResponseDoc struct {
	Body struct {
		Meta res.Meta                      `json:"meta"`
		Data []model.JcteDetailEntityModel `json:"data"`
	} `json:"body"`
}

// GetByID
type JcteDetailGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type JcteDetailGetByIDResponse struct {
	model.JcteDetailEntityModel
}
type JcteDetailGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta                  `json:"meta"`
		Data JcteDetailGetByIDResponse `json:"data"`
	} `json:"body"`
}

// Create
type JcteDetailCreateRequest struct {
	ID int `param:"id" validate:"required,numeric"`
	JcteDetail []model.JcteDetailEntity
}
type JcteDetailCreateResponse struct {
	JcteDetail []model.JcteDetailEntity
}
type JcteDetailCreateResponseDoc struct {
	Body struct {
		Meta res.Meta                       `json:"meta"`
		Data JcteDetailCreateResponse `json:"data"`
	} `json:"body"`
}

// Update

// Delete
type JcteDetailDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type JcteDetailDeleteResponse struct {
	model.JcteDetailEntityModel
}
type JcteDetailDeleteResponseDoc struct {
	Body struct {
		Meta res.Meta                 `json:"meta"`
		Data JcteDetailDeleteResponse `json:"data"`
	} `json:"body"`
}

// export
type JcteDetailExportRequest struct {
	model.JcteDetailFilterModel
}
type JcteDetailExportResponse struct {
	File string `json:"file"`
}
type JcteDetailExportResponseDoc struct {
	Body struct {
		Meta res.Meta                 `json:"meta"`
		Data JcteDetailExportResponse `json:"data"`
	} `json:"body"`
}

// Import
type JcteDetailImportRequest struct {
	Datas []model.JcteDetailEntity
}
type JcteDetailImportResponse struct {
	Datas          []model.JcteDetailEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type JcteDetailImportResponseDoc struct {
	Body struct {
		Meta res.Meta                      `json:"meta"`
		Data []model.JcteDetailEntityModel `json:"data"`
	} `json:"body"`
}

// Update
type JcteDetailUpdateRequest struct {
	ID int `param:"id" validate:"required,numeric"`
	model.JcteDetailEntity
}
type JcteDetailUpdateResponse struct {
	model.JcteDetailEntityModel
}
type JcteDetailUpdateResponseDoc struct {
	Body struct {
		Meta res.Meta           `json:"meta"`
		Data JcteUpdateResponse `json:"data"`
	} `json:"body"`
}
