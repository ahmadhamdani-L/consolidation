package dto

import (
	"worker/internal/abstraction"
	"worker/internal/model"
)

// Get
type TrialBalanceGetRequest struct {
	abstraction.Pagination
	model.TrialBalanceFilterModel
}
type TrialBalanceGetResponse struct {
	Datas          []model.TrialBalanceEntityModel
	PaginationInfo abstraction.PaginationInfo
}
// type TrialBalanceGetResponseDoc struct {
// 	Body struct {
// 		Meta res.Meta                        `json:"meta"`
// 		Data []model.TrialBalanceEntityModel `json:"data"`
// 	} `json:"body"`
// }

// GetByID
type TrialBalanceGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type TrialBalanceGetByIDResponse struct {
	model.TrialBalanceEntityModel
}
// type TrialBalanceGetByIDResponseDoc struct {
// 	Body struct {
// 		Meta res.Meta                    `json:"meta"`
// 		Data TrialBalanceGetByIDResponse `json:"data"`
// 	} `json:"body"`
// }

// Create
type TrialBalanceCreateRequest struct {
	model.TrialBalanceEntity
}
type TrialBalanceCreateResponse struct {
	model.TrialBalanceEntityModel
}
// type TrialBalanceCreateResponseDoc struct {
// 	Body struct {
// 		Meta res.Meta                   `json:"meta"`
// 		Data TrialBalanceCreateResponse `json:"data"`
// 	} `json:"body"`
// }

// Update
type TrialBalanceUpdateRequest struct {
	ID int `param:"id" validate:"required,numeric"`
	model.TrialBalanceEntity
}
type TrialBalanceUpdateResponse struct {
	model.TrialBalanceEntityModel
}
// type TrialBalanceUpdateResponseDoc struct {
// 	Body struct {
// 		Meta res.Meta                   `json:"meta"`
// 		Data TrialBalanceUpdateResponse `json:"data"`
// 	} `json:"body"`
// }

// Delete
type TrialBalanceDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type TrialBalanceDeleteResponse struct {
	model.TrialBalanceEntityModel
}
// type TrialBalanceDeleteResponseDoc struct {
// 	Body struct {
// 		Meta res.Meta                   `json:"meta"`
// 		Data TrialBalanceDeleteResponse `json:"data"`
// 	} `json:"body"`
// }

// export
type TrialBalanceExportRequest struct {
	UserId    int
	Period    string `query:"period" validate:"required"`
	Versions  int    `query:"version" validate:"required"`
	CompanyID int    `query:"company_id"`
}
type TrialBalanceExportResponse struct {
	File string `json:"file"`
}
// type TrialBalanceExportResponseDoc struct {
// 	Body struct {
// 		Meta res.Meta                   `json:"meta"`
// 		Data TrialBalanceExportResponse `json:"data"`
// 	} `json:"body"`
// }

// Import
type TrialBalanceImportRequest struct {
	UserId    int
	CompanyId int
}
type TrialBalanceImportResponse struct {
	Data model.TrialBalanceEntityModel
}
// type TrialBalanceImportResponseDoc struct {
// 	Body struct {
// 		Meta res.Meta                        `json:"meta"`
// 		Data []model.TrialBalanceEntityModel `json:"data"`
// 	} `json:"body"`
// }
