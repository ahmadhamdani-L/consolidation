package model

import (
	"worker/internal/abstraction"
)

type NotificationEntity struct {
	Description string `json:"description" validate:"required"`
	IsOpened    *bool  `json:"is_open" validate:"required"`
	Data        string `json:"data" `
}

type NotificationFilter struct {
	IsOpened *bool   `query:"is_open" validate:"required"`
	Data     *string `query:"-"  filter:"CUSTOM"`
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
