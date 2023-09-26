package dto

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
)

// Get
type EmployeeBenefitDetailGetRequest struct {
	abstraction.Pagination
	model.EmployeeBenefitDetailFilterModel
}
type EmployeeBenefitDetailGetResponse struct {
	Datas model.EmployeeBenefitEntityModel
}
type EmployeeBenefitDetailGetResponseDoc struct {
	Body struct {
		Meta res.Meta                         `json:"meta"`
		Data EmployeeBenefitDetailGetResponse `json:"data"`
	} `json:"body"`
}

// GetByID
type EmployeeBenefitDetailGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type EmployeeBenefitDetailGetByIDResponse struct {
	model.EmployeeBenefitDetailEntityModel
}
type EmployeeBenefitDetailGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta                             `json:"meta"`
		Data EmployeeBenefitDetailGetByIDResponse `json:"data"`
	} `json:"body"`
}

// Create
type EmployeeBenefitDetailCreateRequest struct {
	model.EmployeeBenefitDetailEntity
}
type EmployeeBenefitDetailCreateResponse struct {
	model.EmployeeBenefitDetailEntityModel
}
type EmployeeBenefitDetailCreateResponseDoc struct {
	Body struct {
		Meta res.Meta                            `json:"meta"`
		Data EmployeeBenefitDetailCreateResponse `json:"data"`
	} `json:"body"`
}

// Update
type EmployeeBenefitDetailUpdateRequest struct {
	ID     int      `param:"id" validate:"required,numeric"`
	Amount *float64 `json:"amount" example:"10000.00"`
	Value  *string  `json:"value" example:"10000.00"`
}
type EmployeeBenefitDetailUpdateResponse struct {
	model.EmployeeBenefitDetailEntityModel
}
type EmployeeBenefitDetailUpdateResponseDoc struct {
	Body struct {
		Meta res.Meta                            `json:"meta"`
		Data EmployeeBenefitDetailUpdateResponse `json:"data"`
	} `json:"body"`
}

// Delete
type EmployeeBenefitDetailDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type EmployeeBenefitDetailDeleteResponse struct {
	// model.EmployeeBenefitDetailEntityModel
}
type EmployeeBenefitDetailDeleteResponseDoc struct {
	Body struct {
		Meta res.Meta                            `json:"meta"`
		Data EmployeeBenefitDetailDeleteResponse `json:"data"`
	} `json:"body"`
}
