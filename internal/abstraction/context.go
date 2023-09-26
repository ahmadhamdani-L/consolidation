package abstraction

import (
	"time"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Context struct {
	echo.Context
	Trx  *TrxContext
	Auth *AuthContext
}

type AuthContext struct {
	ID        int
	Name      string
	CompanyID int
	Time      *time.Time
}

type TrxContext struct {
	Db *gorm.DB
}
