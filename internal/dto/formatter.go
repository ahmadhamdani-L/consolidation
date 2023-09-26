package dto

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
)

// Get
type FormatterGetRequest struct {
	abstraction.Pagination
	model.FormatterFilterModel
}
type FormatterGetResponse struct {
	Datas          []model.FormatterEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type FormatterGetResponseDoc struct {
	Body struct {
		Meta res.Meta                     `json:"meta"`
		Data []model.FormatterEntityModel `json:"data"`
	} `json:"body"`
}

// GetByID
type FormatterGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type FormatterGetByIDResponse struct {
	model.FormatterEntityModel
}
type FormatterGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta                 `json:"meta"`
		Data FormatterGetByIDResponse `json:"data"`
	} `json:"body"`
}

// Create
type FormatterCreateRequest struct {
	model.FormatterEntity
}
type FormatterCreateResponse struct {
	model.FormatterEntityModel
}
type FormatterCreateResponseDoc struct {
	Body struct {
		Meta res.Meta                `json:"meta"`
		Data FormatterCreateResponse `json:"data"`
	} `json:"body"`
}

// Update
type FormatterUpdateRequest struct {
	ID int `param:"id" validate:"required,numeric"`
	model.FormatterEntity
}
type FormatterUpdateResponse struct {
	model.FormatterEntityModel
}
type FormatterUpdateResponseDoc struct {
	Body struct {
		Meta res.Meta                `json:"meta"`
		Data FormatterUpdateResponse `json:"data"`
	} `json:"body"`
}

// Delete
type FormatterDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type FormatterDeleteResponse struct {
	// model.FormatterEntityModel
}
type FormatterDeleteResponseDoc struct {
	Body struct {
		Meta res.Meta                `json:"meta"`
		Data FormatterDeleteResponse `json:"data"`
	} `json:"body"`
}

// export
type FormatterExportRequest struct {
	abstraction.Pagination
	model.FormatterFilterModel
}
type FormatterExportResponse struct {
	File string `json:"file"`
}
type FormatterExportResponseDoc struct {
	Body struct {
		Meta res.Meta                `json:"meta"`
		Data FormatterExportResponse `json:"data"`
	} `json:"body"`
}

// Import
type FormatterImportRequest struct {
	UserID    int
	CompanyID int
}
type FormatterImportResponse struct {
	Data model.FormatterEntityModel
}
type FormatterImportResponseDoc struct {
	Body struct {
		Meta res.Meta                     `json:"meta"`
		Data []model.FormatterEntityModel `json:"data"`
	} `json:"body"`
}
