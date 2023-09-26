package dto

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
)

// Get
type MutasiDtaDetailGetRequest struct {
	abstraction.Pagination
	model.MutasiDtaDetailFilterModel
}
type MutasiDtaDetailGetResponse struct {
	Datas          model.MutasiDtaEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type MutasiDtaDetailGetResponseDoc struct {
	Body struct {
		Meta res.Meta                   `json:"meta"`
		Data MutasiDtaDetailGetResponse `json:"data"`
	} `json:"body"`
}

// GetByID
type MutasiDtaDetailGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type MutasiDtaDetailGetByIDResponse struct {
	model.MutasiDtaDetailEntityModel
}
type MutasiDtaDetailGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta                       `json:"meta"`
		Data MutasiDtaDetailGetByIDResponse `json:"data"`
	} `json:"body"`
}

// Create
type MutasiDtaDetailCreateRequest struct {
	model.MutasiDtaDetailEntity
}
type MutasiDtaDetailCreateResponse struct {
	model.MutasiDtaDetailEntityModel
}
type MutasiDtaDetailCreateResponseDoc struct {
	Body struct {
		Meta res.Meta                      `json:"meta"`
		Data MutasiDtaDetailCreateResponse `json:"data"`
	} `json:"body"`
}

// Update
type MutasiDtaDetailUpdateRequest struct {
	ID                  int      `param:"id" validate:"required,numeric"`
	SaldoAwal           *float64 `json:"saldo_awal" validate:"required" example:"10000.00"`
	ManfaatBebanPajak   *float64 `json:"manfaat_beban_pajak" validate:"required" example:"10000.00"`
	Oci                 *float64 `json:"oci" validate:"required" example:"1000.00"`
	AkuisisiEntitasAnak *float64 `json:"akuisisi_entitas_anak" validate:"required" example:"10000.00"`
	DibebankanKeLr      *float64 `json:"dibebankan_ke_lr" validate:"required" example:"10000.00"`
	DibebankanKeOci     *float64 `json:"dibebankan_ke_oci" validate:"required" example:"10000.00"`
}
type MutasiDtaDetailUpdateResponse struct {
	model.MutasiDtaDetailEntityModel
}
type MutasiDtaDetailUpdateResponseDoc struct {
	Body struct {
		Meta res.Meta                      `json:"meta"`
		Data MutasiDtaDetailUpdateResponse `json:"data"`
	} `json:"body"`
}

// Delete
type MutasiDtaDetailDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type MutasiDtaDetailDeleteResponse struct {
	model.MutasiDtaDetailEntityModel
}
type MutasiDtaDetailDeleteResponseDoc struct {
	Body struct {
		Meta res.Meta                      `json:"meta"`
		Data MutasiDtaDetailDeleteResponse `json:"data"`
	} `json:"body"`
}
