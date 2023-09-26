package dto

import (
	"worker/internal/abstraction"
	"worker/internal/model"
)

// Get
type TrialBalanceDetailGetRequest struct {
	abstraction.Pagination
	model.TrialBalanceDetailFilterModel
}
type TrialBalanceDetailGetResponse struct {
	Datas          []model.TrialBalanceDetailEntityModel
	PaginationInfo abstraction.PaginationInfo
}
// type TrialBalanceDetailGetResponseDoc struct {
// 	Body struct {
// 		Meta res.Meta                              `json:"meta"`
// 		Data []model.TrialBalanceDetailEntityModel `json:"data"`
// 	} `json:"body"`
// }

// GetByID
type TrialBalanceDetailGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type TrialBalanceDetailGetByIDResponse struct {
	model.TrialBalanceDetailEntityModel
}
// type TrialBalanceDetailGetByIDResponseDoc struct {
// 	Body struct {
// 		Meta res.Meta                          `json:"meta"`
// 		Data TrialBalanceDetailGetByIDResponse `json:"data"`
// 	} `json:"body"`
// }

// Create
type TrialBalanceDetailCreateRequest struct {
	model.TrialBalanceDetailEntity
}
type TrialBalanceDetailCreateResponse struct {
	model.TrialBalanceDetailEntityModel
}
// type TrialBalanceDetailCreateResponseDoc struct {
// 	Body struct {
// 		Meta res.Meta                         `json:"meta"`
// 		Data TrialBalanceDetailCreateResponse `json:"data"`
// 	} `json:"body"`
// }

// Update
type TrialBalanceDetailUpdateRequest struct {
	ID int `param:"id" validate:"required,numeric"`
	model.TrialBalanceDetailEntity
}
type TrialBalanceDetailUpdateResponse struct {
	model.TrialBalanceDetailEntityModel
}
// type TrialBalanceDetailUpdateResponseDoc struct {
// 	Body struct {
// 		Meta res.Meta                         `json:"meta"`
// 		Data TrialBalanceDetailUpdateResponse `json:"data"`
// 	} `json:"body"`
// }

// Delete
type TrialBalanceDetailDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type TrialBalanceDetailDeleteResponse struct {
	model.TrialBalanceDetailEntityModel
}

// export
type TrialBalanceDetailExportRequest struct {
	model.TrialBalanceDetailFilterModel
}
type TrialBalanceDetailExportResponse struct {
	File string `json:"file"`
}
// type TrialBalanceDetailExportResponseDoc struct {
// 	Body struct {
// 		Meta res.Meta                         `json:"meta"`
// 		Data TrialBalanceDetailExportResponse `json:"data"`
// 	} `json:"body"`
// }

// Import
type TrialBalanceDetailImportRequest struct {
	Datas []model.TrialBalanceDetailEntity
}
type TrialBalanceDetailImportResponse struct {
	Datas          []model.TrialBalanceDetailEntityModel
	PaginationInfo abstraction.PaginationInfo
}
// type TrialBalanceDetailImportResponseDoc struct {
// 	Body struct {
// 		Meta res.Meta                              `json:"meta"`
// 		Data []model.TrialBalanceDetailEntityModel `json:"data"`
// 	} `json:"body"`
// }
