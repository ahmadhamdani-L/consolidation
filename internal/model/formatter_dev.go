package model

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/pkg/util/date"

	"gorm.io/gorm"
)

type FormatterDevEntity struct {
	FormatterFor string `json:"formatter_for" validate:"required" example:"TRIAL-BALANCE"`
	Description  string `json:"description" validate:"required" example:"Template Trial balance"`
}

type FormatterDevFilter struct {
	FormatterFor *string `query:"formatter_for" filter:"ILIKE" example:"TRIAL-BALANCE"`
	Description  *string `query:"description" filter:"ILIKE" example:"Template Trial balance"`
}

type FormatterDevEntityModel struct {
	// abstraction
	abstraction.Entity

	// entity
	FormatterDevEntity

	// relations
	FormatterDetail  []FormatterDetailEntityModel  `json:"formatter_detail" gorm:"ForeignKey:FormatterID"`
	FormatterBridges []FormatterBridgesEntityModel `json:"-" gorm:"foreignKey:FormatterID"`
	UserRelationModel

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type FormatterDevFilterModel struct {
	// abstraction
	abstraction.Filter

	// filter
	FormatterFilter
}

func (FormatterDevEntityModel) TableName() string {
	return "m_formatter_dev"
}

func (m *FormatterDevEntityModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = *date.DateTodayLocal()
	m.CreatedBy = m.Context.Auth.ID
	return
}

func (m *FormatterDevEntityModel) BeforeUpdate(tx *gorm.DB) (err error) {
	m.ModifiedAt = date.DateTodayLocal()
	m.ModifiedBy = &m.Context.Auth.ID
	return
}
