package model

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/pkg/util/date"

	"gorm.io/gorm"
)

type RoleEntity struct {
	Code string `json:"code"`
	Name string `json:"name" validate:"required"`
}

type RoleFilter struct {
	Search *string `query:"search" filter:"CUSTOM" example:"1"`
}

type RoleEntityModel struct {
	// abstraction
	abstraction.Entity

	// entity
	RoleEntity

	// relations

	RolePermissionApi []RolePermissionApiEntityModel `json:"role_permission_api" gorm:"foreignKey:RoleID"`
	RolePermission    []RolePermissionEntityModel    `json:"role_permission" gorm:"foreignKey:RoleID"`
	User              []UserEntityModel              `json:"user" gorm:"foreignKey:RoleID"`
	// context
	Context *abstraction.Context `json:"-" gorm:"-"`

	UserRelationModel
}

type RoleFilterModel struct {
	abstraction.Filter
	// filter
	RoleFilter
}

func (RoleEntityModel) TableName() string {
	return "roles"
}

func (m *RoleEntityModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = *date.DateTodayLocal()
	m.CreatedBy = m.Context.Auth.ID
	return
}

func (m *RoleEntityModel) BeforeUpdate(tx *gorm.DB) (err error) {
	m.ModifiedAt = date.DateTodayLocal()
	m.ModifiedBy = &m.Context.Auth.ID
	return
}
