package dto

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
	"mime/multipart"
)

// Get
type MutasiFaGetRequest struct {
	abstraction.Pagination
	model.MutasiFaFilterModel
}
type MutasiFaGetResponse struct {
	Datas          []model.MutasiFaEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type MutasiFaGetResponseDoc struct {
	Body struct {
		Meta res.Meta                    `json:"meta"`
		Data []model.MutasiFaEntityModel `json:"data"`
	} `json:"body"`
}

// GetByID
type MutasiFaGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type MutasiFaGetByIDResponse struct {
	model.MutasiFaEntityModel
}
type MutasiFaGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta                `json:"meta"`
		Data MutasiFaGetByIDResponse `json:"data"`
	} `json:"body"`
}

// Create
type MutasiFaCreateRequest struct {
	model.MutasiFaEntity
}
type MutasiFaCreateResponse struct {
	model.MutasiFaEntityModel
}
type MutasiFaCreateResponseDoc struct {
	Body struct {
		Meta res.Meta               `json:"meta"`
		Data MutasiFaCreateResponse `json:"data"`
	} `json:"body"`
}

// Update
type MutasiFaUpdateRequest struct {
	ID int `param:"id" validate:"required,numeric"`
	model.MutasiFaEntity
}
type MutasiFaUpdateResponse struct {
	model.MutasiFaEntityModel
}
type MutasiFaUpdateResponseDoc struct {
	Body struct {
		Meta res.Meta               `json:"meta"`
		Data MutasiFaUpdateResponse `json:"data"`
	} `json:"body"`
}

// Delete
type MutasiFaDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type MutasiFaDeleteResponse struct {
	// model.MutasiFaEntityModel
}
type MutasiFaDeleteResponseDoc struct {
	Body struct {
		Meta res.Meta               `json:"meta"`
		Data MutasiFaDeleteResponse `json:"data"`
	} `json:"body"`
}

// export
type MutasiFaExportRequest struct {
	// UserID    int
	// Period    string `query:"period" validate:"required"`
	// Versions  int    `query:"versions" validate:"required"`
	// CompanyID int    `query:"company_id"`
	MutasiFaID int `query:"mutasi_fa_id" validate:"required"`
}
type MutasiFaExportResponse struct {
	FileName string `json:"filename"`
	Path     string `json:"path"`
}
type MutasiFaExportResponseDoc struct {
	Body struct {
		Meta res.Meta               `json:"meta"`
		Data MutasiFaExportResponse `json:"data"`
	} `json:"body"`
}

// Import
type MutasiFaImportRequest struct {
	UserID    int
	CompanyID int
	File      multipart.File
}
type MutasiFaImportResponse struct {
	Data model.MutasiFaEntityModel
}
type MutasiFaImportResponseDoc struct {
	Body struct {
		Meta res.Meta               `json:"meta"`
		Data MutasiFaImportResponse `json:"data"`
	} `json:"body"`
}
