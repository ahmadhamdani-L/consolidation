package model

import (
	"worker-consol/internal/abstraction"
	"worker-consol/pkg/util/date"

	"gorm.io/gorm"
)

type ConsolidationEntity struct {
	Period                string `json:"period" validate:"required"`
	Versions              int    `json:"versions" validate:"required"`
	CompanyID             int    `json:"company_id" validate:"required"`
	ConsolidationVersions int    `json:"consolidation_versions" validate:"required"`
	Status                int    `json:"status"`
	IsDuplicated          *bool  `json:"is_duplicated"`
}

type ConsolidationFilter struct {
	Period                *string `query:"period"`
	Versions              *int    `query:"versions"`
	ConsolidationVersions *int    `query:"consolidation_versions"`
	// FormatterID *int    `query:"formatter_id"`
	Status *int `query:"status"`
}

type ConsolidationEntityModel struct {
	// abstraction
	abstraction.Entity

	// entity
	ConsolidationEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleId"`
	// Formatter               FormatterEntityModel                 `json:"formatter" gorm:"foreignKey:FormatterID"`
	Company             CompanyEntityModel               `json:"company" gorm:"foreignKey:CompanyID"`
	ConsolidationBridge []ConsolidationBridgeEntityModel `json:"consolidation_bridge" gorm:"foreignKey:ConsolidationID"`
	ConsolidationDetail []ConsolidationDetailEntityModel `json:"consolidation_detail" gorm:"foreignKey:ConsolidationID"`

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type ConsolidationFilterModel struct {
	// abstraction
	abstraction.Filter

	// filter
	ConsolidationFilter
	CompanyCustomFilter
}

func (ConsolidationEntityModel) TableName() string {
	return "consolidation"
}

func (m *ConsolidationEntityModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = *date.DateTodayLocal()
	return
}

func (m *ConsolidationEntityModel) BeforeUpdate(tx *gorm.DB) (err error) {
	m.ModifiedAt = date.DateTodayLocal()
	return
}

type ConsolidationFullDetail struct {
	Code            string   `json:"code" gorm:"column:code"`
	AmountBeforeJpm *float64 `json:"amount_before_jpm" gorm:"amount_before_jpm"`
	AmountJpmDr     *float64 `json:"amount_jpm_dr" gorm:"amount_jpm_dr"`
	AmountJpmCr     *float64 `json:"amount_jpm_cr" gorm:"amount_jpm_cr"`
	AmountAfterJpm  *float64 `json:"amount_after_jpm" gorm:"amount_after_jpm"`
	AmountJcteDr    *float64 `json:"amount_jcte_dr" gorm:"amount_jcte_dr"`
	AmountJcteCr    *float64 `json:"amount_jcte_cr" gorm:"amount_jcte_cr"`
	AmountJelimDr   *float64 `json:"amount_jelim_dr" gorm:"amount_jelim_dr"`
	AmountJelimCr   *float64 `json:"amount_jelim_cr" gorm:"amount_jelim_cr"`
	Amount          *float64 `json:"amount" gorm:"amount"`
	// AmountAfterJcte         *float64 `json:"amount_after_jcte" gorm:"amount_after_jcte"`
	// AmountCombineSubsidiary *float64 `json:"amount_combine_subsidiary" gorm:"amount_combine_subsidiary"`
	// AmountConsole           *float64 `json:"amount_console" gorm:"amount_console"`
}
