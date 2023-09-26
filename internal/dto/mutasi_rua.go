package dto

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
	"mime/multipart"
)

// Get
type MutasiRuaGetRequest struct {
	abstraction.Pagination
	model.MutasiRuaFilterModel
}
type MutasiRuaGetResponse struct {
	Datas          []model.MutasiRuaEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type MutasiRuaGetResponseDoc struct {
	Body struct {
		Meta res.Meta                     `json:"meta"`
		Data []model.MutasiRuaEntityModel `json:"data"`
	} `json:"body"`
}

// GetByID
type MutasiRuaGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type MutasiRuaGetByIDResponse struct {
	model.MutasiRuaEntityModel
}
type MutasiRuaGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta                 `json:"meta"`
		Data MutasiRuaGetByIDResponse `json:"data"`
	} `json:"body"`
}

// Create
type MutasiRuaCreateRequest struct {
	model.MutasiRuaEntity
}
type MutasiRuaCreateResponse struct {
	model.MutasiRuaEntityModel
}
type MutasiRuaCreateResponseDoc struct {
	Body struct {
		Meta res.Meta                `json:"meta"`
		Data MutasiRuaCreateResponse `json:"data"`
	} `json:"body"`
}

// Update
type MutasiRuaUpdateRequest struct {
	ID int `param:"id" validate:"required,numeric"`
	model.MutasiRuaEntity
}
type MutasiRuaUpdateResponse struct {
	model.MutasiRuaEntityModel
}
type MutasiRuaUpdateResponseDoc struct {
	Body struct {
		Meta res.Meta                `json:"meta"`
		Data MutasiRuaUpdateResponse `json:"data"`
	} `json:"body"`
}

// Delete
type MutasiRuaDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type MutasiRuaDeleteResponse struct {
	// model.MutasiRuaEntityModel
}
type MutasiRuaDeleteResponseDoc struct {
	Body struct {
		Meta res.Meta                `json:"meta"`
		Data MutasiRuaDeleteResponse `json:"data"`
	} `json:"body"`
}

// export
type MutasiRuaExportRequest struct {
	// UserID    int
	// Period    string `query:"period" validate:"required"`
	// Versions  int    `query:"versions" validate:"required"`
	// CompanyID int    `query:"company_id"`
	MutasiRuaID int `query:"mutasi_rua_id" validate:"required"`
}
type MutasiRuaExportResponse struct {
	FileName string `json:"filename"`
	Path     string `json:"path"`
}
type MutasiRuaExportResponseDoc struct {
	Body struct {
		Meta res.Meta                `json:"meta"`
		Data MutasiRuaExportResponse `json:"data"`
	} `json:"body"`
}

// Import
type MutasiRuaImportRequest struct {
	UserID    int
	CompanyID int
	File      multipart.File
}
type MutasiRuaImportResponse struct {
	Data model.MutasiRuaEntityModel
}
type MutasiRuaImportResponseDoc struct {
	Body struct {
		Meta res.Meta                `json:"meta"`
		Data MutasiRuaImportResponse `json:"data"`
	} `json:"body"`
}
