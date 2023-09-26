package dto

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
)

// Get
type ConsolidationGetRequest struct {
	abstraction.Pagination
	model.ConsolidationFilterModel
}
type ConsolidationGetResponse struct {
	Datas          []model.ConsolidationEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type ConsolidationGetResponseDoc struct {
	Body struct {
		Meta res.Meta                 `json:"meta"`
		Data ConsolidationGetResponse `json:"data"`
	} `json:"body"`
}

// GetListCompany
type ConsolidationGetListAvailable struct {
	CompanyID int    `param:"company_id" validate:"required"`
	Period    string `query:"period" validate:"required"`
	abstraction.Pagination
}
type ConsolidationGetListAvaibleResponse struct {
	Parent         []model.TrialBalanceEntityModel
	ChildOnly      []model.TrialBalanceEntityModel
	ChildParent    []model.ConsolidationEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type ConsolidationGetListAvailableResponseDoc struct {
	Body struct {
		Meta res.Meta                            `json:"meta"`
		Data ConsolidationGetListAvaibleResponse `json:"data"`
	} `json:"body"`
}

// conolidation
type ConsolidationConsolidateRequest struct {
	ConsolidationID         int   `json:"consolidation_id" validate:"required"`
	ConsolidationMasterID   int   `json:"tb_id" validate:"required"`
	ListConsolidation       []int `json:"list_to_consolidation" validate:"required"`
	ListConsolidationParent []int `json:"list_to_consolidation_parent"`
}

type ConsolidationCombaineRequest struct {
	ConsolidationMasterID   int   `json:"tb_id" validate:"required"`
	ListConsolidation       []int `json:"list_to_consolidation" validate:"required"`
	ListConsolidationParent []int `json:"list_to_consolidation_parent"`
}
type ConsolidationConsolidateResponse struct {
	model.ConsolidationEntityModel
}
type ConsolidationConsolidateResponseDoc struct {
	Body struct {
		Meta res.Meta                         `json:"meta"`
		Data ConsolidationConsolidateResponse `json:"data"`
	} `json:"body"`
}

// GetListCompanyDuplicate
type ConsolidationGetListDuplicateAvailable struct {
	ConsolidationID int `param:"consolidation_id" validate:"required"`
	abstraction.Pagination
}
type ConsolidationGetListDuplicateAvaibleResponse struct {
	Parent                       []model.TrialBalanceEntityModel
	ChildOnly                    []model.TrialBalanceEntityModel
	ChildParent                  []model.ConsolidationEntityModel
	ConsolidationParent          model.TrialBalanceEntityModel
	ConsolidationChildOnly       []model.TrialBalanceEntityModel
	ConsolidationChildParentOnly []model.ConsolidationEntityModel
	General                      model.ConsolidationEntityModel
	PaginationInfo               abstraction.PaginationInfo
}

// ConsolidationChildOnly
// ConsolidationChildParentOnly
type ConsolidationGetListDuplicateAvailableResponseDoc struct {
	Body struct {
		Meta res.Meta                                     `json:"meta"`
		Data ConsolidationGetListDuplicateAvaibleResponse `json:"data"`
	} `json:"body"`
}

// Get
type FindListCompanyCreateNewCombineGetRequest struct {
	abstraction.Pagination
	model.CompanyFilterModel
	// ChildCompany *bool `query:"include_child_company"`
}
type FindListCompanyCreateNewCombineGetResponse struct {
	Datas          []model.CompanyEntityModel
	PaginationInfo abstraction.PaginationInfo
}
type FindListCompanyCreateNewCombineGetResponseDoc struct {
	Body struct {
		Meta res.Meta                   `json:"meta"`
		Data []model.CompanyEntityModel `json:"data"`
	} `json:"body"`
}

// Delete
type ConsolidationDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type ConsolidationDeleteResponse struct {
	model.ConsolidationEntityModel
}
type ConsolidationDeleteResponseDoc struct {
	Body struct {
		Meta res.Meta                    `json:"meta"`
		Data ConsolidationDeleteResponse `json:"data"`
	} `json:"body"`
}

// GetByID
type ConsolidationGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type ConsolidationGetByIDResponse struct {
	model.ConsolidationEntityModel
}
type ConsolidationGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta                     `json:"meta"`
		Data ConsolidationGetByIDResponse `json:"data"`
	} `json:"body"`
}

// Get Control
type ConsolidationGetControlRequest struct {
	ConsolidationID int `param:"consolidation_id" validate:"required,numeric"`
}
type ConsolidationGetControlResponse struct {
	Datas model.ConsolidationDetailEntityModel
}
type ConsolidationGetControlResponseDoc struct {
	Body struct {
		Meta res.Meta                        `json:"meta"`
		Data ConsolidationGetControlResponse `json:"data"`
	} `json:"body"`
}
