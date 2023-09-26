package model

import (
	"mcash-finance-console-core/internal/abstraction"
)

type RolePermissionApiEntity struct {
	RoleID    int    `json:"role_id" validate:"required"`
	ApiPath   string `json:"api_path" validate:"required"`
	ApiMethod string `json:"api_method" validate:"required"`
}

type RolePermissionApiFilter struct {
	RoleID    *int    `query:"role_id" example:"1"`
	ApiPath   *string `query:"api_path" filter:"ILIKE" example:"/endpoint/action"`
	ApiMethod *string `query:"api_method" filter:"ILIKE" example:"GET"`
}

type RolePermissionApiEntityModel struct {
	// abstraction
	ID int `json:"id" gorm:"primaryKey;autoIncrement;"`

	// entity
	RolePermissionApiEntity

	// relations
	Role RoleEntityModel `json:"role" gorm:"foreignKey:RoleID"`

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type RolePermissionApiFilterModel struct {

	// filter
	RolePermissionApiFilter
}

func (RolePermissionApiEntityModel) TableName() string {
	return "role_permissions_api"
}

// func (m *RolePermissionApiEntityModel) BeforeCreate(tx *gorm.DB) (err error) {
// 	m.CreatedAt = *date.DateTodayLocal()
// 	m.CreatedBy = m.Context.Auth.ID
// 	return
// }

// func (m *RolePermissionApiEntityModel) BeforeUpdate(tx *gorm.DB) (err error) {
// 	m.ModifiedAt = date.DateTodayLocal()
// 	m.ModifiedBy = &m.Context.Auth.ID
// 	return
// }
