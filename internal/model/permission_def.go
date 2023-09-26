package model

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/pkg/util/date"

	"gorm.io/gorm"
)

type PermissionDefEntity struct {
	FunctionalID  string `json:"functional_id" validate:"required" gorm:"column:functional_id;index:permissions_def_functional_id_idx,unique"`
	Label         string `json:"label" validate:"required"`
	AllowCreate   *bool  `json:"allow_create" validate:"required"`
	ApiPathCreate string `json:"api_path-" validate:"required"`
	AllowRead     *bool  `json:"allow_read" validate:"required"`
	ApiPathRead   string `json:"-" validate:"required"`
	AllowUpdate   *bool  `json:"allow_update" validate:"required"`
	ApiPathUpdate string `json:"-" validate:"required"`
	AllowDelete   *bool  `json:"allow_delete" validate:"required"`
	ApiPathDelete string `json:"-" validate:"required"`
}

type PermissionDefFilter struct {
	FunctionalID  *string `query:"functional_id" filter:"ILIKE" example:"TRIAL-BALANCE"`
	Label         *string `query:"label" filter:"ILIKE" example:"TRIAL BALANCE"`
	AllowCreate   *bool   `query:"allow_create" example:"true"`
	ApiPathCreate *string `query:"-" filter:"ILIKE" example:"/trial-balance/create"`
	AllowRead     *bool   `query:"allow_read" example:"true"`
	ApiPathRead   *string `query:"-" filter:"ILIKE" example:"/trial-balance/view"`
	AllowUpdate   *bool   `query:"allow_update" example:"true"`
	ApiPathUpdate *string `query:"-" filter:"ILIKE" example:"/trial-balance/create"`
	AllowDelete   *bool   `query:"allow_delete" example:"false"`
	ApiPathDelete *string `query:"-" filter:"ILIKE" example:"/trial-balance/delete"`
}

type PermissionDefEntityModel struct {
	// abstraction
	abstraction.Entity

	// entity
	PermissionDefEntity

	// relations
	RolePermission []RolePermissionEntityModel `json:"role_permissions" gorm:"foreignKey:FunctionalID;references:FunctionalID"`

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`

	UserRelationModel
}

type PermissionDefFilterModel struct {
	abstraction.Filter
	// filter
	PermissionDefFilter
}

func (PermissionDefEntityModel) TableName() string {
	return "permissions_def"
}

func (m *PermissionDefEntityModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = *date.DateTodayLocal()
	m.CreatedBy = m.Context.Auth.ID
	return
}

func (m *PermissionDefEntityModel) BeforeUpdate(tx *gorm.DB) (err error) {
	m.ModifiedAt = date.DateTodayLocal()
	m.ModifiedBy = &m.Context.Auth.ID
	return
}
