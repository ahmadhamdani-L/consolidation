package dto

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
)

// Get
type JcteGetRequest struct {
	abstraction.Pagination
	model.JcteFilterModel
}
type JcteGetResponse struct {
	Datas          []model.JcteEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type JcteGetResponseDoc struct {
	Body struct {
		Meta res.Meta                `json:"meta"`
		Data []model.JcteEntityModel `json:"data"`
	} `json:"body"`
}

// GetByID
type JcteGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type JcteGetByIDResponse struct {
	model.JcteEntityModel
}
type JcteGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta            `json:"meta"`
		Data JcteGetByIDResponse `json:"data"`
	} `json:"body"`
}

// Create
type JcteCreateRequest struct {
	model.JcteEntity
	JcteDetail []model.JcteDetailEntity
}
type JcteCreateResponse struct {
	model.JcteEntityModel
}
type JcteCreateResponseDoc struct {
	Body struct {
		Meta res.Meta           `json:"meta"`
		Data JcteCreateResponse `json:"data"`
	} `json:"body"`
}

// Update
type JcteUpdateRequest struct {
	ID int `param:"id" validate:"required,numeric"`
	model.JcteEntity
}
type JcteUpdateResponse struct {
	model.JcteEntityModel
}
type JcteUpdateResponseDoc struct {
	Body struct {
		Meta res.Meta           `json:"meta"`
		Data JcteUpdateResponse `json:"data"`
	} `json:"body"`
}

// Delete
type JcteDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type JcteDeleteResponse struct {
	model.JcteEntityModel
}
type JcteDeleteResponseDoc struct {
	Body struct {
		Meta res.Meta           `json:"meta"`
		Data JcteDeleteResponse `json:"data"`
	} `json:"body"`
}

// export
type JcteExportRequest struct {
	JcteID int `query:"jcte_id" validate:"required"`
}
type JcteExportResponse struct {
	FileName string `json:"filename"`
	Path     string `json:"path"`
}
type JcteExportResponseDoc struct {
	Body struct {
		Meta res.Meta           `json:"meta"`
		Data JcteExportResponse `json:"data"`
	} `json:"body"`
}
