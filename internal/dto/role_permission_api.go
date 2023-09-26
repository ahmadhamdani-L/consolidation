package dto

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
)

// Get
type RolePermissionApiGetRequest struct {
	abstraction.Pagination
	model.RolePermissionApiFilterModel
}
type RolePermissionApiGetResponse struct {
	Datas          []model.RolePermissionApiEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type RolePermissionApiGetResponseDoc struct {
	Body struct {
		Meta res.Meta                             `json:"meta"`
		Data []model.RolePermissionApiEntityModel `json:"data"`
	} `json:"body"`
}

// GetByID
type RolePermissionApiGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type RolePermissionApiGetByIDResponse struct {
	model.RolePermissionApiEntityModel
}
type RolePermissionApiGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta                         `json:"meta"`
		Data RolePermissionApiGetByIDResponse `json:"data"`
	} `json:"body"`
}

// Create
type RolePermissionApiCreateRequest struct {
	model.RolePermissionApiEntity
}
type RolePermissionApiCreateResponse struct {
	model.RolePermissionApiEntityModel
}
type RolePermissionApiCreateResponseDoc struct {
	Body struct {
		Meta res.Meta                        `json:"meta"`
		Data RolePermissionApiCreateResponse `json:"data"`
	} `json:"body"`
}

// Update
type RolePermissionApiUpdateRequest struct {
	ID int `param:"id" validate:"required,numeric"`
	model.RolePermissionApiEntity
}
type RolePermissionApiUpdateResponse struct {
	model.RolePermissionApiEntityModel
}
type RolePermissionApiUpdateResponseDoc struct {
	Body struct {
		Meta res.Meta                        `json:"meta"`
		Data RolePermissionApiUpdateResponse `json:"data"`
	} `json:"body"`
}

// Delete
type RolePermissionApiDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type RolePermissionApiDeleteResponse struct {
	model.RolePermissionApiEntityModel
}
type RolePermissionApiDeleteResponseDoc struct {
	Body struct {
		Meta res.Meta                        `json:"meta"`
		Data RolePermissionApiDeleteResponse `json:"data"`
	} `json:"body"`
}
