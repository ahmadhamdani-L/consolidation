package model

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/pkg/util/date"

	"gorm.io/gorm"
)

type NotificationEntity struct {
	Description string `json:"description" validate:"required"`
	IsOpened    *bool  `json:"is_opened" validate:"required"`
	Data        string `json:"data" validate:"required"`
}

type NotificationFilter struct {
	IsOpened *bool   `query:"is_opened"`
	Data     *string `query:"-" filter:"CUSTOM"`
}

type NotificationEntityModel struct {
	// abstraction
	abstraction.Entity

	// entity
	NotificationEntity
	UserRelationModel

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type NotificationFilterModel struct {
	// abstraction
	abstraction.Filter

	// filter
	NotificationFilter
}

func (NotificationEntityModel) TableName() string {
	return "notification"
}

func (m *NotificationEntityModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = *date.DateTodayLocal()
	m.CreatedBy = m.Context.Auth.ID
	return
}

func (m *NotificationEntityModel) BeforeUpdate(tx *gorm.DB) (err error) {
	m.ModifiedAt = date.DateTodayLocal()
	m.ModifiedBy = &m.Context.Auth.ID
	return
}
