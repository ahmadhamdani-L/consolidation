package dto

import (
	"worker/internal/abstraction"
	"worker/internal/model"
	res "worker/pkg/util/response"
	"mime/multipart"
)

// Get
type AgingUtangPiutangGetRequest struct {
	abstraction.Pagination
	model.AgingUtangPiutangFilterModel
}
type AgingUtangPiutangGetResponse struct {
	Datas          []model.AgingUtangPiutangEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type AgingUtangPiutangGetResponseDoc struct {
	Body struct {
		Meta res.Meta                             `json:"meta"`
		Data []model.AgingUtangPiutangEntityModel `json:"data"`
	} `json:"body"`
}

// GetByID
type AgingUtangPiutangGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type AgingUtangPiutangGetByIDResponse struct {
	model.AgingUtangPiutangEntityModel
}
type AgingUtangPiutangGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta                         `json:"meta"`
		Data AgingUtangPiutangGetByIDResponse `json:"data"`
	} `json:"body"`
}

// Create
type AgingUtangPiutangCreateRequest struct {
	model.AgingUtangPiutangEntity
}
type AgingUtangPiutangCreateResponse struct {
	model.AgingUtangPiutangEntityModel
}
type AgingUtangPiutangCreateResponseDoc struct {
	Body struct {
		Meta res.Meta                        `json:"meta"`
		Data AgingUtangPiutangCreateResponse `json:"data"`
	} `json:"body"`
}

// Update
type AgingUtangPiutangUpdateRequest struct {
	ID int `param:"id" validate:"required,numeric"`
	model.AgingUtangPiutangEntity
}
type AgingUtangPiutangUpdateResponse struct {
	model.AgingUtangPiutangEntityModel
}
type AgingUtangPiutangUpdateResponseDoc struct {
	Body struct {
		Meta res.Meta                        `json:"meta"`
		Data AgingUtangPiutangUpdateResponse `json:"data"`
	} `json:"body"`
}

// Delete
type AgingUtangPiutangDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type AgingUtangPiutangDeleteResponse struct {
	model.AgingUtangPiutangEntityModel
}
type AgingUtangPiutangDeleteResponseDoc struct {
	Body struct {
		Meta res.Meta                        `json:"meta"`
		Data AgingUtangPiutangDeleteResponse `json:"data"`
	} `json:"body"`
}

// export
type AgingUtangPiutangExportRequest struct {
	UserId    int
	Period    string `query:"period" validate:"required"`
	Versions  int    `query:"version" validate:"required"`
	CompanyID int    `query:"company_id"`
}
type AgingUtangPiutangExportResponse struct {
	File string `json:"file"`
}
type AgingUtangPiutangExportResponseDoc struct {
	Body struct {
		Meta res.Meta                        `json:"meta"`
		Data AgingUtangPiutangExportResponse `json:"data"`
	} `json:"body"`
}

// Import
type AgingUtangPiutangImportRequest struct {
	UserId    int
	CompanyId int
	File      multipart.File
}
type AgingUtangPiutangImportResponse struct {
	Data []model.AgingUtangPiutangEntityModel
}
type AgingUtangPiutangImportResponseDoc struct {
	Body struct {
		Meta res.Meta                           `json:"meta"`
		Data model.AgingUtangPiutangEntityModel `json:"data"`
	} `json:"body"`
}
