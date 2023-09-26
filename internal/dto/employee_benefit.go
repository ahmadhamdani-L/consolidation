package dto

import (
	"worker/internal/abstraction"
	"worker/internal/model"
	res "worker/pkg/util/response"
	"mime/multipart"
)

// Get
type EmployeeBenefitGetRequest struct {
	abstraction.Pagination
	model.EmployeeBenefitFilterModel
}
type EmployeeBenefitGetResponse struct {
	Datas          []model.EmployeeBenefitEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type EmployeeBenefitGetResponseDoc struct {
	Body struct {
		Meta res.Meta                            `json:"meta"`
		Data []model.EmployeeBenefitEntityModel `json:"data"`
	} `json:"body"`
}

// GetByID
type EmployeeBenefitGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type EmployeeBenefitGetByIDResponse struct {
	model.EmployeeBenefitEntityModel
}
type EmployeeBenefitGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta                        `json:"meta"`
		Data EmployeeBenefitGetByIDResponse `json:"data"`
	} `json:"body"`
}

// Create
type EmployeeBenefitCreateRequest struct {
	model.EmployeeBenefitEntity
}
type EmployeeBenefitCreateResponse struct {
	model.EmployeeBenefitEntityModel
}
type EmployeeBenefitCreateResponseDoc struct {
	Body struct {
		Meta res.Meta                       `json:"meta"`
		Data EmployeeBenefitCreateResponse `json:"data"`
	} `json:"body"`
}

// Update
type EmployeeBenefitUpdateRequest struct {
	ID int `param:"id" validate:"required,numeric"`
	model.EmployeeBenefitEntity
}
type EmployeeBenefitUpdateResponse struct {
	model.EmployeeBenefitEntityModel
}
type EmployeeBenefitUpdateResponseDoc struct {
	Body struct {
		Meta res.Meta                       `json:"meta"`
		Data EmployeeBenefitUpdateResponse `json:"data"`
	} `json:"body"`
}

// Delete
type EmployeeBenefitDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type EmployeeBenefitDeleteResponse struct {
	model.EmployeeBenefitEntityModel
}
type EmployeeBenefitDeleteResponseDoc struct {
	Body struct {
		Meta res.Meta                       `json:"meta"`
		Data EmployeeBenefitDeleteResponse `json:"data"`
	} `json:"body"`
}

// export
type EmployeeBenefitExportRequest struct {
	UserId    int
	Period    string `query:"period" validate:"required"`
	Versions  int    `query:"version" validate:"required"`
	CompanyID int    `query:"company_id"`
}
type EmployeeBenefitExportResponse struct {
	File string `json:"file"`
}
type EmployeeBenefitExportResponseDoc struct {
	Body struct {
		Meta res.Meta                       `json:"meta"`
		Data EmployeeBenefitExportResponse `json:"data"`
	} `json:"body"`
}

// Import
type EmployeeBenefitImportRequest struct {
	File      multipart.File
	UserId    int
	CompanyId int
}
type EmployeeBenefitImportResponse struct {
	Data           []model.EmployeeBenefitEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type EmployeeBenefitImportResponseDoc struct {
	Body struct {
		Meta res.Meta                          `json:"meta"`
		Data model.EmployeeBenefitEntityModel `json:"data"`
	} `json:"body"`
}
