package dto

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
)

// Get
type ConsolidationDetailGetRequest struct {
	abstraction.Pagination
	model.ConsolidationDetailFilterModel
}
type ConsolidationDetailGetResponse struct {
	Datas          model.ConsolidationEntityModel
	ChildCompany   []model.ConsolidationBridgeEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type ConsolidationDetailGetResponseDoc struct {
	Body struct {
		Meta res.Meta                       `json:"meta"`
		Data ConsolidationDetailGetResponse `json:"data"`
	} `json:"body"`
}

// GetByID
type ConsolidationDetailGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type ConsolidationDetailGetByIDResponse struct {
	model.ConsolidationDetailEntityModel
}
type ConsolidationDetailGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta                           `json:"meta"`
		Data ConsolidationDetailGetByIDResponse `json:"data"`
	} `json:"body"`
}

// Create
type ConsolidationDetailCreateRequest struct {
	model.ConsolidationDetailEntity
}
type ConsolidationDetailCreateResponse struct {
	model.ConsolidationDetailEntityModel
}
type ConsolidationDetailCreateResponseDoc struct {
	Body struct {
		Meta res.Meta                          `json:"meta"`
		Data ConsolidationDetailCreateResponse `json:"data"`
	} `json:"body"`
}

// Update
type ConsolidationDetailUpdateRequest struct {
	ID              int      `param:"id" validate:"required,numeric"`
	AmountBeforeAje *float64 `json:"amount_before_aje" validate:"required" min:"0" example:"10000.00"`
}
type ConsolidationDetailUpdateResponse struct {
	model.ConsolidationDetailEntityModel
}
type ConsolidationDetailUpdateResponseDoc struct {
	Body struct {
		Meta res.Meta                          `json:"meta"`
		Data ConsolidationDetailUpdateResponse `json:"data"`
	} `json:"body"`
}

// Delete
type ConsolidationDetailDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type ConsolidationDetailDeleteResponse struct {
	// model.ConsolidationDetailEntityModel
}
type ConsolidationDetailDeleteResponseDoc struct {
	Body struct {
		Meta res.Meta                          `json:"meta"`
		Data ConsolidationDetailDeleteResponse `json:"data"`
	} `json:"body"`
}

// export
type ConsolidationDetailExportRequest struct {
	model.ConsolidationDetailFilterModel
}
type ConsolidationDetailExportResponse struct {
	File string `json:"file"`
}
type ConsolidationDetailExportResponseDoc struct {
	Body struct {
		Meta res.Meta                          `json:"meta"`
		Data ConsolidationDetailExportResponse `json:"data"`
	} `json:"body"`
}

// Import
type ConsolidationDetailImportRequest struct {
	Datas []model.ConsolidationDetailEntity
}
type ConsolidationDetailImportResponse struct {
	Datas          []model.ConsolidationDetailEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type ConsolidationDetailImportResponseDoc struct {
	Body struct {
		Meta res.Meta                               `json:"meta"`
		Data []model.ConsolidationDetailEntityModel `json:"data"`
	} `json:"body"`
}

type ConsolidationDetailGetByParentRequest struct {
	ParentID       int    `query:"parent_id"`
	ParentCode     string `query:"parent"`
	TrialBalanceID int    `query:"consolidation_id" validate:"required,numeric"`
	model.ConsolidationDetailFilterModel
}

type ConsolidationDetailGetByParentRequests struct {
	ParentID       int    `query:"parent_id"`
	ParentCode     string `query:"parent"`
	ConsolidationID int    `query:"consolidation_id" validate:"required,numeric"`
}

type ConsolidationDetailGetByParentResponse struct {
	Data *[]model.ConsolidationDetailFmtEntityModel
	// Data map[int]*model.TrialBalanceDetailFmtEntityModel
}

type ConsolidationDetailGetByParentResponseDoc struct {
	Body struct {
		Meta res.Meta                              `json:"meta"`
		Data ConsolidationDetailGetByParentResponse `json:"data"`
	} `json:"body"`
}

type ConsolidationDetailGetAllResponse struct {
	// Data model.TrialBalanceEntityModel
	Data map[int]model.ConsolidationDetailFmtEntityModel
}

type ConsolidationDetailGetAllResponseDoc struct {
	Body struct {
		Meta res.Meta                         `json:"meta"`
		Data ConsolidationDetailGetAllResponse `json:"data"`
	} `json:"body"`
}

type ConsolidationDetailGetByParentResponses struct {
	Data model.ConsolidationEntityModel
	// Data map[int]*model.TrialBalanceDetailFmtEntityModel
}