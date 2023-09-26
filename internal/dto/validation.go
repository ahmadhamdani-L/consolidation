package dto

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
)

// Get
type ValidationGetRequest struct {
	abstraction.Pagination
	model.ValidationFilterModel
}
type ValidationGetResponse struct {
	Datas          []model.TrialBalanceEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type ValidationGetResponseDoc struct {
	Body struct {
		Meta res.Meta                        `json:"meta"`
		Data []model.TrialBalanceEntityModel `json:"data"`
	} `json:"body"`
}

// GetByID
type ValidationGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type ValidationGetByIDResponse struct {
	Data model.ValidationEntityModel
}
type ValidationGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta                  `json:"meta"`
		Data ValidationGetByIDResponse `json:"data"`
	} `json:"body"`
}

// validate
type ValidationValidateRequest struct {
	CompanyID      int    `json:"company_id" validate:"required"`
	Period         string `json:"period" validate:"required"`
	ListValidation []int  `json:"list_to_validate" validate:"required"`
}
type ValidationValidateResponse struct {
	model.ValidationEntityModel
}
type ValidationValidateResponseDoc struct {
	Body struct {
		Meta res.Meta                   `json:"meta"`
		Data ValidationValidateResponse `json:"data"`
	} `json:"body"`
}

// Get
type ValidationGetListAvailable struct {
	CompanyID int    `query:"company_id" validate:"required"`
	Period    string `query:"period" validate:"required"`
	abstraction.Pagination
}
type ValidationGetListAvailableResponse struct {
	Datas          []model.TrialBalanceEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type ValidationGetListAvailableResponseDoc struct {
	Body struct {
		Meta res.Meta                        `json:"meta"`
		Data []model.TrialBalanceEntityModel `json:"data"`
	} `json:"body"`
}

// validate
type ValidationValidateModulRequest struct {
	ValidationMasterID int   `param:"validation_id" validate:"required"`   //trial_balance_id, parent dari yang mau di validate
	ListValidation     []int `json:"list_to_validate" validate:"required"` //harus ada minimal si masternya itu dimasukin ke list ini
}
type ValidationValidateModulResponse struct {
	model.ValidationEntityModel
}
type ValidationValidateModulResponseDoc struct {
	Body struct {
		Meta res.Meta                   `json:"meta"`
		Data ValidationValidateResponse `json:"data"`
	} `json:"body"`
}
