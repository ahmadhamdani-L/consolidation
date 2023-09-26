package dto

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
)

// Get
type JpmGetRequest struct {
	abstraction.Pagination
	model.JpmFilterModel
}
type JpmGetResponse struct {
	Datas          []model.JpmEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type JpmGetResponseDoc struct {
	Body struct {
		Meta res.Meta               `json:"meta"`
		Data []model.JpmEntityModel `json:"data"`
	} `json:"body"`
}

// GetByID
type JpmGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type JpmGetByIDResponse struct {
	model.JpmEntityModel
}
type JpmGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta           `json:"meta"`
		Data JpmGetByIDResponse `json:"data"`
	} `json:"body"`
}

// Create
type JpmCreateRequest struct {
	model.JpmEntity
	JpmDetail []model.JpmDetailEntity
}
type JpmCreateResponse struct {
	model.JpmEntityModel
}
type JpmCreateResponseDoc struct {
	Body struct {
		Meta res.Meta          `json:"meta"`
		Data JpmCreateResponse `json:"data"`
	} `json:"body"`
}

// Update
type JpmUpdateRequest struct {
	ID int `param:"id" validate:"required,numeric"`
	model.JpmEntity
}
type JpmUpdateResponse struct {
	model.JpmEntityModel
}
type JpmUpdateResponseDoc struct {
	Body struct {
		Meta res.Meta          `json:"meta"`
		Data JpmUpdateResponse `json:"data"`
	} `json:"body"`
}

// Delete
type JpmDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type JpmDeleteResponse struct {
	model.JpmEntityModel
}
type JpmDeleteResponseDoc struct {
	Body struct {
		Meta res.Meta          `json:"meta"`
		Data JpmDeleteResponse `json:"data"`
	} `json:"body"`
}

// export
type JpmExportRequest struct {
	JpmID int `query:"jpm_id" validate:"required"`
}
type JpmExportResponse struct {
	FileName string `json:"filename"`
	Path     string `json:"path"`
}
type JpmExportResponseDoc struct {
	Body struct {
		Meta res.Meta          `json:"meta"`
		Data JpmExportResponse `json:"data"`
	} `json:"body"`
}
