package model

import (
	"mcash-finance-console-core/internal/abstraction"
)

type AccessScopeEntity struct {
	UserID       int    `json:"user_id" validate:"required"`
	UserString   string `json:"user" gorm:"-"`
	AccessAll    *bool  `json:"access_all" validate:"required"`
	CompanyList  string `json:"company_list" gorm:"-"`
	JmlCompany   int    `json:"total_company,omitempty" gorm:"-"`
	UserIsActive *bool  `json:"status,omitempty" gorm:"-"`
}

type AccessScopeFilter struct {
	UserID       *int    `query:"user_id" example:"1"`
	UserString   *string `query:"user" filter:"CUSTOM" example:"Lutfi Ramadhan"`
	AccessAll    *bool   `query:"access_all" example:"1"`
	UserIsActive *bool   `query:"status" filter:"CUSTOM"`
}

type AccessScopeEntityModel struct {
	// abstraction
	ID int `json:"id" gorm:"primaryKey;autoIncrement;"`

	// entity
	AccessScopeEntity

	// relations
	User              UserEntityModel                `json:"-" gorm:"foreignKey:UserID"`
	AccessScopeDetail []AccessScopeDetailEntityModel `json:"access_scope_detail" gorm:"foreignKey:AccessScopeID"`
	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type AccessScopeFilterModel struct {

	// filter
	AccessScopeFilter

	CompanyCustomFilter
}

func (AccessScopeEntityModel) TableName() string {
	return "access_scope"
}

// func (m *AccessScopeEntityModel) BeforeCreate(tx *gorm.DB) (err error) {
// 	m.CreatedAt = *date.DateTodayLocal()
// 	m.CreatedBy = m.Context.Auth.ID
// 	return
// }

// func (m *AccessScopeEntityModel) BeforeUpdate(tx *gorm.DB) (err error) {
// 	m.ModifiedAt = date.DateTodayLocal()
// 	m.ModifiedBy = &m.Context.Auth.ID
// 	return
// }
