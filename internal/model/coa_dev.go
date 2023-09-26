package model

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/pkg/util/date"

	"gorm.io/gorm"
)

type CoaDevEntity struct {
	Code      string `json:"code" validate:"required" example:"110101003"`
	Name      string `json:"name" validate:"required" example:"CA - Kas Di Tangan - IDR - Kas Akun Digital"`
	CoaTypeID int    `json:"coa_type_id" validate:"required" example:"1"`
}

type CoaDevFilter struct {
	Code          *string `query:"code" filter:"LIKE" example:"110101003"`
	Name          *string `query:"name" filter:"ILIKE" example:"CA - Kas Di Tangan - IDR - Kas Akun Digital"`
	CoaTypeID     *int    `query:"coa_type_id" example:"1"`
	ArrCoaGroupID *[]int  `example:"1" filter:"CUSTOM"`
	Search        *string `query:"s" filter:"CUSTOM" example:"CA / 110101003"`
}

type CoaDevEntityModel struct {
	// abstraction
	abstraction.Entity

	// entity
	CoaDevEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleID"`
	CoaType CoaTypeEntityModel `json:"coa_type" gorm:"foreignKey:CoaTypeID"`
	UserRelationModel

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type CoaDevFilterModel struct {
	// abstraction
	abstraction.Filter

	// filter
	CoaFilter
}

func (CoaDevEntityModel) TableName() string {
	return "m_coa_dev"
}

func (m *CoaDevEntityModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = *date.DateTodayLocal()
	m.CreatedBy = m.Context.Auth.ID
	return
}

func (m *CoaDevEntityModel) BeforeUpdate(tx *gorm.DB) (err error) {
	m.ModifiedAt = date.DateTodayLocal()
	m.ModifiedBy = &m.Context.Auth.ID
	return
}
