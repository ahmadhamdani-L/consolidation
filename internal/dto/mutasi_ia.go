package dto

import (
	"worker/internal/abstraction"
	"worker/internal/model"
	res "worker/pkg/util/response"
	"mime/multipart"
)

// Get
type MutasiIaGetRequest struct {
	abstraction.Pagination
	model.MutasiIaFilterModel
}
type MutasiIaGetResponse struct {
	Datas          []model.MutasiIaEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type MutasiIaGetResponseDoc struct {
	Body struct {
		Meta res.Meta                    `json:"meta"`
		Data []model.MutasiIaEntityModel `json:"data"`
	} `json:"body"`
}

// GetByID
type MutasiIaGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type MutasiIaGetByIDResponse struct {
	model.MutasiIaEntityModel
}
type MutasiIaGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta                `json:"meta"`
		Data MutasiIaGetByIDResponse `json:"data"`
	} `json:"body"`
}

// Create
type MutasiIaCreateRequest struct {
	model.MutasiIaEntity
}
type MutasiIaCreateResponse struct {
	model.MutasiIaEntityModel
}
type MutasiIaCreateResponseDoc struct {
	Body struct {
		Meta res.Meta               `json:"meta"`
		Data MutasiIaCreateResponse `json:"data"`
	} `json:"body"`
}

// Update
type MutasiIaUpdateRequest struct {
	ID int `param:"id" validate:"required,numeric"`
	model.MutasiIaEntity
}
type MutasiIaUpdateResponse struct {
	model.MutasiIaEntityModel
}
type MutasiIaUpdateResponseDoc struct {
	Body struct {
		Meta res.Meta               `json:"meta"`
		Data MutasiIaUpdateResponse `json:"data"`
	} `json:"body"`
}

// Delete
type MutasiIaDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type MutasiIaDeleteResponse struct {
	model.MutasiIaEntityModel
}
type MutasiIaDeleteResponseDoc struct {
	Body struct {
		Meta res.Meta               `json:"meta"`
		Data MutasiIaDeleteResponse `json:"data"`
	} `json:"body"`
}

// export
type MutasiIaExportRequest struct {
	UserId    int
	Period    string `query:"period" validate:"required"`
	Versions  int    `query:"versions" validate:"required"`
	CompanyID int    `query:"company_id"`
}
type MutasiIaExportResponse struct {
	File string `json:"file"`
}
type MutasiIaExportResponseDoc struct {
	Body struct {
		Meta res.Meta               `json:"meta"`
		Data MutasiIaExportResponse `json:"data"`
	} `json:"body"`
}

// Import
type MutasiIaImportRequest struct {
	UserId    int
	CompanyId int
	File      multipart.File
}
type MutasiIaImportResponse struct {
	Datas          []model.MutasiIaEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type MutasiIaImportResponseDoc struct {
	Body struct {
		Meta res.Meta                    `json:"meta"`
		Data []model.MutasiIaEntityModel `json:"data"`
	} `json:"body"`
}
