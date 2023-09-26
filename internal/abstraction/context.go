package abstraction

import (
	"time"

	"gorm.io/gorm"
)

type Context struct {
	Trx  *TrxContext
	Auth *AuthContext
}

type AuthContext struct {
	ID        int
	Name      string
	CompanyID int
	Time      *time.Time
	// RoleID    int
}

type TrxContext struct {
	Db *gorm.DB
}
