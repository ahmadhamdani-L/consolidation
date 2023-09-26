package model

import (
	"mcash-finance-console-core/internal/abstraction"
)

type PembelianPenjualanBerelasiDetailEntity struct {
	PembelianPenjualanBerelasiID int      `json:"pembelian_penjualan_berelasi_id" validate:"required" example:"1"`
	Code                         string   `json:"code" validate:"required" example:"001"`
	Name                         string   `json:"description" validate:"required" example:"PT ABC"`
	BoughtAmount                 *float64 `json:"bought_amount" validate:"required" example:"10000.00"`
	SalesAmount                  *float64 `json:"sales_amount" validate:"required" example:"10000.00"`
	SortID                       int      `json:"sort_id" validate:"required" example:"1"`
}

type PembelianPenjualanBerelasiDetailFilter struct {
	PembelianPenjualanBerelasiID *int     `query:"pembelian_penjualan_berelasi_id" validate:"required" example:"1"`
	Code                         *string  `query:"code" filter:"ILIKE" example:"001"`
	Name                         *string  `query:"description" filter:"ILIKE" example:"PT ABC"`
	SalesAmount                  *float64 `query:"bought_amount" example:"10000.00"`
	BoughtAmount                 *float64 `query:"sales_amount" example:"10000.00"`
	SortID                       *int     `query:"sort_id" example:"1"`
	Period      *string `query:"period" example:"2022-01-01" filter:"DATESTRING"`
	Versions    *int    `query:"versions" example:"1"`
	ArrVersions *[]int  `filter:"CUSTOM" example:"1"`
	// FormatterID *int    `query:"formatter_id" example:"1"`
	Status *int `query:"status" example:"1"`
}

type PembelianPenjualanBerelasiDetailEntityModel struct {
	ID int `json:"id" gorm:"primaryKey;autoIncrement;"`

	// entity
	PembelianPenjualanBerelasiDetailEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleID"`
	PembelianPenjualanBerelasi PembelianPenjualanBerelasiEntityModel `json:"-" gorm:"foreignKey:PembelianPenjualanBerelasiID"`
	Company                    CompanyEntityModel                    `json:"company" gorm:"foreignKey:Code;references:Code"`

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type PembelianPenjualanBerelasiDetailFilterModel struct {

	// filter
	PembelianPenjualanBerelasiDetailFilter

	CompanyCustomFilter
}

func (PembelianPenjualanBerelasiDetailEntityModel) TableName() string {
	return "pembelian_penjualan_berelasi_detail"
}
