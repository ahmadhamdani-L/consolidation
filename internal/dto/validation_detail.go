package dto

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
)

// Get
type ValidationDetailGetRequest struct {
	abstraction.Pagination
	model.ValidationDetailFilterModel
}
type ValidationDetailGetResponse struct {
	Datas          []model.ValidationDetailEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type ValidationDetailGetResponseDoc struct {
	Body struct {
		Meta res.Meta                            `json:"meta"`
		Data []model.ValidationDetailEntityModel `json:"data"`
	} `json:"body"`
}

// GetByID
type ValidationDetailGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type ValidationDetailGetByIDResponse struct {
	model.ValidationDetailEntityModel
}
type ValidationDetailGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta                        `json:"meta"`
		Data ValidationDetailGetByIDResponse `json:"data"`
	} `json:"body"`
}
