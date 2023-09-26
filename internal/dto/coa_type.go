package dto

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
)

// Get
type CoaTypeGetRequest struct {
	abstraction.Pagination
	model.CoaTypeFilterModel
}
type CoaTypeGetResponse struct {
	Datas          []model.CoaTypeEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type CoaTypeGetResponseDoc struct {
	Body struct {
		Meta res.Meta                   `json:"meta"`
		Data []model.CoaTypeEntityModel `json:"data"`
	} `json:"body"`
}

// GetByID
type CoaTypeGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type CoaTypeGetByIDResponse struct {
	model.CoaTypeEntityModel
}
type CoaTypeGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta               `json:"meta"`
		Data CoaTypeGetByIDResponse `json:"data"`
	} `json:"body"`
}

// Create
type CoaTypeCreateRequest struct {
	model.CoaTypeEntity
}
type CoaTypeCreateResponse struct {
	model.CoaTypeEntityModel
}
type CoaTypeCreateResponseDoc struct {
	Body struct {
		Meta res.Meta              `json:"meta"`
		Data CoaTypeCreateResponse `json:"data"`
	} `json:"body"`
}

// Update
type CoaTypeUpdateRequest struct {
	ID int `param:"id" validate:"required,numeric"`
	model.CoaTypeEntity
}
type CoaTypeUpdateResponse struct {
	model.CoaTypeEntityModel
}
type CoaTypeUpdateResponseDoc struct {
	Body struct {
		Meta res.Meta              `json:"meta"`
		Data CoaTypeUpdateResponse `json:"data"`
	} `json:"body"`
}

// Delete
type CoaTypeDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type CoaTypeDeleteResponse struct {
	// model.CoaTypeEntityModel
}
type CoaTypeDeleteResponseDoc struct {
	Body struct {
		Meta res.Meta              `json:"meta"`
		Data CoaTypeDeleteResponse `json:"data"`
	} `json:"body"`
}
