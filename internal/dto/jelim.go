package dto

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
)

// Get
type JelimGetRequest struct {
	abstraction.Pagination
	model.JelimFilterModel
}
type JelimGetResponse struct {
	Datas          []model.JelimEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type JelimGetResponseDoc struct {
	Body struct {
		Meta res.Meta                 `json:"meta"`
		Data []model.JelimEntityModel `json:"data"`
	} `json:"body"`
}

// GetByID
type JelimGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type JelimGetByIDResponse struct {
	model.JelimEntityModel
}
type JelimGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta             `json:"meta"`
		Data JelimGetByIDResponse `json:"data"`
	} `json:"body"`
}

// Create
type JelimCreateRequest struct {
	model.JelimEntity
	JelimDetail []model.JelimDetailEntity
}
type JelimCreateResponse struct {
	model.JelimEntityModel
}
type JelimCreateResponseDoc struct {
	Body struct {
		Meta res.Meta            `json:"meta"`
		Data JelimCreateResponse `json:"data"`
	} `json:"body"`
}

// Update
type JelimUpdateRequest struct {
	ID int `param:"id" validate:"required,numeric"`
	model.JelimEntity
}
type JelimUpdateResponse struct {
	model.JelimEntityModel
}
type JelimUpdateResponseDoc struct {
	Body struct {
		Meta res.Meta            `json:"meta"`
		Data JelimUpdateResponse `json:"data"`
	} `json:"body"`
}

// Delete
type JelimDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type JelimDeleteResponse struct {
	model.JelimEntityModel
}
type JelimDeleteResponseDoc struct {
	Body struct {
		Meta res.Meta            `json:"meta"`
		Data JelimDeleteResponse `json:"data"`
	} `json:"body"`
}

// export
type JelimExportRequest struct {
	JelimID int `query:"jelim_id" validate:"required"`
}
type JelimExportResponse struct {
	FileName string `json:"filename"`
	Path     string `json:"path"`
}
type JelimExportResponseDoc struct {
	Body struct {
		Meta res.Meta            `json:"meta"`
		Data JelimExportResponse `json:"data"`
	} `json:"body"`
}
