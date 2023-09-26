package dto

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
	"mime/multipart"
)

// Get
type PembelianPenjualanBerelasiGetRequest struct {
	abstraction.Pagination
	model.PembelianPenjualanBerelasiFilterModel
}
type PembelianPenjualanBerelasiGetResponse struct {
	Datas          []model.PembelianPenjualanBerelasiEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type PembelianPenjualanBerelasiGetResponseDoc struct {
	Body struct {
		Meta res.Meta                                      `json:"meta"`
		Data []model.PembelianPenjualanBerelasiEntityModel `json:"data"`
	} `json:"body"`
}

// GetByID
type PembelianPenjualanBerelasiGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type PembelianPenjualanBerelasiGetByIDResponse struct {
	model.PembelianPenjualanBerelasiEntityModel
}
type PembelianPenjualanBerelasiGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta                                  `json:"meta"`
		Data PembelianPenjualanBerelasiGetByIDResponse `json:"data"`
	} `json:"body"`
}

// Create
type PembelianPenjualanBerelasiCreateRequest struct {
	model.PembelianPenjualanBerelasiEntity
}
type PembelianPenjualanBerelasiCreateResponse struct {
	model.PembelianPenjualanBerelasiEntityModel
}
type PembelianPenjualanBerelasiCreateResponseDoc struct {
	Body struct {
		Meta res.Meta                                 `json:"meta"`
		Data PembelianPenjualanBerelasiCreateResponse `json:"data"`
	} `json:"body"`
}

// Update
type PembelianPenjualanBerelasiUpdateRequest struct {
	ID int `param:"id" validate:"required,numeric"`
	model.PembelianPenjualanBerelasiEntity
}
type PembelianPenjualanBerelasiUpdateResponse struct {
	model.PembelianPenjualanBerelasiEntityModel
}
type PembelianPenjualanBerelasiUpdateResponseDoc struct {
	Body struct {
		Meta res.Meta                                 `json:"meta"`
		Data PembelianPenjualanBerelasiUpdateResponse `json:"data"`
	} `json:"body"`
}

// Delete
type PembelianPenjualanBerelasiDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type PembelianPenjualanBerelasiDeleteResponse struct {
	// model.PembelianPenjualanBerelasiEntityModel
}
type PembelianPenjualanBerelasiDeleteResponseDoc struct {
	Body struct {
		Meta res.Meta                                 `json:"meta"`
		Data PembelianPenjualanBerelasiDeleteResponse `json:"data"`
	} `json:"body"`
}

// export
type PembelianPenjualanBerelasiExportRequest struct {
	// UserID int
	// Period                       string `query:"period" validate:"required"`
	// Versions                     int    `query:"versions" validate:"required"`
	// CompanyID                    int    `query:"company_id"`
	PembelianPenjualanBerelasiID int `query:"pembelian_penjualan_berelasi_id" validate:"required"`
}
type PembelianPenjualanBerelasiExportResponse struct {
	FileName string `json:"filename"`
	Path     string `json:"path"`
}
type PembelianPenjualanBerelasiExportResponseDoc struct {
	Body struct {
		Meta res.Meta                                 `json:"meta"`
		Data PembelianPenjualanBerelasiExportResponse `json:"data"`
	} `json:"body"`
}

// Import
type PembelianPenjualanBerelasiImportRequest struct {
	UserID    int
	CompanyID int
	File      multipart.File
}
type PembelianPenjualanBerelasiImportResponse struct {
	Data model.PembelianPenjualanBerelasiEntityModel
}
type PembelianPenjualanBerelasiImportResponseDoc struct {
	Body struct {
		Meta res.Meta                                 `json:"meta"`
		Data PembelianPenjualanBerelasiImportResponse `json:"data"`
	} `json:"body"`
}
