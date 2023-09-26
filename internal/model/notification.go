package model

import (
	"worker-validation/internal/abstraction"
)

type NotificationEntity struct {
	Description string `json:"description" validate:"required"`
	IsOpened    *bool  `json:"is_open" validate:"required"`
	Data        string `json:"data" validate:"required"`
}

type NotificationFilter struct {
	IsOpened *bool   `query:"is_open" validate:"required"`
	Data     *string `query:"-" validate:"required" filter:"CUSTOM"`
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
