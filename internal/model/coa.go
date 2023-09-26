package model

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/pkg/util/date"

	"gorm.io/gorm"
)

type CoaEntity struct {
	Code      string `json:"code" validate:"required" example:"110101003"`
	Name      string `json:"name" validate:"required" example:"CA - Kas Di Tangan - IDR - Kas Akun Digital"`
	CoaTypeID int    `json:"coa_type_id" validate:"required" example:"1"`
}

type CoaFilter struct {
	Code          *string `query:"code" filter:"LIKE" example:"110101003"`
	Name          *string `query:"name" filter:"ILIKE" example:"CA - Kas Di Tangan - IDR - Kas Akun Digital"`
	CoaTypeID     *int    `query:"coa_type_id" example:"1"`
	ArrCoaGroupID *[]int  `example:"1" filter:"CUSTOM"`
	Search        *string `query:"s" filter:"CUSTOM" example:"CA / 110101003"`
}

type CoaEntityModel struct {
	// abstraction
	abstraction.Entity

	// entity
	CoaEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleID"`
	CoaType CoaTypeEntityModel `json:"coa_type" gorm:"foreignKey:CoaTypeID"`
	UserRelationModel

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type CoaFilterModel struct {
	// abstraction
	abstraction.Filter

	// filter
	CoaFilter
}

func (CoaEntityModel) TableName() string {
	return "m_coa"
}

func (m *CoaEntityModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = *date.DateTodayLocal()
	m.CreatedBy = m.Context.Auth.ID
	return
}

func (m *CoaEntityModel) BeforeUpdate(tx *gorm.DB) (err error) {
	m.ModifiedAt = date.DateTodayLocal()
	m.ModifiedBy = &m.Context.Auth.ID
	return
}
