package dto

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
)

// Get
type ParameterGetRequest struct {
	abstraction.Pagination
	model.ParameterFilterModel
}
type ParameterGetResponse struct {
	Datas          []model.ParameterEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type ParameterGetResponseDoc struct {
	Body struct {
		Meta res.Meta                     `json:"meta"`
		Data []model.ParameterEntityModel `json:"data"`
	} `json:"body"`
}

// GetByID
type ParameterGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type ParameterGetByIDResponse struct {
	model.ParameterEntityModel
}
type ParameterGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta                 `json:"meta"`
		Data ParameterGetByIDResponse `json:"data"`
	} `json:"body"`
}

// Create
type ParameterCreateRequest struct {
	model.ParameterEntity
}
type ParameterCreateResponse struct {
	model.ParameterEntityModel
}
type ParameterCreateResponseDoc struct {
	Body struct {
		Meta res.Meta                `json:"meta"`
		Data ParameterCreateResponse `json:"data"`
	} `json:"body"`
}

// Update
type ParameterUpdateRequest struct {
	ID int `param:"id" validate:"required,numeric"`
	model.ParameterEntity
}
type ParameterUpdateResponse struct {
	model.ParameterEntityModel
}
type ParameterUpdateResponseDoc struct {
	Body struct {
		Meta res.Meta                `json:"meta"`
		Data ParameterUpdateResponse `json:"data"`
	} `json:"body"`
}

// Delete
type ParameterDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type ParameterDeleteResponse struct {
	// model.ParameterEntityModel
}
type ParameterDeleteResponseDoc struct {
	Body struct {
		Meta res.Meta                `json:"meta"`
		Data ParameterDeleteResponse `json:"data"`
	} `json:"body"`
}
