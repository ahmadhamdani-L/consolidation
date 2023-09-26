package model

import (
	"worker-validation/internal/abstraction"
)

type PembelianPenjualanBerelasiDetailEntity struct {
	PembelianPenjualanBerelasiID int      `json:"pembelian_penjualan_berelasi_id" validate:"required"`
	Code                         string   `json:"code" validate:"required"`
	Name                         string   `json:"description" validate:"required"`
	BoughtAmount                 *float64 `json:"bought_amount" validate:"required"`
	SalesAmount                  *float64 `json:"sales_amount" validate:"required"`
	SortID                       int      `json:"sort_id" validate:"required"`
}

type PembelianPenjualanBerelasiDetailFilter struct {
	PembelianPenjualanBerelasiID *int     `query:"pembelian_penjualan_berelasi_id" validate:"required"`
	Code                         *string  `query:"code"`
	Name                         *string  `query:"description" filter:"ILIKE"`
	SalesAmount                  *float64 `query:"bought_amount"`
	BoughtAmount                 *float64 `query:"sales_amount"`
	SortID                       *int     `query:"sort_id"`
}

type PembelianPenjualanBerelasiDetailEntityModel struct {
	ID int `json:"id" gorm:"primaryKey;autoIncrement;"`

	// entity
	PembelianPenjualanBerelasiDetailEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleId"`
	PembelianPenjualanBerelasi PembelianPenjualanBerelasiEntityModel `json:"-" gorm:"foreignKey:PembelianPenjualanBerelasiID"`
	Company                    CompanyEntityModel                    `json:"company" gorm:"foreignKey:Code;references:Code"`

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type PembelianPenjualanBerelasiDetailFilterModel struct {

	// filter
	PembelianPenjualanBerelasiDetailFilter
}

func (PembelianPenjualanBerelasiDetailEntityModel) TableName() string {
	return "pembelian_penjualan_berelasi_detail"
}
