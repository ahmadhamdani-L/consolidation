package model

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/pkg/util/date"

	"gorm.io/gorm"
)

type RolePermissionEntity struct {
	FunctionalID string `json:"functional_id" validate:"required" gorm:"column:functional_id"`
	RoleID       int    `json:"role_id" validate:"required"`
	Create       *bool  `json:"create" validate:"required" gorm:"default:false"`
	Read         *bool  `json:"read" validate:"required" gorm:"default:false"`
	Update       *bool  `json:"update" validate:"required" gorm:"default:false"`
	Delete       *bool  `json:"delete" validate:"required" gorm:"default:false"`
}

type RolePermissionFilter struct {
	FunctionalID *string `query:"functional_id" filter:"ILIKE" example:"TRIAL-BALANCE"`
	RoleID       *int    `query:"role_id" example:"1"`
	Create       *bool   `query:"create" example:"true"`
	Read         *bool   `query:"read" example:"true"`
	Update       *bool   `query:"update" example:"true"`
	Delete       *bool   `query:"delete" example:"true"`
}

type RolePermissionEntityModel struct {
	// abstraction
	abstraction.Entity

	// entity
	RolePermissionEntity

	// relations
	Role          RoleEntityModel          `json:"role" gorm:"foreignKey:RoleID"`
	PermissionDef PermissionDefEntityModel `json:"permissions_def" gorm:"foreignKey:FunctionalID;references:FunctionalID"`

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`

	UserRelationModel
}

type RolePermissionFilterModel struct {
	abstraction.Filter
	// filter
	RolePermissionFilter
}

func (RolePermissionEntityModel) TableName() string {
	return "role_permissions"
}

func (m *RolePermissionEntityModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = *date.DateTodayLocal()
	m.CreatedBy = m.Context.Auth.ID
	return
}

func (m *RolePermissionEntityModel) BeforeUpdate(tx *gorm.DB) (err error) {
	m.ModifiedAt = date.DateTodayLocal()
	m.ModifiedBy = &m.Context.Auth.ID
	return
}
