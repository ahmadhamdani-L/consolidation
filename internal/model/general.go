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

type FilterModul struct {
	CompanyID *int    `query:"company_id"`
	Period    *string `query:"period" filter:"DATE"`
	Versions  *int    `query:"versions"`
	Status    *int    `query:"status"`
}
