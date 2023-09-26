package dto

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
)

// Get
type MutasiFaDetailGetRequest struct {
	abstraction.Pagination
	model.MutasiFaDetailFilterModel
}
type MutasiFaDetailGetResponse struct {
	Datas model.MutasiFaEntityModel
}
type MutasiFaDetailGetResponseDoc struct {
	Body struct {
		Meta res.Meta                  `json:"meta"`
		Data MutasiFaDetailGetResponse `json:"data"`
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
