package dto

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
)

// Get
type FormatterDetailGetRequest struct {
	abstraction.Pagination
	model.FormatterDetailFilterModel
}
type FormatterDetailGetResponse struct {
	Datas          []model.FormatterDetailEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type FormatterDetailGetResponseDoc struct {
	Body struct {
		Meta res.Meta                           `json:"meta"`
		Data []model.FormatterDetailEntityModel `json:"data"`
	} `json:"body"`
}

// GetByID
type FormatterDetailGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type FormatterDetailGetByIDResponse struct {
	model.FormatterDetailEntityModel
}
type FormatterDetailGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta                       `json:"meta"`
		Data FormatterDetailGetByIDResponse `json:"data"`
	} `json:"body"`
}

// Create
type FormatterDetailCreateRequest struct {
	model.FormatterDetailEntity
	CoaTypeID *int `json:"coa_type_id"`
}
type FormatterDetailCreateResponse struct {
	model.FormatterDetailEntityModel
}
type FormatterDetailCreateResponseDoc struct {
	Body struct {
		Meta res.Meta                      `json:"meta"`
		Data FormatterDetailCreateResponse `json:"data"`
	} `json:"body"`
}

// Update
type FormatterDetailUpdateRequest struct {
	ID int `param:"id" validate:"required,numeric"`
	model.FormatterDetailEntity
}
type FormatterDetailUpdateResponse struct {
	model.FormatterDetailEntityModel
}
type FormatterDetailUpdateResponseDoc struct {
	Body struct {
		Meta res.Meta                      `json:"meta"`
		Data FormatterDetailUpdateResponse `json:"data"`
	} `json:"body"`
}

// Delete
type FormatterDetailDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type FormatterDetailDeleteResponse struct {
	// model.FormatterDetailEntityModel
}
type FormatterDetailDeleteResponseDoc struct {
	Body struct {
		Meta res.Meta                      `json:"meta"`
		Data FormatterDetailDeleteResponse `json:"data"`
	} `json:"body"`
}

// export
type FormatterDetailExportRequest struct {
	abstraction.Pagination
	model.FormatterDetailFilterModel
}
type FormatterDetailExportResponse struct {
	File string `json:"file"`
}
type FormatterDetailExportResponseDoc struct {
	Body struct {
		Meta res.Meta                      `json:"meta"`
		Data FormatterDetailExportResponse `json:"data"`
	} `json:"body"`
}

// Import
type FormatterDetailImportRequest struct {
	UserID    int
	CompanyID int
}
type FormatterDetailImportResponse struct {
	Data model.FormatterDetailEntityModel
}
type FormatterDetailImportResponseDoc struct {
	Body struct {
		Meta res.Meta                           `json:"meta"`
		Data []model.FormatterDetailEntityModel `json:"data"`
	} `json:"body"`
}

type FormatterDetailDetailGetByParentRequests struct {
	Datas []model.FormatterDetailFmtEntityModel
}
type FormatterDetailDetailGetByParentResponses struct {
	Data []model.FormatterDetailFmtEntityModel
	// Data map[int]*model.TrialBalanceDetailFmtEntityModel
}

type FormatterDragAndDropRequest struct {
	ParentID int `param:"parent_id" validate:"required,numeric"`
	Datas []model.FormatterDetailFmtEntityModel `json:"Datas"`
}
type FormatterDragAndDropResponse struct {
	Datas          []model.FormatterDetailFmtEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type FormatterDragAndDropResponseDoc struct {
	Body struct {
		Meta res.Meta          `json:"meta"`
		Data FormatterDragAndDropResponse `json:"data"`
	} `json:"body"`
}

type FormatterDetailFmtEntityModel struct {
    ID       int              `json:"id"`
    SortID   float64              `json:"sort_id"`
    Children []model.FormatterDetailFmtEntityModel  `json:"children"`
}

type FormatterDetailGetCoaRequest struct {
	model.FormatterDetailFilterModel
}
type FormatterDetailGetCoaResponse struct {
	Datas []model.CoaEntityModel
	abstraction.Pagination
}
type FormatterDetailGetCoaResponseDoc struct {
	Body struct {
		Meta res.Meta                      `json:"meta"`
		Data FormatterDetailGetCoaResponse `json:"data"`
	} `json:"body"`
}