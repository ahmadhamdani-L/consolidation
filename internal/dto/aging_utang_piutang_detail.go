package dto

import (
	"worker/internal/abstraction"
	"worker/internal/model"
	res "worker/pkg/util/response"
)

// Get
type AgingUtangPiutangDetailGetRequest struct {
	abstraction.Pagination
	model.AgingUtangPiutangDetailFilterModel
}
type AgingUtangPiutangDetailGetResponse struct {
	Datas          []model.AgingUtangPiutangDetailEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type AgingUtangPiutangDetailGetResponseDoc struct {
	Body struct {
		Meta res.Meta                                   `json:"meta"`
		Data []model.AgingUtangPiutangDetailEntityModel `json:"data"`
	} `json:"body"`
}

// GetByID
type AgingUtangPiutangDetailGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type AgingUtangPiutangDetailGetByIDResponse struct {
	model.AgingUtangPiutangDetailEntityModel
}
type AgingUtangPiutangDetailGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta                               `json:"meta"`
		Data AgingUtangPiutangDetailGetByIDResponse `json:"data"`
	} `json:"body"`
}

// Create
type AgingUtangPiutangDetailCreateRequest struct {
	model.AgingUtangPiutangDetailEntity
}
type AgingUtangPiutangDetailCreateResponse struct {
	model.AgingUtangPiutangDetailEntityModel
}
type AgingUtangPiutangDetailCreateResponseDoc struct {
	Body struct {
		Meta res.Meta                              `json:"meta"`
		Data AgingUtangPiutangDetailCreateResponse `json:"data"`
	} `json:"body"`
}

// Update
type AgingUtangPiutangDetailUpdateRequest struct {
	ID int `param:"id" validate:"required,numeric"`
	model.AgingUtangPiutangDetailEntity
}
type AgingUtangPiutangDetailUpdateResponse struct {
	model.AgingUtangPiutangDetailEntityModel
}
type AgingUtangPiutangDetailUpdateResponseDoc struct {
	Body struct {
		Meta res.Meta                              `json:"meta"`
		Data AgingUtangPiutangDetailUpdateResponse `json:"data"`
	} `json:"body"`
}

// Delete
type AgingUtangPiutangDetailDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type AgingUtangPiutangDetailDeleteResponse struct {
	model.AgingUtangPiutangDetailEntityModel
}
type AgingUtangPiutangDetailDeleteResponseDoc struct {
	Body struct {
		Meta res.Meta                              `json:"meta"`
		Data AgingUtangPiutangDetailDeleteResponse `json:"data"`
	} `json:"body"`
}

// export
type AgingUtangPiutangDetailExportRequest struct {
	model.AgingUtangPiutangDetailFilterModel
}
type AgingUtangPiutangDetailExportResponse struct {
	File string `json:"file"`
}
type AgingUtangPiutangDetailExportResponseDoc struct {
	Body struct {
		Meta res.Meta                              `json:"meta"`
		Data AgingUtangPiutangDetailExportResponse `json:"data"`
	} `json:"body"`
}

// Import
type AgingUtangPiutangDetailImportRequest struct {
	Datas []model.AgingUtangPiutangDetailEntity
}
type AgingUtangPiutangDetailImportResponse struct {
	Datas          []model.AgingUtangPiutangDetailEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type AgingUtangPiutangDetailImportResponseDoc struct {
	Body struct {
		Meta res.Meta                                   `json:"meta"`
		Data []model.AgingUtangPiutangDetailEntityModel `json:"data"`
	} `json:"body"`
}
