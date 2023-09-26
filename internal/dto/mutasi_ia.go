package dto

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
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
	// model.MutasiIaEntityModel
}
type MutasiIaDeleteResponseDoc struct {
	Body struct {
		Meta res.Meta               `json:"meta"`
		Data MutasiIaDeleteResponse `json:"data"`
	} `json:"body"`
}

// export
type MutasiIaExportRequest struct {
	// UserID    int
	// Period    string `query:"period" validate:"required"`
	// Versions  int    `query:"versions" validate:"required"`
	// CompanyID int    `query:"company_id"`
	MutasiIaID int `query:"mutasi_ia_id" validate:"required"`
}
type MutasiIaExportResponse struct {
	FileName string `json:"filename"`
	Path     string `json:"path"`
}
type MutasiIaExportResponseDoc struct {
	Body struct {
		Meta res.Meta               `json:"meta"`
		Data MutasiIaExportResponse `json:"data"`
	} `json:"body"`
}

// Import
type MutasiIaImportRequest struct {
	UserID    int
	CompanyID int
	File      multipart.File
}
type MutasiIaImportResponse struct {
	Data model.MutasiIaEntityModel
}
type MutasiIaImportResponseDoc struct {
	Body struct {
		Meta res.Meta               `json:"meta"`
		Data MutasiIaImportResponse `json:"data"`
	} `json:"body"`
}
