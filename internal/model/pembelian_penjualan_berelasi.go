package model

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/pkg/util/date"

	"gorm.io/gorm"
)

type PembelianPenjualanBerelasiEntity struct {
	Period    string `json:"period" validate:"required" example:"2022-01-01"`
	Versions  int    `json:"versions" validate:"required" example:"1"`
	CompanyID int    `json:"company_id" validate:"required" example:"1"`
	Status    int    `json:"status,omitempty" example:"1"`
}

type PembelianPenjualanBerelasiFilter struct {
	Period      *string `query:"period" example:"2022-01-01" filter:"DATESTRING"`
	Versions    *int    `query:"versions" example:"1"`
	ArrVersions *[]int  `filter:"CUSTOM" example:"1"`
	Status      *int    `query:"status" example:"1"`
}

type PembelianPenjualanBerelasiEntityModel struct {
	// abstraction
	abstraction.Entity

	// entity
	PembelianPenjualanBerelasiEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleID"`
	Company                          CompanyEntityModel                            `json:"company" gorm:"foreignKey:CompanyID"`
	PembelianPenjualanBerelasiDetail []PembelianPenjualanBerelasiDetailEntityModel `json:"pembelian_penjualan_berelasi_detail" gorm:"foreignKey:PembelianPenjualanBerelasiID"`
	UserRelationModel

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type PembelianPenjualanBerelasiFilterModel struct {
	// abstraction
	abstraction.Filter

	// filter
	PembelianPenjualanBerelasiFilter
	CompanyCustomFilter
}

type PembelianPenjualanBerelasiVersionModel struct {
	Version []map[int]string `json:"versions"`
}

func (PembelianPenjualanBerelasiEntityModel) TableName() string {
	return "pembelian_penjualan_berelasi"
}

func (m *PembelianPenjualanBerelasiEntityModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = *date.DateTodayLocal()
	m.CreatedBy = m.Context.Auth.ID
	return
}

func (m *PembelianPenjualanBerelasiEntityModel) BeforeUpdate(tx *gorm.DB) (err error) {
	m.ModifiedAt = date.DateTodayLocal()
	m.ModifiedBy = &m.Context.Auth.ID
	return
}
