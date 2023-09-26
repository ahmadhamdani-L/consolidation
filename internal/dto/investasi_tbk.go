package dto

import (
	"worker/internal/abstraction"
	"worker/internal/model"
	res "worker/pkg/util/response"
	"mime/multipart"
)

// Get
type InvestasiTbkGetRequest struct {
	abstraction.Pagination
	model.InvestasiTbkFilterModel
}
type InvestasiTbkGetResponse struct {
	Datas          []model.InvestasiTbkEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type InvestasiTbkGetResponseDoc struct {
	Body struct {
		Meta res.Meta                        `json:"meta"`
		Data []model.InvestasiTbkEntityModel `json:"data"`
	} `json:"body"`
}

// GetByID
type InvestasiTbkGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type InvestasiTbkGetByIDResponse struct {
	model.InvestasiTbkEntityModel
}
type InvestasiTbkGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta                    `json:"meta"`
		Data InvestasiTbkGetByIDResponse `json:"data"`
	} `json:"body"`
}

// Create
type InvestasiTbkCreateRequest struct {
	model.InvestasiTbkEntity
}
type InvestasiTbkCreateResponse struct {
	model.InvestasiTbkEntityModel
}
type InvestasiTbkCreateResponseDoc struct {
	Body struct {
		Meta res.Meta                   `json:"meta"`
		Data InvestasiTbkCreateResponse `json:"data"`
	} `json:"body"`
}

// Update
type InvestasiTbkUpdateRequest struct {
	ID int `param:"id" validate:"required,numeric"`
	model.InvestasiTbkEntity
}
type InvestasiTbkUpdateResponse struct {
	model.InvestasiTbkEntityModel
}
type InvestasiTbkUpdateResponseDoc struct {
	Body struct {
		Meta res.Meta                   `json:"meta"`
		Data InvestasiTbkUpdateResponse `json:"data"`
	} `json:"body"`
}

// Delete
type InvestasiTbkDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type InvestasiTbkDeleteResponse struct {
	model.InvestasiTbkEntityModel
}
type InvestasiTbkDeleteResponseDoc struct {
	Body struct {
		Meta res.Meta                   `json:"meta"`
		Data InvestasiTbkDeleteResponse `json:"data"`
	} `json:"body"`
}

// export
type InvestasiTbkExportRequest struct {
	UserId    int
	Period    string `query:"period" validate:"required"`
	Versions  int    `query:"versions" validate:"required"`
	CompanyID int    `query:"company_id"`
}
type InvestasiTbkExportResponse struct {
	File string `json:"file"`
}
type InvestasiTbkExportResponseDoc struct {
	Body struct {
		Meta res.Meta                   `json:"meta"`
		Data InvestasiTbkExportResponse `json:"data"`
	} `json:"body"`
}

// Import
type InvestasiTbkImportRequest struct {
	UserId    int
	CompanyId int
	File      multipart.File
}
type InvestasiTbkImportResponse struct {
	Data           []model.InvestasiTbkEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type InvestasiTbkImportResponseDoc struct {
	Body struct {
		Meta res.Meta                        `json:"meta"`
		Data []model.InvestasiTbkEntityModel `json:"data"`
	} `json:"body"`
}
