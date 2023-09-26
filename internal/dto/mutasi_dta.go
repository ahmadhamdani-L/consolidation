package dto

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
	"mime/multipart"
)

// Get
type MutasiDtaGetRequest struct {
	abstraction.Pagination
	model.MutasiDtaFilterModel
}
type MutasiDtaGetResponse struct {
	Datas          []model.MutasiDtaEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type MutasiDtaGetResponseDoc struct {
	Body struct {
		Meta res.Meta                     `json:"meta"`
		Data []model.MutasiDtaEntityModel `json:"data"`
	} `json:"body"`
}

// GetByID
type MutasiDtaGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type MutasiDtaGetByIDResponse struct {
	model.MutasiDtaEntityModel
}
type MutasiDtaGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta                 `json:"meta"`
		Data MutasiDtaGetByIDResponse `json:"data"`
	} `json:"body"`
}

// Create
type MutasiDtaCreateRequest struct {
	model.MutasiDtaEntity
}
type MutasiDtaCreateResponse struct {
	model.MutasiDtaEntityModel
}
type MutasiDtaCreateResponseDoc struct {
	Body struct {
		Meta res.Meta                `json:"meta"`
		Data MutasiDtaCreateResponse `json:"data"`
	} `json:"body"`
}

// Update
type MutasiDtaUpdateRequest struct {
	ID int `param:"id" validate:"required,numeric"`
	model.MutasiDtaEntity
}
type MutasiDtaUpdateResponse struct {
	model.MutasiDtaEntityModel
}
type MutasiDtaUpdateResponseDoc struct {
	Body struct {
		Meta res.Meta                `json:"meta"`
		Data MutasiDtaUpdateResponse `json:"data"`
	} `json:"body"`
}

// Delete
type MutasiDtaDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type MutasiDtaDeleteResponse struct {
	// model.MutasiDtaEntityModel
}
type MutasiDtaDeleteResponseDoc struct {
	Body struct {
		Meta res.Meta                `json:"meta"`
		Data MutasiDtaDeleteResponse `json:"data"`
	} `json:"body"`
}

// export
type MutasiDtaExportRequest struct {
	// UserID    int
	// Period    string `query:"period" validate:"required"`
	// Versions  int    `query:"versions" validate:"required"`
	// CompanyID int    `query:"company_id"`
	MutasiDtaID int `query:"mutasi_dta_id" validate:"required"`
}
type MutasiDtaExportResponse struct {
	FileName string `json:"filename"`
	Path     string `json:"path"`
}
type MutasiDtaExportResponseDoc struct {
	Body struct {
		Meta res.Meta                `json:"meta"`
		Data MutasiDtaExportResponse `json:"data"`
	} `json:"body"`
}

// Import
type MutasiDtaImportRequest struct {
	UserID    int
	CompanyID int
	File      multipart.File
}
type MutasiDtaImportResponse struct {
	Data model.MutasiDtaEntityModel
}
type MutasiDtaImportResponseDoc struct {
	Body struct {
		Meta res.Meta                `json:"meta"`
		Data MutasiDtaImportResponse `json:"data"`
	} `json:"body"`
}
