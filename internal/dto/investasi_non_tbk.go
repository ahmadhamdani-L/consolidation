package dto

import (
	"worker/internal/abstraction"
	"worker/internal/model"
	res "worker/pkg/util/response"
	"mime/multipart"
)

// Get
type InvestasiNonTbkGetRequest struct {
	abstraction.Pagination
	model.InvestasiNonTbkFilterModel
}
type InvestasiNonTbkGetResponse struct {
	Datas          []model.InvestasiNonTbkEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type InvestasiNonTbkGetResponseDoc struct {
	Body struct {
		Meta res.Meta                           `json:"meta"`
		Data []model.InvestasiNonTbkEntityModel `json:"data"`
	} `json:"body"`
}

// GetByID
type InvestasiNonTbkGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type InvestasiNonTbkGetByIDResponse struct {
	model.InvestasiNonTbkEntityModel
}
type InvestasiNonTbkGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta                       `json:"meta"`
		Data InvestasiNonTbkGetByIDResponse `json:"data"`
	} `json:"body"`
}

// Create
type InvestasiNonTbkCreateRequest struct {
	model.InvestasiNonTbkEntity
}
type InvestasiNonTbkCreateResponse struct {
	model.InvestasiNonTbkEntityModel
}
type InvestasiNonTbkCreateResponseDoc struct {
	Body struct {
		Meta res.Meta                      `json:"meta"`
		Data InvestasiNonTbkCreateResponse `json:"data"`
	} `json:"body"`
}

// Update
type InvestasiNonTbkUpdateRequest struct {
	ID int `param:"id" validate:"required,numeric"`
	model.InvestasiNonTbkEntity
}
type InvestasiNonTbkUpdateResponse struct {
	model.InvestasiNonTbkEntityModel
}
type InvestasiNonTbkUpdateResponseDoc struct {
	Body struct {
		Meta res.Meta                      `json:"meta"`
		Data InvestasiNonTbkUpdateResponse `json:"data"`
	} `json:"body"`
}

// Delete
type InvestasiNonTbkDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type InvestasiNonTbkDeleteResponse struct {
	model.InvestasiNonTbkEntityModel
}
type InvestasiNonTbkDeleteResponseDoc struct {
	Body struct {
		Meta res.Meta                      `json:"meta"`
		Data InvestasiNonTbkDeleteResponse `json:"data"`
	} `json:"body"`
}

// export
type InvestasiNonTbkExportRequest struct {
	UserId    int
	Period    string `query:"period" validate:"required"`
	Versions  int    `query:"version" validate:"required"`
	CompanyID int    `query:"company_id"`
}
type InvestasiNonTbkExportResponse struct {
	File string `json:"file"`
}
type InvestasiNonTbkExportResponseDoc struct {
	Body struct {
		Meta res.Meta                      `json:"meta"`
		Data InvestasiNonTbkExportResponse `json:"data"`
	} `json:"body"`
}

// Import
type InvestasiNonTbkImportRequest struct {
	UserId    int
	CompanyId int
	File      multipart.File
}
type InvestasiNonTbkImportResponse struct {
	Datas          []model.InvestasiNonTbkEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type InvestasiNonTbkImportResponseDoc struct {
	Body struct {
		Meta res.Meta                           `json:"meta"`
		Data []model.InvestasiNonTbkEntityModel `json:"data"`
	} `json:"body"`
}
