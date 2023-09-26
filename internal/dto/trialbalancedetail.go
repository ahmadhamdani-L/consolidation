package dto

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
)

// Get
type TrialBalanceDetailGetRequest struct {
	abstraction.Pagination
	model.TrialBalanceDetailFilterModel
}
type TrialBalanceDetailGetResponse struct {
	Datas          model.TrialBalanceEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type TrialBalanceDetailGetResponseDoc struct {
	Body struct {
		Meta res.Meta                      `json:"meta"`
		Data TrialBalanceDetailGetResponse `json:"data"`
	} `json:"body"`
}

// GetByID
type TrialBalanceDetailGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type TrialBalanceDetailGetByIDResponse struct {
	model.TrialBalanceDetailEntityModel
}
type TrialBalanceDetailGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta                          `json:"meta"`
		Data TrialBalanceDetailGetByIDResponse `json:"data"`
	} `json:"body"`
}

// Create
type TrialBalanceDetailCreateRequest struct {
	model.TrialBalanceDetailEntity
}
type TrialBalanceDetailCreateResponse struct {
	model.TrialBalanceDetailEntityModel
}
type TrialBalanceDetailCreateResponseDoc struct {
	Body struct {
		Meta res.Meta                         `json:"meta"`
		Data TrialBalanceDetailCreateResponse `json:"data"`
	} `json:"body"`
}

// Update
type TrialBalanceDetailUpdateRequest struct {
	ID              int      `param:"id" validate:"required,numeric"`
	AmountBeforeAje *float64 `json:"amount_before_aje" validate:"required" min:"0" example:"10000.00"`
}
type TrialBalanceDetailUpdateResponse struct {
	model.TrialBalanceDetailEntityModel
}
type TrialBalanceDetailUpdateResponseDoc struct {
	Body struct {
		Meta res.Meta                         `json:"meta"`
		Data TrialBalanceDetailUpdateResponse `json:"data"`
	} `json:"body"`
}

// Delete
type TrialBalanceDetailDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type TrialBalanceDetailDeleteResponse struct {
	// model.TrialBalanceDetailEntityModel
}
type TrialBalanceDetailDeleteResponseDoc struct {
	Body struct {
		Meta res.Meta                         `json:"meta"`
		Data TrialBalanceDetailDeleteResponse `json:"data"`
	} `json:"body"`
}

// export
type TrialBalanceDetailExportRequest struct {
	model.TrialBalanceDetailFilterModel
}
type TrialBalanceDetailExportResponse struct {
	File string `json:"file"`
}
type TrialBalanceDetailExportResponseDoc struct {
	Body struct {
		Meta res.Meta                         `json:"meta"`
		Data TrialBalanceDetailExportResponse `json:"data"`
	} `json:"body"`
}

// Import
type TrialBalanceDetailImportRequest struct {
	Datas []model.TrialBalanceDetailEntity
}
type TrialBalanceDetailImportResponse struct {
	Datas          []model.TrialBalanceDetailEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type TrialBalanceDetailImportResponseDoc struct {
	Body struct {
		Meta res.Meta                              `json:"meta"`
		Data []model.TrialBalanceDetailEntityModel `json:"data"`
	} `json:"body"`
}

type TrialBalanceDetailGetByParentRequest struct {
	ParentID       int    `query:"parent_id"`
	ParentCode     string `query:"parent"`
	TrialBalanceID int    `query:"trial_balance_id" validate:"required,numeric"`
}

type TrialBalanceDetailGetByParentResponse struct {
	AmountBeforeAje *float64                      `json:"amount_before_aje"`
	AmountAjeCr     *float64                      `json:"amount_aje_cr"`
	AmountAjeDr     *float64                      `json:"amount_aje_dr"`
	AmountAfterAje  *float64                      `json:"amount_after_aje"`
	Data            model.TrialBalanceEntityModel `json:"trial_balance"`
	// Data map[int]*model.TrialBalanceDetailFmtEntityModel
}

type TrialBalanceDetailGetByParentResponseDoc struct {
	Body struct {
		Meta res.Meta                              `json:"meta"`
		Data TrialBalanceDetailGetByParentResponse `json:"data"`
	} `json:"body"`
}

type TrialBalanceDetailGetAllResponse struct {
	// Data model.TrialBalanceEntityModel
	Data map[int]model.TrialBalanceDetailFmtEntityModel
}

type TrialBalanceDetailGetAllResponseDoc struct {
	Body struct {
		Meta res.Meta                         `json:"meta"`
		Data TrialBalanceDetailGetAllResponse `json:"data"`
	} `json:"body"`
}
