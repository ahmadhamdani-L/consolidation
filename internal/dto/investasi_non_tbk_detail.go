package dto

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
)

// Get
type InvestasiNonTbkDetailGetRequest struct {
	abstraction.Pagination
	model.InvestasiNonTbkDetailFilterModel
}
type InvestasiNonTbkDetailGetResponse struct {
	Datas          model.InvestasiNonTbkEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type InvestasiNonTbkDetailGetResponseDoc struct {
	Body struct {
		Meta res.Meta                         `json:"meta"`
		Data InvestasiNonTbkDetailGetResponse `json:"data"`
	} `json:"body"`
}

// GetByID
type InvestasiNonTbkDetailGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type InvestasiNonTbkDetailGetByIDResponse struct {
	model.InvestasiNonTbkDetailEntityModel
}
type InvestasiNonTbkDetailGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta                             `json:"meta"`
		Data InvestasiNonTbkDetailGetByIDResponse `json:"data"`
	} `json:"body"`
}

// Create
type InvestasiNonTbkDetailCreateRequest struct {
	model.InvestasiNonTbkDetailEntity
}
type InvestasiNonTbkDetailCreateResponse struct {
	model.InvestasiNonTbkDetailEntityModel
}
type InvestasiNonTbkDetailCreateResponseDoc struct {
	Body struct {
		Meta res.Meta                            `json:"meta"`
		Data InvestasiNonTbkDetailCreateResponse `json:"data"`
	} `json:"body"`
}

// Update
type InvestasiNonTbkDetailUpdateRequest struct {
	ID                int      `param:"id" validate:"required,numeric"`
	LbrSahamOwnership *float64 `json:"lbr_saham_ownership" validate:"required" example:"10"`
	TotalLbrSaham     *float64 `json:"total_lbr_saham" validate:"required" example:"10000.00" min:"1"`
	HargaPar          *float64 `json:"harga_par" validate:"required" example:"10000.00"`
	HargaBeli         *float64 `json:"harga_beli" validate:"required" example:"10000.00"`
	// TotalHargaPar     *float64 `json:"total_harga_par" validate:"required" example:"10000.00"`
	// PercentageOwnership *float64 `json:"percentage_ownership" validate:"required" example:"10"`
	// TotalHargaBeli *float64 `json:"total_harga_beli" validate:"required" example:"10000.00"`
}
type InvestasiNonTbkDetailUpdateResponse struct {
	model.InvestasiNonTbkDetailEntityModel
}
type InvestasiNonTbkDetailUpdateResponseDoc struct {
	Body struct {
		Meta res.Meta                            `json:"meta"`
		Data InvestasiNonTbkDetailUpdateResponse `json:"data"`
	} `json:"body"`
}

// Delete
type InvestasiNonTbkDetailDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type InvestasiNonTbkDetailDeleteResponse struct {
	model.InvestasiNonTbkDetailEntityModel
}
type InvestasiNonTbkDetailDeleteResponseDoc struct {
	Body struct {
		Meta res.Meta                            `json:"meta"`
		Data InvestasiNonTbkDetailDeleteResponse `json:"data"`
	} `json:"body"`
}
