package dto

import (
	"worker/internal/abstraction"
	"worker/internal/model"
	res "worker/pkg/util/response"
)

// Get
type MutasiFaDetailGetRequest struct {
	abstraction.Pagination
	model.MutasiFaDetailFilterModel
}
type MutasiFaDetailGetResponse struct {
	Datas          []model.MutasiFaDetailEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type MutasiFaDetailGetResponseDoc struct {
	Body struct {
		Meta res.Meta                          `json:"meta"`
		Data []model.MutasiFaDetailEntityModel `json:"data"`
	} `json:"body"`
}

// GetByID
type MutasiFaDetailGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type MutasiFaDetailGetByIDResponse struct {
	model.MutasiFaDetailEntityModel
}
type MutasiFaDetailGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta                      `json:"meta"`
		Data MutasiFaDetailGetByIDResponse `json:"data"`
	} `json:"body"`
}

// Create
type MutasiFaDetailCreateRequest struct {
	model.MutasiFaDetailEntity
}
type MutasiFaDetailCreateResponse struct {
	model.MutasiFaDetailEntityModel
}
type MutasiFaDetailCreateResponseDoc struct {
	Body struct {
		Meta res.Meta                     `json:"meta"`
		Data MutasiFaDetailCreateResponse `json:"data"`
	} `json:"body"`
}

// Update
type MutasiFaDetailUpdateRequest struct {
	ID int `param:"id" validate:"required,numeric"`
	model.MutasiFaDetailEntity
}
type MutasiFaDetailUpdateResponse struct {
	model.MutasiFaDetailEntityModel
}
type MutasiFaDetailUpdateResponseDoc struct {
	Body struct {
		Meta res.Meta                     `json:"meta"`
		Data MutasiFaDetailUpdateResponse `json:"data"`
	} `json:"body"`
}

// Delete
type MutasiFaDetailDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type MutasiFaDetailDeleteResponse struct {
	model.MutasiFaDetailEntityModel
}
type MutasiFaDetailDeleteResponseDoc struct {
	Body struct {
		Meta res.Meta                     `json:"meta"`
		Data MutasiFaDetailDeleteResponse `json:"data"`
	} `json:"body"`
}

// export
type MutasiFaDetailExportRequest struct {
	model.MutasiFaDetailFilterModel
}
type MutasiFaDetailExportResponse struct {
	File string `json:"file"`
}
type MutasiFaDetailExportResponseDoc struct {
	Body struct {
		Meta res.Meta                     `json:"meta"`
		Data MutasiFaDetailExportResponse `json:"data"`
	} `json:"body"`
}

// Import
type MutasiFaDetailImportRequest struct {
	Datas []model.MutasiFaDetailEntity
}
type MutasiFaDetailImportResponse struct {
	Datas          []model.MutasiFaDetailEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type MutasiFaDetailImportResponseDoc struct {
	Body struct {
		Meta res.Meta                          `json:"meta"`
		Data []model.MutasiFaDetailEntityModel `json:"data"`
	} `json:"body"`
}
