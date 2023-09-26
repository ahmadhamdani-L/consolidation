package repository

import (
	"fmt"
	"worker-validation/internal/abstraction"
	"worker-validation/internal/model"

	"gorm.io/gorm"
)

type Export interface {
	GetData(ctx *abstraction.Context, m model.FilterData, tipe string) ([]interface{}, error)
}

type export struct {
	abstraction.Repository
}

func NewExport(db *gorm.DB) *company {
	return &company{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *export) GetData(ctx *abstraction.Context, m model.FilterData, tipe string) ([]interface{}, error) {
	conn := r.CheckTrx(ctx)
	var datas []interface{}
	var tableMain, src string

	switch tipe {
	case "AUP":
		tableMain = "aging_utang_piutang"
		src = "AGING-UTANG-PIUTANG"
	case "IT":
		tableMain = "investasi_tbk"
		src = "INVESTASI-TBK"
	case "MUTASI-FA":
		tableMain = "mutasi_fa"
		src = "MUTASI-FA"
	case "TB":
		tableMain = "trial_balance"
		src = "TRIAL-BALANCE"
	}

	// query := conn.Model(&model.AgingUtangPiutangEntityEntityModel{})
	query := conn.Table("? tbm", tableMain).
		Joins("INNER JOIN formatter_bridges fb ON tbm.id = fb.trx_ref_id AND fb.source = ?", src).
		Joins("INNER JOIN m_formatter_detail fd ON fb.formatter_id = fd.formatter_id").
		Joins(fmt.Sprintf("INNER JOIN %s_detail aupd ON aupd.formatter_bridges_id = fb.id AND aupd.code = fd.code", tableMain)).
		Where("aup.company_id = ? AND aup.period = ? AND aup.versions = ?", m.CompanyID, m.Period, m.Versions).
		Order("fd.sort_id ASC")
	if err := query.Find(&datas).Error; err != nil {
		return datas, err
	}

	return datas, nil
}
