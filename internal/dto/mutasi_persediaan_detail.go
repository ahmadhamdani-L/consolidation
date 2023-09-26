package dto

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
)

// Get
type MutasiPersediaanDetailGetRequest struct {
	abstraction.Pagination
	model.MutasiPersediaanDetailFilterModel
}
type MutasiPersediaanDetailGetResponse struct {
	Datas          model.MutasiPersediaanEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type MutasiPersediaanDetailGetResponseDoc struct {
	Body struct {
		Meta res.Meta                          `json:"meta"`
		Data MutasiPersediaanDetailGetResponse `json:"data"`
	} `json:"body"`
}

// GetByID
type MutasiPersediaanDetailGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type MutasiPersediaanDetailGetByIDResponse struct {
	model.MutasiPersediaanDetailEntityModel
}
type MutasiPersediaanDetailGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta                              `json:"meta"`
		Data MutasiPersediaanDetailGetByIDResponse `json:"data"`
	} `json:"body"`
}

// Create
type MutasiPersediaanDetailCreateRequest struct {
	model.MutasiPersediaanDetailEntity
}
type MutasiPersediaanDetailCreateResponse struct {
	model.MutasiPersediaanDetailEntityModel
}
type MutasiPersediaanDetailCreateResponseDoc struct {
	Body struct {
		Meta res.Meta                             `json:"meta"`
		Data MutasiPersediaanDetailCreateResponse `json:"data"`
	} `json:"body"`
}

// Update
type MutasiPersediaanDetailUpdateRequest struct {
	ID     int      `param:"id" validate:"required,numeric"`
	Amount *float64 `json:"amount" validate:"required" example:"10000.00"`
}
type MutasiPersediaanDetailUpdateResponse struct {
	model.MutasiPersediaanDetailEntityModel
}
type MutasiPersediaanDetailUpdateResponseDoc struct {
	Body struct {
		Meta res.Meta                             `json:"meta"`
		Data MutasiPersediaanDetailUpdateResponse `json:"data"`
	} `json:"body"`
}

// Delete
type MutasiPersediaanDetailDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type MutasiPersediaanDetailDeleteResponse struct {
	// model.MutasiPersediaanDetailEntityModel
}
type MutasiPersediaanDetailDeleteResponseDoc struct {
	Body struct {
		Meta res.Meta                             `json:"meta"`
		Data MutasiPersediaanDetailDeleteResponse `json:"data"`
	} `json:"body"`
}
