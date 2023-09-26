package dto

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
)

// Get
type AccessScopeGetRequest struct {
	abstraction.Pagination
	model.AccessScopeFilterModel
}
type AccessScopeGetResponse struct {
	Datas          []model.AccessScopeEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type AccessScopeGetResponseDoc struct {
	Body struct {
		Meta res.Meta                       `json:"meta"`
		Data []model.AccessScopeEntityModel `json:"data"`
	} `json:"body"`
}

// GetByID
type AccessScopeGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type AccessScopeGetByIDResponse struct {
	model.AccessScopeEntityModel
}
type AccessScopeGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta                   `json:"meta"`
		Data AccessScopeGetByIDResponse `json:"data"`
	} `json:"body"`
}

// Create
type AccessScopeCreateRequest struct {
	model.AccessScopeEntity
}
type AccessScopeCreateResponse struct {
	model.AccessScopeEntityModel
}
type AccessScopeCreateResponseDoc struct {
	Body struct {
		Meta res.Meta                  `json:"meta"`
		Data AccessScopeCreateResponse `json:"data"`
	} `json:"body"`
}

// Update
type AccessScopeUpdateRequest struct {
	ID        int   `param:"id" validate:"required,numeric"`
	AccessAll *bool `json:"access_all" validate:"required"`
	CompanyID []int `json:"company_id" validate:"required"`
}
type AccessScopeUpdateResponse struct {
	model.AccessScopeEntityModel
}
type AccessScopeUpdateResponseDoc struct {
	Body struct {
		Meta res.Meta                  `json:"meta"`
		Data AccessScopeUpdateResponse `json:"data"`
	} `json:"body"`
}

// Delete
type AccessScopeDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type AccessScopeDeleteResponse struct {
	model.AccessScopeEntityModel
}
type AccessScopeDeleteResponseDoc struct {
	Body struct {
		Meta res.Meta                  `json:"meta"`
		Data AccessScopeDeleteResponse `json:"data"`
	} `json:"body"`
}
