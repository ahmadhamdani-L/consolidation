package model

type GetVersionModel struct {
	Version []map[string]string `json:"versions"`
}

type CompanyCustomFilter struct {
	ArrCompanyID       *[]int    `filter:"CUSTOM"`
	ArrCompanyString   *[]string `filter:"CUSTOM"`
	ArrCompanyOperator *[]string `filter:"CUSTOM"`
	CompanyID          *int      `query:"company_id"`
}

type UserRelationModel struct {
	UserCreated  UserEntityModel `json:"-" gorm:"foreignKey:CreatedBy"`
	UserModified UserEntityModel `json:"-" gorm:"foreignKey:ModifiedBy"`
}
