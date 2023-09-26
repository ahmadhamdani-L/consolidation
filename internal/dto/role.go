package dto

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
)

// Get
type RoleGetRequest struct {
	abstraction.Pagination
	model.RoleFilterModel
}
type RoleGetResponse struct {
	Datas          []model.RoleEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type RoleGetResponseDoc struct {
	Body struct {
		Meta res.Meta                `json:"meta"`
		Data []model.RoleEntityModel `json:"data"`
	} `json:"body"`
}

// GetByID
type RoleGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type RoleGetByIDResponse struct {
	model.RoleEntityModel
}
type RoleGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta            `json:"meta"`
		Data RoleGetByIDResponse `json:"data"`
	} `json:"body"`
}

// Create
type RoleCreateRequest struct {
	model.RoleEntity
	RolePermission []struct {
		FunctionalID string `json:"functional_id" validate:"required"`
		Create       *bool  `json:"create" validate:"required"`
		Read         *bool  `json:"read" validate:"required"`
		Update       *bool  `json:"update" validate:"required"`
		Delete       *bool  `json:"delete" validate:"required"`
	} `json:"role_permissions" validate:"required"`
}
type RoleCreateResponse struct {
	model.RoleEntityModel
}
type RoleCreateResponseDoc struct {
	Body struct {
		Meta res.Meta           `json:"meta"`
		Data RoleCreateResponse `json:"data"`
	} `json:"body"`
}

// Update
type RoleUpdateRequest struct {
	ID int `param:"id" validate:"required,numeric"`
	model.RoleEntity
	RolePermission []struct {
		FunctionalID string `json:"functional_id" validate:"required"`
		Create       *bool  `json:"create" validate:"required"`
		Read         *bool  `json:"read" validate:"required"`
		Update       *bool  `json:"update" validate:"required"`
		Delete       *bool  `json:"delete" validate:"required"`
	} `json:"role_permissions" validate:"required"`
}
type RoleUpdateResponse struct {
	model.RoleEntityModel
}
type RoleUpdateResponseDoc struct {
	Body struct {
		Meta res.Meta           `json:"meta"`
		Data RoleUpdateResponse `json:"data"`
	} `json:"body"`
}

// Delete
type RoleDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type RoleDeleteResponse struct {
	model.RoleEntityModel
}
type RoleDeleteResponseDoc struct {
	Body struct {
		Meta res.Meta           `json:"meta"`
		Data RoleDeleteResponse `json:"data"`
	} `json:"body"`
}

type RoleDeletePermissionRequest struct {
	ID           int    `param:"id" validate:"required,numeric"`
	FunctionalID string `query:"functional_id" validate:"required,numeric"`
}
type RoleDeletePermissionResponse struct {
	model.RolePermissionEntityModel
}
type RoleDeletePermissionResponseDoc struct {
	Body struct {
		Meta res.Meta           `json:"meta"`
		Data RoleDeleteResponse `json:"data"`
	} `json:"body"`
}
