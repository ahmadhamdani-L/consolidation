package dto

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
)

// Get
type MutasiIaDetailGetRequest struct {
	abstraction.Pagination
	model.MutasiIaDetailFilterModel
}
type MutasiIaDetailGetResponse struct {
	Datas model.MutasiIaEntityModel
}
type MutasiIaDetailGetResponseDoc struct {
	Body struct {
		Meta res.Meta                  `json:"meta"`
		Data MutasiIaDetailGetResponse `json:"data"`
	} `json:"body"`
}

// GetByID
type MutasiIaDetailGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type MutasiIaDetailGetByIDResponse struct {
	model.MutasiIaDetailEntityModel
}
type MutasiIaDetailGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta                      `json:"meta"`
		Data MutasiIaDetailGetByIDResponse `json:"data"`
	} `json:"body"`
}

// Create
type MutasiIaDetailCreateRequest struct {
	model.MutasiIaDetailEntity
}
type MutasiIaDetailCreateResponse struct {
	model.MutasiIaDetailEntityModel
}
type MutasiIaDetailCreateResponseDoc struct {
	Body struct {
		Meta res.Meta                     `json:"meta"`
		Data MutasiIaDetailCreateResponse `json:"data"`
	} `json:"body"`
}

// Update
type MutasiIaDetailUpdateRequest struct {
	ID                      int      `param:"id" validate:"required,numeric"`
	BeginningBalance        *float64 `json:"beginning_balance" validate:"required" example:"10000.00"`
	AcquisitionOfSubsidiary *float64 `json:"acquisition_of_subsidiary" validate:"required" example:"10000.00"`
	Additions               *float64 `json:"additions" validate:"required" example:"10000.00"`
	Deductions              *float64 `json:"deductions" validate:"required" example:"10000.00"`
	Reclassification        *float64 `json:"reclassification" validate:"required" example:"10000.00"`
	Revaluation             *float64 `json:"revaluation" validate:"required" example:"10000.00"`
	// EndingBalance           *float64 `json:"ending_balance" validate:"required" example:"10000.00"`
	// Control                 *float64 `json:"control" validate:"required" example:"10000.00"`
}
type MutasiIaDetailUpdateResponse struct {
	model.MutasiIaDetailEntityModel
}
type MutasiIaDetailUpdateResponseDoc struct {
	Body struct {
		Meta res.Meta                     `json:"meta"`
		Data MutasiIaDetailUpdateResponse `json:"data"`
	} `json:"body"`
}

// Delete
type MutasiIaDetailDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type MutasiIaDetailDeleteResponse struct {
	model.MutasiIaDetailEntityModel
}
type MutasiIaDetailDeleteResponseDoc struct {
	Body struct {
		Meta res.Meta                     `json:"meta"`
		Data MutasiIaDetailDeleteResponse `json:"data"`
	} `json:"body"`
}

// export
type MutasiIaDetailExportRequest struct {
	model.MutasiIaDetailFilterModel
}
type MutasiIaDetailExportResponse struct {
	File string `json:"file"`
}
type MutasiIaDetailExportResponseDoc struct {
	Body struct {
		Meta res.Meta                     `json:"meta"`
		Data MutasiIaDetailExportResponse `json:"data"`
	} `json:"body"`
}

// Import
type MutasiIaDetailImportRequest struct {
	Datas []model.MutasiIaDetailEntity
}
type MutasiIaDetailImportResponse struct {
	Datas          []model.MutasiIaDetailEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type MutasiIaDetailImportResponseDoc struct {
	Body struct {
		Meta res.Meta                          `json:"meta"`
		Data []model.MutasiIaDetailEntityModel `json:"data"`
	} `json:"body"`
}
