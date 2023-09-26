package model

import (
	"mcash-finance-console-core/internal/abstraction"
)

type AccessScopeDetailEntity struct {
	AccessScopeID int     `json:"access_scope_id" validate:"required"`
	CompanyID     int     `json:"company_id" validate:"required"`
	CompanyString *string `json:"company" gorm:"-"`
}

type AccessScopeDetailFilter struct {
	AccessScopeID *int    `query:"access_scope_id"`
	CompanyID     *int    `query:"company_id"`
	CompanyString *string `query:"company" filter:"CUSTOM" example:"PT ABC"`
}

type AccessScopeDetailEntityModel struct {
	ID int `json:"id" gorm:"primaryKey;autoIncrement;"`

	// entity
	AccessScopeDetailEntity

	// relations
	AccessScope AccessScopeEntityModel `json:"access_scope" gorm:"foreignKey:AccessScopeID"`
	Company     CompanyEntityModel     `json:"company" gorm:"foreignKey:CompanyID"`

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type AccessScopeDetailListEntityModel struct {
	abstraction.Entity
	CompanyEntity
	Child    []AccessScopeDetailListEntityModel `json:"child" gorm:"-"`
	IsParent bool                               `json:"is_parent" gorm:"-"`
	Checked  bool                               `json:"checked"`
}

type AccessScopeDetailFilterModel struct {

	// filter
	AccessScopeDetailFilter
}

func (AccessScopeDetailEntityModel) TableName() string {
	return "access_scope_detail"
}
