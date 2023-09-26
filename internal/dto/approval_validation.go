package dto

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
)

type ApproveValidationRequest struct {
	TrialBalanceID int `json:"validation_id" validate:"required"`
}

type ApproveValidationResponse struct {
	Success bool
}

type ApproveValidationResponseDoc struct {
	Body struct {
		Meta res.Meta                  `json:"meta"`
		Data ApproveValidationResponse `json:"data"`
	} `json:"body"`
}

// Get
type ApprovalValidationGetRequest struct {
	abstraction.Pagination
	model.TrialBalanceFilterModel
}
type ApprovalValidationGetResponse struct {
	Datas          []model.TrialBalanceEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type ApprovalValidationGetResponseDoc struct {
	Body struct {
		Meta res.Meta                        `json:"meta"`
		Data []model.TrialBalanceEntityModel `json:"data"`
	} `json:"body"`
}
