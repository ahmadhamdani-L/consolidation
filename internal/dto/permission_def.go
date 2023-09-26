package dto

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
)

// Get
type PermissionDefGetRequest struct {
	abstraction.Pagination
	model.PermissionDefFilterModel
}
type PermissionDefGetResponse struct {
	Datas          []model.PermissionDefEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type PermissionDefGetResponseDoc struct {
	Body struct {
		Meta res.Meta                         `json:"meta"`
		Data []model.PermissionDefEntityModel `json:"data"`
	} `json:"body"`
}

// GetByID
type PermissionDefGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type PermissionDefGetByIDResponse struct {
	model.PermissionDefEntityModel
}
type PermissionDefGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta                     `json:"meta"`
		Data PermissionDefGetByIDResponse `json:"data"`
	} `json:"body"`
}

// Create
type PermissionDefCreateRequest struct {
	model.PermissionDefEntity
}
type PermissionDefCreateResponse struct {
	model.PermissionDefEntityModel
}
type PermissionDefCreateResponseDoc struct {
	Body struct {
		Meta res.Meta                    `json:"meta"`
		Data PermissionDefCreateResponse `json:"data"`
	} `json:"body"`
}

// Update
type PermissionDefUpdateRequest struct {
	ID int `param:"id" validate:"required,numeric"`
	model.PermissionDefEntity
}
type PermissionDefUpdateResponse struct {
	model.PermissionDefEntityModel
}
type PermissionDefUpdateResponseDoc struct {
	Body struct {
		Meta res.Meta                    `json:"meta"`
		Data PermissionDefUpdateResponse `json:"data"`
	} `json:"body"`
}

// Delete
type PermissionDefDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type PermissionDefDeleteResponse struct {
	model.PermissionDefEntityModel
}
type PermissionDefDeleteResponseDoc struct {
	Body struct {
		Meta res.Meta                    `json:"meta"`
		Data PermissionDefDeleteResponse `json:"data"`
	} `json:"body"`
}
