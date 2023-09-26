package model

import (
	"mcash-finance-console-core/internal/abstraction"
	"time"
)

type ValidationEntity struct {
	CompanyID      int        `json:"company_id" `
	Period         string     `json:"period" `
	Versions       int        `json:"versions" `
	ValidationNote string     `json:"description" `
	Status         int        `json:"status" `
	CreatedAt      time.Time  `json:"created_at"`
	ModifiedAt     *time.Time `json:"modified_at"`
}

type ValidationFilter struct {
	Period         *string `query:"period"`
	Versions       *int    `query:"versions"`
	Status         *int    `query:"status"`
	TrialBalanceID *int    `query:"id"`
}

type ValidationEntityModel struct {
	// abstraction
	ID int `json:"id"`

	// entity
	ValidationEntity

	ValidationDetail []ValidationDetailEntityModel `json:"validation_detail" gorm:"-"`
	Company          CompanyEntityModel            `json:"company" gorm:"-"`

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type ValidationFilterModel struct {
	// abstraction
	// abstraction.Filter

	// filter
	ValidationFilter
	CompanyCustomFilter
}
