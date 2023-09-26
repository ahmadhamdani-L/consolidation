package dto

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
)

// Get
type RolePermissionGetRequest struct {
	abstraction.Pagination
	model.RolePermissionFilterModel
}
type RolePermissionGetResponse struct {
	Datas          []model.RolePermissionEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type RolePermissionGetResponseDoc struct {
	Body struct {
		Meta res.Meta                          `json:"meta"`
		Data []model.RolePermissionEntityModel `json:"data"`
	} `json:"body"`
}

// GetByID
type RolePermissionGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type RolePermissionGetByIDResponse struct {
	model.RolePermissionEntityModel
}
type RolePermissionGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta                      `json:"meta"`
		Data RolePermissionGetByIDResponse `json:"data"`
	} `json:"body"`
}

// Create
type RolePermissionCreateRequest struct {
	model.RolePermissionEntity
}
type RolePermissionCreateResponse struct {
	model.RolePermissionEntityModel
}
type RolePermissionCreateResponseDoc struct {
	Body struct {
		Meta res.Meta                     `json:"meta"`
		Data RolePermissionCreateResponse `json:"data"`
	} `json:"body"`
}

// Update
type RolePermissionUpdateRequest struct {
	ID int `param:"id" validate:"required,numeric"`
	model.RolePermissionEntity
}
type RolePermissionUpdateResponse struct {
	model.RolePermissionEntityModel
}
type RolePermissionUpdateResponseDoc struct {
	Body struct {
		Meta res.Meta                     `json:"meta"`
		Data RolePermissionUpdateResponse `json:"data"`
	} `json:"body"`
}

// Delete
type RolePermissionDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type RolePermissionDeleteResponse struct {
	model.RolePermissionEntityModel
}
type RolePermissionDeleteResponseDoc struct {
	Body struct {
		Meta res.Meta                     `json:"meta"`
		Data RolePermissionDeleteResponse `json:"data"`
	} `json:"body"`
}
