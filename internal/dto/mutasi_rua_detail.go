package dto

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
)

// Get
type MutasiRuaDetailGetRequest struct {
	abstraction.Pagination
	model.MutasiRuaDetailFilterModel
}
type MutasiRuaDetailGetResponse struct {
	Datas          model.MutasiRuaEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type MutasiRuaDetailGetResponseDoc struct {
	Body struct {
		Meta res.Meta                   `json:"meta"`
		Data MutasiRuaDetailGetResponse `json:"data"`
	} `json:"body"`
}

// GetByID
type MutasiRuaDetailGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type MutasiRuaDetailGetByIDResponse struct {
	model.MutasiRuaDetailEntityModel
}
type MutasiRuaDetailGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta                       `json:"meta"`
		Data MutasiRuaDetailGetByIDResponse `json:"data"`
	} `json:"body"`
}

// Create
type MutasiRuaDetailCreateRequest struct {
	model.MutasiRuaDetailEntity
}
type MutasiRuaDetailCreateResponse struct {
	model.MutasiRuaDetailEntityModel
}
type MutasiRuaDetailCreateResponseDoc struct {
	Body struct {
		Meta res.Meta                      `json:"meta"`
		Data MutasiRuaDetailCreateResponse `json:"data"`
	} `json:"body"`
}

// Update
type MutasiRuaDetailUpdateRequest struct {
	ID                      int      `param:"id" validate:"required,numeric"`
	BeginningBalance        *float64 `json:"beginning_balance" validate:"required" example:"10000.00"`
	AcquisitionOfSubsidiary *float64 `json:"acquisition_of_subsidiary" validate:"required" example:"10000.00"`
	Additions               *float64 `json:"additions" validate:"required" example:"10000.00"`
	Deductions              *float64 `json:"deductions" validate:"required" example:"10000.00"`
	Reclassification        *float64 `json:"reclassification" validate:"required" example:"10000.00"`
	Remeasurement           *float64 `json:"remeasurement" validate:"required" example:"10000.00"`
}
type MutasiRuaDetailUpdateResponse struct {
	model.MutasiRuaDetailEntityModel
}
type MutasiRuaDetailUpdateResponseDoc struct {
	Body struct {
		Meta res.Meta                      `json:"meta"`
		Data MutasiRuaDetailUpdateResponse `json:"data"`
	} `json:"body"`
}

// Delete
type MutasiRuaDetailDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type MutasiRuaDetailDeleteResponse struct {
	// model.MutasiRuaDetailEntityModel
}
type MutasiRuaDetailDeleteResponseDoc struct {
	Body struct {
		Meta res.Meta                      `json:"meta"`
		Data MutasiRuaDetailDeleteResponse `json:"data"`
	} `json:"body"`
}

// export
type MutasiRuaDetailExportRequest struct {
	model.MutasiRuaDetailFilterModel
}
type MutasiRuaDetailExportResponse struct {
	File string `json:"file"`
}
type MutasiRuaDetailExportResponseDoc struct {
	Body struct {
		Meta res.Meta                      `json:"meta"`
		Data MutasiRuaDetailExportResponse `json:"data"`
	} `json:"body"`
}

// Import
type MutasiRuaDetailImportRequest struct {
	Datas []model.MutasiRuaDetailEntity
}
type MutasiRuaDetailImportResponse struct {
	Datas          []model.MutasiRuaDetailEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type MutasiRuaDetailImportResponseDoc struct {
	Body struct {
		Meta res.Meta                           `json:"meta"`
		Data []model.MutasiRuaDetailEntityModel `json:"data"`
	} `json:"body"`
}
