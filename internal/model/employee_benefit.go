package model

import (
	"worker-validation/internal/abstraction"
)

type EmployeeBenefitEntity struct {
	Period    string `json:"period" validate:"required" example:"2022-01-01"`
	Versions  int    `json:"versions" validate:"required" example:"1"`
	CompanyID int    `json:"company_id" validate:"required" example:"1"`
	// FormatterID int    `json:"formatter_id" validate:"required" example:"1"`
	Status int `json:"status" validate:"required" example:"1"`
}

type EmployeeBenefitFilter struct {
	Period      *string `query:"period" example:"2022-01-01" filter:"DATESTRING"`
	Versions    *int    `query:"versions" example:"1"`
	ArrVersions *[]int  `filter:"CUSTOM" example:"1"`
	// FormatterID *int    `query:"formatter_id" example:"1"`
	Status *int `query:"status" example:"1"`
}

type EmployeeBenefitEntityModel struct {
	// abstraction
	abstraction.Entity

	// entity
	EmployeeBenefitEntity

	// relations
	// SampleChilds []SampleChildEntityModel `json:"sample_childs" gorm:"foreignKey:SampleID"`
	Company CompanyEntityModel `json:"company" gorm:"foreignKey:CompanyID"`
	// Formatter               FormatterEntityModel                 `json:"formatter" gorm:"foreignKey:FormatterID"`
	EmployeeBenefitDetailAsumsi         []EmployeeBenefitDetailFmtEntityModel `json:"employee_benefit_detail_asumsi" gorm:"-"`
	EmployeeBenefitDetailRekonsiliasi   []EmployeeBenefitDetailFmtEntityModel `json:"employee_benefit_detail_rekonsiliasi" gorm:"-"`
	EmployeeBenefitDetailRincianLaporan []EmployeeBenefitDetailFmtEntityModel `json:"employee_benefit_detail_rincian_laporan" gorm:"-"`
	EmployeeBenefitDetailRincianEkuitas []EmployeeBenefitDetailFmtEntityModel `json:"employee_benefit_detail_rincian_ekuitas" gorm:"-"`
	EmployeeBenefitDetailMutasi         []EmployeeBenefitDetailFmtEntityModel `json:"employee_benefit_detail_mutasi" gorm:"-"`
	EmployeeBenefitDetailInformasi      []EmployeeBenefitDetailFmtEntityModel `json:"employee_benefit_detail_informasi" gorm:"-"`
	EmployeeBenefitDetailAnalisis       []EmployeeBenefitDetailFmtEntityModel `json:"employee_benefit_detail_analisis" gorm:"-"`
	UserRelationModel

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

type EmployeeBenefitFilterModel struct {
	// abstraction
	abstraction.Filter

	// filter
	EmployeeBenefitFilter
	CompanyCustomFilter
}

type EmployeeBenefitVersionModel struct {
	Versions []map[int]string `json:"versions"`
}

func (EmployeeBenefitEntityModel) TableName() string {
	return "employee_benefit"
}
