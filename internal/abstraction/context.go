package abstraction

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Context struct {
	echo.Context
	Auth *AuthContext
	Trx  *TrxContext
}

type AuthContext struct {
	ID int
	Name      string
	Username  string
	Email	 string
	CompanyID int
	RoleID    int
}

type TrxContext struct {
	Db *gorm.DB
}
