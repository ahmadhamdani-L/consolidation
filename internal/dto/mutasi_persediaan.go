package dto

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
	"mime/multipart"
)

// Get
type MutasiPersediaanGetRequest struct {
	abstraction.Pagination
	model.MutasiPersediaanFilterModel
}
type MutasiPersediaanGetResponse struct {
	Datas          []model.MutasiPersediaanEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type MutasiPersediaanGetResponseDoc struct {
	Body struct {
		Meta res.Meta                            `json:"meta"`
		Data []model.MutasiPersediaanEntityModel `json:"data"`
	} `json:"body"`
}

// GetByID
type MutasiPersediaanGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type MutasiPersediaanGetByIDResponse struct {
	model.MutasiPersediaanEntityModel
}
type MutasiPersediaanGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta                        `json:"meta"`
		Data MutasiPersediaanGetByIDResponse `json:"data"`
	} `json:"body"`
}

// Create
type MutasiPersediaanCreateRequest struct {
	model.MutasiPersediaanEntity
}
type MutasiPersediaanCreateResponse struct {
	model.MutasiPersediaanEntityModel
}
type MutasiPersediaanCreateResponseDoc struct {
	Body struct {
		Meta res.Meta                       `json:"meta"`
		Data MutasiPersediaanCreateResponse `json:"data"`
	} `json:"body"`
}

// Update
type MutasiPersediaanUpdateRequest struct {
	ID int `param:"id" validate:"required,numeric"`
	model.MutasiPersediaanEntity
}
type MutasiPersediaanUpdateResponse struct {
	model.MutasiPersediaanEntityModel
}
type MutasiPersediaanUpdateResponseDoc struct {
	Body struct {
		Meta res.Meta                       `json:"meta"`
		Data MutasiPersediaanUpdateResponse `json:"data"`
	} `json:"body"`
}

// Delete
type MutasiPersediaanDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type MutasiPersediaanDeleteResponse struct {
	// model.MutasiPersediaanEntityModel
}
type MutasiPersediaanDeleteResponseDoc struct {
	Body struct {
		Meta res.Meta                       `json:"meta"`
		Data MutasiPersediaanDeleteResponse `json:"data"`
	} `json:"body"`
}

// export
type MutasiPersediaanExportRequest struct {
	// UserID int
	// Period             string `query:"period" validate:"required"`
	// Versions           int    `query:"versions" validate:"required"`
	// CompanyID          int    `query:"company_id"`
	MutasiPersediaanID int `query:"mutasi_persediaan_id" validate:"required"`
}

type MutasiPersediaanExportAsyncRequest struct {
	UserID    int
	Period    string `query:"period" validate:"required"`
	Versions  int    `query:"versions" validate:"required"`
	CompanyID int    `query:"company_id"`
}
type MutasiPersediaanExportResponse struct {
	FileName string `json:"filename"`
	Path     string `json:"path"`
}
type MutasiPersediaanExportResponseDoc struct {
	Body struct {
		Meta res.Meta                       `json:"meta"`
		Data MutasiPersediaanExportResponse `json:"data"`
	} `json:"body"`
}

// Import
type MutasiPersediaanImportRequest struct {
	File      multipart.File
	UserID    int
	CompanyID int
}
type MutasiPersediaanImportResponse struct {
	Data model.MutasiPersediaanEntityModel
}
type MutasiPersediaanImportResponseDoc struct {
	Body struct {
		Meta res.Meta                       `json:"meta"`
		Data MutasiPersediaanImportResponse `json:"data"`
	} `json:"body"`
}
