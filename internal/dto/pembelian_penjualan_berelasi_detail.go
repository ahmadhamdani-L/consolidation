package dto

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
)

// Get
type PembelianPenjualanBerelasiDetailGetRequest struct {
	abstraction.Pagination
	model.PembelianPenjualanBerelasiDetailFilterModel
}
type PembelianPenjualanBerelasiDetailGetResponse struct {
	Datas struct {
		TotalPenjualan float64                                     `json:"total_penjualan"`
		TotalPembelian float64                                     `json:"total_pembelian"`
		Data           model.PembelianPenjualanBerelasiEntityModel `json:"pembelian_penjualan_berelasi"`
	}
	PaginationInfo abstraction.PaginationInfo
}
type PembelianPenjualanBerelasiDetailGetResponseDoc struct {
	Body struct {
		Meta res.Meta                                    `json:"meta"`
		Data PembelianPenjualanBerelasiDetailGetResponse `json:"data"`
	} `json:"body"`
}

// GetByID
type PembelianPenjualanBerelasiDetailGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type PembelianPenjualanBerelasiDetailGetByIDResponse struct {
	model.PembelianPenjualanBerelasiDetailEntityModel
}
type PembelianPenjualanBerelasiDetailGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta                                        `json:"meta"`
		Data PembelianPenjualanBerelasiDetailGetByIDResponse `json:"data"`
	} `json:"body"`
}

// Create
type PembelianPenjualanBerelasiDetailCreateRequest struct {
	model.PembelianPenjualanBerelasiDetailEntity
}
type PembelianPenjualanBerelasiDetailCreateResponse struct {
	model.PembelianPenjualanBerelasiDetailEntityModel
}
type PembelianPenjualanBerelasiDetailCreateResponseDoc struct {
	Body struct {
		Meta res.Meta                                       `json:"meta"`
		Data PembelianPenjualanBerelasiDetailCreateResponse `json:"data"`
	} `json:"body"`
}

// Update
type PembelianPenjualanBerelasiDetailUpdateRequest struct {
	ID           int      `param:"id" validate:"required,numeric"`
	BoughtAmount *float64 `json:"bought_amount" validate:"required" example:"10000.00"`
	SalesAmount  *float64 `json:"sales_amount" validate:"required" example:"10000.00"`
}
type PembelianPenjualanBerelasiDetailUpdateResponse struct {
	model.PembelianPenjualanBerelasiDetailEntityModel
}
type PembelianPenjualanBerelasiDetailUpdateResponseDoc struct {
	Body struct {
		Meta res.Meta                                       `json:"meta"`
		Data PembelianPenjualanBerelasiDetailUpdateResponse `json:"data"`
	} `json:"body"`
}

// Delete
type PembelianPenjualanBerelasiDetailDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type PembelianPenjualanBerelasiDetailDeleteResponse struct {
}
type PembelianPenjualanBerelasiDetailDeleteResponseDoc struct {
	Body struct {
		Meta res.Meta                                       `json:"meta"`
		Data PembelianPenjualanBerelasiDetailDeleteResponse `json:"data"`
	} `json:"body"`
}
