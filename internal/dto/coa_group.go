package dto

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
)

// Get
type CoaGroupGetRequest struct {
	abstraction.Pagination
	model.CoaGroupFilterModel
}
type CoaGroupGetResponse struct {
	Datas          []model.CoaGroupEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type CoaGroupGetResponseDoc struct {
	Body struct {
		Meta res.Meta                    `json:"meta"`
		Data []model.CoaGroupEntityModel `json:"data"`
	} `json:"body"`
}

// GetByID
type CoaGroupGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type CoaGroupGetByIDResponse struct {
	model.CoaGroupEntityModel
}
type CoaGroupGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta                `json:"meta"`
		Data CoaGroupGetByIDResponse `json:"data"`
	} `json:"body"`
}

// Create
type CoaGroupCreateRequest struct {
	model.CoaGroupEntity
}
type CoaGroupCreateResponse struct {
	model.CoaGroupEntityModel
}
type CoaGroupCreateResponseDoc struct {
	Body struct {
		Meta res.Meta               `json:"meta"`
		Data CoaGroupCreateResponse `json:"data"`
	} `json:"body"`
}

// Update
type CoaGroupUpdateRequest struct {
	ID int `param:"id" validate:"required,numeric"`
	model.CoaGroupEntity
}
type CoaGroupUpdateResponse struct {
	model.CoaGroupEntityModel
}
type CoaGroupUpdateResponseDoc struct {
	Body struct {
		Meta res.Meta               `json:"meta"`
		Data CoaGroupUpdateResponse `json:"data"`
	} `json:"body"`
}

// Delete
type CoaGroupDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type CoaGroupDeleteResponse struct {
	// model.CoaGroupEntityModel
}
type CoaGroupDeleteResponseDoc struct {
	Body struct {
		Meta res.Meta               `json:"meta"`
		Data CoaGroupDeleteResponse `json:"data"`
	} `json:"body"`
}
