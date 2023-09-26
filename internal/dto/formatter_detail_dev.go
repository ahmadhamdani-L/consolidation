package dto

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
)

// Get
type FormatterDetailDevGetRequest struct {
	abstraction.Pagination
	model.FormatterDetailDevFilterModel
}
type FormatterDetailDevGetResponse struct {
	Datas          []model.FormatterDetailDevEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type FormatterDetailDevGetResponseDoc struct {
	Body struct {
		Meta res.Meta                           `json:"meta"`
		Data []model.FormatterDetailDevEntityModel `json:"data"`
	} `json:"body"`
}

// GetByID
type FormatterDetailDevGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type FormatterDetailDevGetByIDResponse struct {
	model.FormatterDetailDevEntityModel
}
type FormatterDetailDevGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta                       `json:"meta"`
		Data FormatterDetailDevGetByIDResponse `json:"data"`
	} `json:"body"`
}

// Create
type FormatterDetailDevCreateRequest struct {
	model.FormatterDetailDevEntity
	CoaTypeID *int `json:"coa_type_id"`
}
type FormatterDetailDevCreateResponse struct {
	model.FormatterDetailDevEntityModel
}
type FormatterDetailDevCreateResponseDoc struct {
	Body struct {
		Meta res.Meta                      `json:"meta"`
		Data FormatterDetailDevCreateResponse `json:"data"`
	} `json:"body"`
}

// Update
type FormatterDetailDevUpdateRequest struct {
	ID int `param:"id" validate:"required,numeric"`
	model.FormatterDetailDevEntity
}
type FormatterDetailDevUpdateResponse struct {
	model.FormatterDetailDevEntityModel
}
type FormatterDetailDevUpdateResponseDoc struct {
	Body struct {
		Meta res.Meta                      `json:"meta"`
		Data FormatterDetailDevUpdateResponse `json:"data"`
	} `json:"body"`
}

// Delete
type FormatterDetailDevDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type FormatterDetailDevDeleteResponse struct {
	// model.FormatterDetailDevEntityModel
}
type FormatterDetailDevDeleteResponseDoc struct {
	Body struct {
		Meta res.Meta                      `json:"meta"`
		Data FormatterDetailDevDeleteResponse `json:"data"`
	} `json:"body"`
}

// export
type FormatterDetailDevExportRequest struct {
	abstraction.Pagination
	model.FormatterDetailDevFilterModel
}
type FormatterDetailDevExportResponse struct {
	File string `json:"file"`
}
type FormatterDetailDevExportResponseDoc struct {
	Body struct {
		Meta res.Meta                      `json:"meta"`
		Data FormatterDetailDevExportResponse `json:"data"`
	} `json:"body"`
}

// Import
type FormatterDetailDevImportRequest struct {
	UserID    int
	CompanyID int
}
type FormatterDetailDevImportResponse struct {
	Data model.FormatterDetailDevEntityModel
}
type FormatterDetailDevImportResponseDoc struct {
	Body struct {
		Meta res.Meta                           `json:"meta"`
		Data []model.FormatterDetailDevEntityModel `json:"data"`
	} `json:"body"`
}

type FormatterDetailDevDetailGetByParentRequests struct {
	Datas []model.FormatterDetailDevFmtEntityModel
}
type FormatterDetailDevDetailGetByParentResponses struct {
	Data []model.FormatterDetailDevFmtEntityModel
	// Data map[int]*model.TrialBalanceDetailFmtEntityModel
}

type FormatterDragAndDropDevRequest struct {
	ParentID int `param:"parent_id" validate:"required,numeric"`
	Datas []model.FormatterDetailDevFmtEntityModel `json:"Datas"`
}
type FormatterDragAndDropDevResponse struct {
	Datas          []model.FormatterDetailDevFmtEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type FormatterDragAndDropDevResponseDoc struct {
	Body struct {
		Meta res.Meta          `json:"meta"`
		Data FormatterDragAndDropResponse `json:"data"`
	} `json:"body"`
}

type FormatterDetailDevFmtEntityModel struct {
    ID       int              `json:"id"`
    SortID   float64              `json:"sort_id"`
    Children []model.FormatterDetailDevFmtEntityModel  `json:"children"`
}

type FormatterDetailDevGetCoaRequest struct {
	model.FormatterDetailDevFilterModel
}
type FormatterDetailDevGetCoaResponse struct {
	Datas []model.CoaDevEntityModel
	abstraction.Pagination
}
type FormatterDetailDevGetCoaResponseDoc struct {
	Body struct {
		Meta res.Meta                      `json:"meta"`
		Data FormatterDetailDevGetCoaResponse `json:"data"`
	} `json:"body"`
}