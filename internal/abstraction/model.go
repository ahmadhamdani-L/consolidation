package abstraction

import (
	"time"
	"worker/pkg/util/date"

	"gorm.io/gorm"
)

type EntityImportedWorksheetDetail struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement;"`
	ModifiedAt         *time.Time `json:"modified_at"`
}

type Entity struct {
	ID int `json:"id" gorm:"primaryKey;autoIncrement;"`

	CreatedAt          time.Time  `json:"created_at"`
	CreatedBy          int        `json:"created_by"`
	UserCreatedString  string     `json:"user_created" gorm:"-"`
	ModifiedAt         *time.Time `json:"modified_at"`
	ModifiedBy         *int       `json:"modified_by"`
	UserModifiedString *string    `json:"user_modified" gorm:"-"`

	// DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Filter struct {
	CreatedAt          *time.Time `query:"created_at" filter:"DATE" example:"2022-08-17T15:04:05Z"`
	CreatedBy          *int       `query:"created_by" example:"1"`
	UserCreatedString  *string    `query:"user_created" filter:"CUSTOM" example:"Lutfi Ramadhan"`
	ModifiedAt         *time.Time `query:"modified_at" filter:"DATE" example:"2022-08-17T15:04:05Z"`
	ModifiedBy         *int       `query:"modified_by" example:"1"`
	UserModifiedString *string    `query:"user_modified" filter:"CUSTOM" example:"Lutfi Ramadhan"`
}

type EntityFormatter struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement;"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy int       `json:"created_by"`

	// DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
type FilterFormatter struct {
	CreatedAt *time.Time `query:"created_at"`
	CreatedBy *int       `query:"created_by"`
}

func (m *Entity) BeforeUpdate(tx *gorm.DB) (err error) {
	m.ModifiedAt = date.DateTodayLocal()
	return
}

func (m *Entity) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = *date.DateTodayLocal()
	return
}

type JsonDataImport struct {
	TrialBalance               string
	AgingUtangPiutang          string
	InvestasiTbk               string
	InvestasiNonTbk            string
	MutasiFA                   string
	MutasiDta                  string
	MutasiIa                   string
	MutasiRua                  string
	MutasiPersediaan           string
	PembelianPenjualanBerelasi string
	EmployeeBenefit			   string
	FNTrialBalance               string
	FNAgingUtangPiutang          string
	FNInvestasiTbk               string
	FNInvestasiNonTbk            string
	FNMutasiFA                   string
	FNMutasiDta                  string
	FNMutasiIa                   string
	FNMutasiRua                  string
	FNMutasiPersediaan           string
	FNPembelianPenjualanBerelasi string
	FNEmployeeBenefit			 string
	CompanyID                  int
	UserID                     int
	Version                    int
	ImportedWorkSheetID        int
	Period					   string
	Data                       string
}

type JsonDataReUpload struct {

	TrialBalance                 string
	AgingUtangPiutang            string
	InvestasiTbk                 string
	InvestasiNonTbk              string
	MutasiFA                     string
	MutasiDta                    string
	MutasiIa                     string
	MutasiRua                    string
	MutasiPersediaan             string
	PembelianPenjualanBerelasi   string
	EmployeeBenefit			     string
	IDTrialBalance               int
	IDAgingUtangPiutang          int
	IDInvestasiTbk               int
	IDInvestasiNonTbk            int
	IDMutasiFA                   int
	IDMutasiDta                  int
	IDMutasiIa                   int
	IDMutasiRua                  int
	IDMutasiPersediaan           int
	IDPembelianPenjualanBerelasi int
	IDEmployeeBenefit 			 int
	CompanyID                    int
	UserID                       int
	Version                      int
	ImportedWorkSheetID          int
	IDWorksheetDetailTrialBalance               int
	IDWorksheetDetailAgingUtangPiutang          int
	IDWorksheetDetailInvestasiTbk               int
	IDWorksheetDetailInvestasiNonTbk            int
	IDWorksheetDetailMutasiFA                   int
	IDWorksheetDetailMutasiDta   				int
	IDWorksheetDetailMutasiIa                   int
	IDWorksheetDetailMutasiRua                  int
	IDWorksheetDetailMutasiPersediaan           int
	IDWorksheetDetailPembelianPenjualanBerelasi int
	IDWorksheetDetailEmployeeBenefit            int
	FNTrialBalance								string
}

