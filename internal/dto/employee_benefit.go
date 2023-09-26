package dto

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
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
		Meta res.Meta                           `json:"meta"`
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
		Meta res.Meta                       `json:"meta"`
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
		Meta res.Meta                      `json:"meta"`
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
		Meta res.Meta                      `json:"meta"`
		Data EmployeeBenefitUpdateResponse `json:"data"`
	} `json:"body"`
}

// Delete
type EmployeeBenefitDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type EmployeeBenefitDeleteResponse struct {
	// model.EmployeeBenefitEntityModel
}
type EmployeeBenefitDeleteResponseDoc struct {
	Body struct {
		Meta res.Meta                      `json:"meta"`
		Data EmployeeBenefitDeleteResponse `json:"data"`
	} `json:"body"`
}

// export
type EmployeeBenefitExportRequest struct {
	EmployeeBenefitID int `query:"employee_benefit_id" validate:"required"`
}

type EmployeeBenefitExportAsyncRequest struct {
	UserID    int
	Period    string `query:"period" validate:"required"`
	Versions  int    `query:"versions" validate:"required"`
	CompanyID int    `query:"company_id"`
}
type EmployeeBenefitExportResponse struct {
	FileName string `json:"filename"`
	Path     string `json:"path"`
}
type EmployeeBenefitExportResponseDoc struct {
	Body struct {
		Meta res.Meta                      `json:"meta"`
		Data EmployeeBenefitExportResponse `json:"data"`
	} `json:"body"`
}

// Import
type EmployeeBenefitImportRequest struct {
	UserID    int
	CompanyID int
	File      multipart.File
}
type EmployeeBenefitImportResponse struct {
	Data model.EmployeeBenefitEntityModel
}
type EmployeeBenefitImportResponseDoc struct {
	Body struct {
		Meta res.Meta                      `json:"meta"`
		Data EmployeeBenefitImportResponse `json:"data"`
	} `json:"body"`
}
