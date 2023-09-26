package model

import (
	"worker-validation/internal/abstraction"
	"worker-validation/pkg/util/date"

	"gorm.io/gorm"
)

type PembelianPenjualanBerelasiEntity struct {
	Period    string `json:"period" validate:"required"`
	Versions  int    `json:"versions" validate:"required"`
	CompanyID int    `json:"company_id" validate:"required"`
	Status    int    `json:"status" validate:"required"`
}

type PembelianPenjualanBerelasiFilter struct {
	Period   *string `query:"period"`
	Versions *int    `query:"versions"`
	Status   *int    `query:"status"`
}

type PembelianPenjualanBerelasiEntityModel struct {
	// abstraction
	abstraction.Entity

	// entity
	PembelianPenjualanBerelasiEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleId"`
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
	Versions []map[int]string `json:"versions"`
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
