package dto

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
)

// Get
type TrialBalanceGetRequest struct {
	abstraction.Pagination
	model.TrialBalanceFilterModel
}
type TrialBalanceGetResponse struct {
	Datas          []model.TrialBalanceEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type TrialBalanceGetResponseDoc struct {
	Body struct {
		Meta res.Meta                        `json:"meta"`
		Data []model.TrialBalanceEntityModel `json:"data"`
	} `json:"body"`
}

// GetByID
type TrialBalanceGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type TrialBalanceGetByIDResponse struct {
	model.TrialBalanceEntityModel
}
type TrialBalanceGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta                    `json:"meta"`
		Data TrialBalanceGetByIDResponse `json:"data"`
	} `json:"body"`
}

// Create
type TrialBalanceCreateRequest struct {
	model.TrialBalanceEntity
}
type TrialBalanceCreateResponse struct {
	model.TrialBalanceEntityModel
}
type TrialBalanceCreateResponseDoc struct {
	Body struct {
		Meta res.Meta                   `json:"meta"`
		Data TrialBalanceCreateResponse `json:"data"`
	} `json:"body"`
}

// Update
type TrialBalanceUpdateRequest struct {
	ID int `param:"id" validate:"required,numeric"`
	model.TrialBalanceEntity
}
type TrialBalanceUpdateResponse struct {
	model.TrialBalanceEntityModel
}
type TrialBalanceUpdateResponseDoc struct {
	Body struct {
		Meta res.Meta                   `json:"meta"`
		Data TrialBalanceUpdateResponse `json:"data"`
	} `json:"body"`
}

// Delete
type TrialBalanceDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type TrialBalanceDeleteResponse struct {
	model.TrialBalanceEntityModel
}
type TrialBalanceDeleteResponseDoc struct {
	Body struct {
		Meta res.Meta                   `json:"meta"`
		Data TrialBalanceDeleteResponse `json:"data"`
	} `json:"body"`
}

// export
type TrialBalanceExportRequest struct {
	// UserID    int
	// Period    string `query:"period" validate:"required"`
	// Versions  int    `query:"versions" validate:"required"`
	// CompanyID int    `query:"company_id"`
	TrialBalanceID int `query:"trial_balance_id" validate:"required"`
}

type TrialBalanceExportAsyncRequest struct {
	UserID    int
	Period    string `query:"period" validate:"required"`
	Versions  int    `query:"versions" validate:"required"`
	CompanyID int    `query:"company_id"`
}
type TrialBalanceExportResponse struct {
	FileName string `json:"filename"`
	Path     string `json:"path"`
}
type TrialBalanceExportResponseDoc struct {
	Body struct {
		Meta res.Meta                   `json:"meta"`
		Data TrialBalanceExportResponse `json:"data"`
	} `json:"body"`
}

// Import
type TrialBalanceImportRequest struct {
	UserID    int
	CompanyID int
}
type TrialBalanceImportResponse struct {
	// Data model.TrialBalanceEntityModel
}
type TrialBalanceImportResponseDoc struct {
	Body struct {
		Meta res.Meta                   `json:"meta"`
		Data TrialBalanceImportResponse `json:"data"`
	} `json:"body"`
}
