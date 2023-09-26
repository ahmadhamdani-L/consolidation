package model

type GetVersionModel struct {
	Version []map[int]string `json:"versions"`
}

type CompanyCustomFilter struct {
	ArrCompanyID       *[]int    `filter:"CUSTOM"`
	ArrCompanyString   *[]string `filter:"CUSTOM"`
	ArrCompanyOperator *[]string `filter:"CUSTOM"`
	CompanyID          *int      `query:"company_id"`
	CompanyString      *string   `query:"company" filter:"CUSTOM"`
	CompanyOperator    *string   `query:"company_operator" filter:"CUSTOM"`
}

type ExportFilter struct {
	CompanyID int
	Period    string
	Version   int
}

type UserRelationModel struct {
	UserCreated  UserEntityModel `json:"-" gorm:"foreignKey:CreatedBy"`
	UserModified UserEntityModel `json:"-" gorm:"foreignKey:ModifiedBy"`
}
