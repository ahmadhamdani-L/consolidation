package dto

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
)

// Get
type AdjustmentGetRequest struct {
	abstraction.Pagination
	model.AdjustmentFilterModel
}
type AdjustmentGetResponse struct {
	Datas          []model.AdjustmentEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type AdjustmentGetResponseDoc struct {
	Body struct {
		Meta res.Meta                      `json:"meta"`
		Data []model.AdjustmentEntityModel `json:"data"`
	} `json:"body"`
}

// GetByID
type AdjustmentGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type AdjustmentGetByIDResponse struct {
	model.AdjustmentEntityModel
}
type AdjustmentGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta                  `json:"meta"`
		Data AdjustmentGetByIDResponse `json:"data"`
	} `json:"body"`
}

// Create
type AdjustmentCreateRequest struct {
	model.AdjustmentEntity
	AdjustmentDetail []model.AdjustmentDetailEntity
}
type AdjustmentCreateResponse struct {
	model.AdjustmentEntityModel
}
type AdjustmentCreateResponseDoc struct {
	Body struct {
		Meta res.Meta                 `json:"meta"`
		Data AdjustmentCreateResponse `json:"data"`
	} `json:"body"`
}

// Update
type AdjustmentUpdateRequest struct {
	ID int `param:"id" validate:"required,numeric"`
	model.AdjustmentEntity
}
type AdjustmentUpdateResponse struct {
	model.AdjustmentEntityModel
}
type AdjustmentUpdateResponseDoc struct {
	Body struct {
		Meta res.Meta                 `json:"meta"`
		Data AdjustmentUpdateResponse `json:"data"`
	} `json:"body"`
}

// Delete
type AdjustmentDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type AdjustmentDeleteResponse struct {
	model.AdjustmentEntityModel
}
type AdjustmentDeleteResponseDoc struct {
	Body struct {
		Meta res.Meta                 `json:"meta"`
		Data AdjustmentDeleteResponse `json:"data"`
	} `json:"body"`
}

// export
type AdjustmentExportRequest struct {
	AdjustmentID int `query:"adjustment_id" validate:"required"`
}
type AdjustmentExportResponse struct {
	FileName string `json:"filename"`
	Path     string `json:"path"`
}
type AdjustmentExportResponseDoc struct {
	Body struct {
		Meta res.Meta                 `json:"meta"`
		Data AdjustmentExportResponse `json:"data"`
	} `json:"body"`
}
