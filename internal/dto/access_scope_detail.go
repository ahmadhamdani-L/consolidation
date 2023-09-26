package dto

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
)

// Get
type AccessScopeDetailGetRequest struct {
	abstraction.Pagination
	model.AccessScopeDetailFilterModel
}
type AccessScopeDetailGetResponse struct {
	Datas          []model.AccessScopeDetailEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type AccessScopeDetailGetResponseDoc struct {
	Body struct {
		Meta res.Meta                             `json:"meta"`
		Data []model.AccessScopeDetailEntityModel `json:"data"`
	} `json:"body"`
}

// GetByID
type AccessScopeDetailGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type AccessScopeDetailGetByIDResponse struct {
	model.AccessScopeDetailEntityModel
}
type AccessScopeDetailGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta                         `json:"meta"`
		Data AccessScopeDetailGetByIDResponse `json:"data"`
	} `json:"body"`
}

// Create
type AccessScopeDetailCreateRequest struct {
	model.AccessScopeDetailEntity
}
type AccessScopeDetailCreateResponse struct {
	model.AccessScopeDetailEntityModel
}
type AccessScopeDetailCreateResponseDoc struct {
	Body struct {
		Meta res.Meta                        `json:"meta"`
		Data AccessScopeDetailCreateResponse `json:"data"`
	} `json:"body"`
}

// Update
type AccessScopeDetailUpdateRequest struct {
	ID int `param:"id" validate:"required,numeric"`
	model.AccessScopeDetailEntity
}
type AccessScopeDetailUpdateResponse struct {
	model.AccessScopeDetailEntityModel
}
type AccessScopeDetailUpdateResponseDoc struct {
	Body struct {
		Meta res.Meta                        `json:"meta"`
		Data AccessScopeDetailUpdateResponse `json:"data"`
	} `json:"body"`
}

// Delete
type AccessScopeDetailDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type AccessScopeDetailDeleteResponse struct {
	model.AccessScopeDetailEntityModel
}
type AccessScopeDetailDeleteResponseDoc struct {
	Body struct {
		Meta res.Meta                        `json:"meta"`
		Data AccessScopeDetailDeleteResponse `json:"data"`
	} `json:"body"`
}

// GetByID
type AccessScopeDetailGetCompanyListRequest struct {
	ID int `query:"access_scope_id" validate:"required,numeric"`
}
type AccessScopeDetailGetCompanyListResponse struct {
	Data []model.AccessScopeDetailListEntityModel
}
type AccessScopeDetailGetCompanyListResponseDoc struct {
	Body struct {
		Meta res.Meta                                 `json:"meta"`
		Data []model.AccessScopeDetailListEntityModel `json:"data"`
	} `json:"body"`
}
