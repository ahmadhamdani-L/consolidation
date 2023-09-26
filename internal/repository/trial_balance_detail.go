package repository

import (
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	"mcash-finance-console-core/pkg/util/helper"
	"strconv"

	"gorm.io/gorm"
)

type TrialBalanceDetail interface {
	Find(ctx *abstraction.Context, m *model.TrialBalanceDetailFilterModel, p *abstraction.Pagination) (*[]model.TrialBalanceDetailEntityModel, *abstraction.PaginationInfo, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.TrialBalanceDetailEntityModel, error)
	Create(ctx *abstraction.Context, e *model.TrialBalanceDetailEntityModel) (*model.TrialBalanceDetailEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.TrialBalanceDetailEntityModel) (*model.TrialBalanceDetailEntityModel, error)
	Delete(ctx *abstraction.Context, id *int, e *model.TrialBalanceDetailEntityModel) (*model.TrialBalanceDetailEntityModel, error)
	FindToExport(ctx *abstraction.Context, code *string, fmtBridgeID *int) (*[]model.TrialBalanceDetailEntityModel, error)
	Import(ctx *abstraction.Context, e *[]model.TrialBalanceDetailEntityModel) (*[]model.TrialBalanceDetailEntityModel, error)
	// FindWithFormatter(ctx *abstraction.Context, m *model.TrialBalanceDetailFilterModel) (*[]model.TrialBalanceDetailFmtEntityModel, error)
	FindWithFormatter(ctx *abstraction.Context, m *model.TrialBalanceDetailFilterModel, p *abstraction.Pagination) (*[]model.TrialBalanceDetailFmtEntityModel, *abstraction.PaginationInfo, error)
	FindSummary(ctx *abstraction.Context, codeCoa *string, formatterBridgesID *int, isCoa *bool) (*model.TrialBalanceDetailEntityModel, error)
	FindDetail(ctx *abstraction.Context, trialBalanceID *int, parentID *int) (*[]model.TrialBalanceDetailFmtEntityModel, error)
	FindAllDetail(ctx *abstraction.Context, trialBalanceID *int) (*[]model.TrialBalanceDetailFmtEntityModel, error)
	FindByCode(ctx *abstraction.Context, e *model.TrialBalanceDetailFilterModel) (*model.TrialBalanceDetailEntityModel, error)
	SummaryByCodes(ctx *abstraction.Context, bridgeID *int, codes []string) (*model.TrialBalanceDetailEntityModel, error)
	FindByExactCode(ctx *abstraction.Context, e *model.TrialBalanceDetailFilterModel) (*model.TrialBalanceDetailEntityModel, error)
	FindControlWbs1(ctx *abstraction.Context, trialBalanceID *int) (*model.TrialBalanceDetailEntityModel, error)
	FindByCriteriaTb(ctx *abstraction.Context, filter *model.TrialBalanceFilterModel) (data *model.TrialBalanceEntityModel, err error)
	FindSummaryTb(ctx *abstraction.Context, m *model.TrialBalanceDetailFilterModel) (*model.TrialBalanceDetailEntityModel, error)
	
}

type trialbalancedetail struct {
	abstraction.Repository
}

func NewTrialBalanceDetail(db *gorm.DB) *trialbalancedetail {
	return &trialbalancedetail{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *trialbalancedetail) FindSummaryTb(ctx *abstraction.Context, m *model.TrialBalanceDetailFilterModel) (*model.TrialBalanceDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas model.TrialBalanceDetailEntityModel

	query := conn.Model(&model.TrialBalanceDetailEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)
	query = query.Select("SUM(amount_before_aje) amount_before_aje, SUM(amount_after_aje) amount_after_aje, SUM(amount_aje_cr) amount_aje_cr, SUM(amount_aje_dr) amount_aje_dr")
	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}
func (r *trialbalancedetail) FindByCriteriaTb(ctx *abstraction.Context, filter *model.TrialBalanceFilterModel) (data *model.TrialBalanceEntityModel, err error) {
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.TrialBalanceEntityModel{})
	if err = query.Where("company_id = ?", filter.CompanyID).Where("period = ?", filter.Period).Where("versions = ?", filter.Versions).First(&data).Error; err != nil {
		return
	}
	// if data.ID == 0 {
	// 	err = errors.New("Data Not Found")
	// }
	return
}

func (r *trialbalancedetail) Find(ctx *abstraction.Context, m *model.TrialBalanceDetailFilterModel, p *abstraction.Pagination) (*[]model.TrialBalanceDetailEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.TrialBalanceDetailEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.TrialBalanceDetailEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	//sort
	if p.Sort == nil {
		sort := "asc"
		p.Sort = &sort
	}
	if p.SortBy == nil {
		sortBy := "id"
		p.SortBy = &sortBy
	}

	sort := fmt.Sprintf("%s %s", *p.SortBy, *p.Sort)
	query = query.Order(sort)

	//pagination
	if p.Page == nil {
		page := 1
		p.Page = &page
	}
	if p.PageSize == nil {
		pageSize := 10
		p.PageSize = &pageSize
	}
	info = abstraction.PaginationInfo{
		Pagination: p,
	}
	limit := *p.PageSize
	offset := limit * (*p.Page - 1)
	if m.TrialBalanceID != nil {
		query = query.Where("formatter_bridges_id = (?)", conn.Model(&model.FormatterBridgesEntityModel{}).Select("id").Where("source = ?", "TRIAL-BALANCE").Where("trx_ref_id = ?", m.TrialBalanceID).Limit(1))
	}

	var totalData int64
	if *p.PageSize != -1 {
		query = query.Count(&totalData).Limit(limit).Offset(offset)
	} else {
		query = query.Count(&totalData)
	}
	if err := query.Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
		return &datas, &info, err
	}

	info.Count = int(totalData)
	info.MoreRecords = false
	if len(datas) > *p.PageSize {
		info.MoreRecords = true
		// info.Count -= 1
		// datas = datas[:len(datas)-1]
	}

	return &datas, &info, nil
}

func (r *trialbalancedetail) FindToExport(ctx *abstraction.Context, code *string, fmtBridgeID *int) (*[]model.TrialBalanceDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.TrialBalanceDetailEntityModel
	if err := conn.Where("code ILIKE ?", *code+"%").Where("formatter_bridges_id = ?", *fmtBridgeID).Find(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *trialbalancedetail) FindByID(ctx *abstraction.Context, id *int) (*model.TrialBalanceDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.TrialBalanceDetailEntityModel
	if err := conn.Where("id = ?", &id).First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *trialbalancedetail) Create(ctx *abstraction.Context, e *model.TrialBalanceDetailEntityModel) (*model.TrialBalanceDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}

	return e, nil
}

func (r *trialbalancedetail) Update(ctx *abstraction.Context, id *int, e *model.TrialBalanceDetailEntityModel) (*model.TrialBalanceDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Where("id = ?", &id).First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil

}

func (r *trialbalancedetail) Delete(ctx *abstraction.Context, id *int, e *model.TrialBalanceDetailEntityModel) (*model.TrialBalanceDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id =?", id).Delete(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil
}

func (r *trialbalancedetail) Import(ctx *abstraction.Context, e *[]model.TrialBalanceDetailEntityModel) (*[]model.TrialBalanceDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}

	return e, nil
}

// func (r *trialbalancedetail) FindWithFormatter(ctx *abstraction.Context, m *model.TrialBalanceDetailFilterModel) (*[]model.TrialBalanceDetailFmtEntityModel, error) {
// 	conn := r.CheckTrx(ctx)
// 	var data []model.TrialBalanceDetailFmtEntityModel
// 	err := conn.Model(&model.FormatterBridgesEntityModel{}).
// 		Joins("INNER JOIN m_formatter f ON formatter_bridges.formatter_id = f.id").
// 		Joins("INNER JOIN m_formatter_detail fd ON f.id = fd.formatter_id").
// 		Joins("INNER JOIN trial_balance_detail tb ON fd.code = tb.code AND tb.formatter_bridges_id = formatter_bridges.id").
// 		Where("trx_ref_id = ? AND source = ?", m.TrialBalanceID, "TRIAL-BALANCE").
// 		Order("fd.sort_id ASC").
// 		Select("tb.*, is_total, is_control, control_formula").Find(&data).Error
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &data, nil
// }

func (r *trialbalancedetail) FindWithFormatter(ctx *abstraction.Context, m *model.TrialBalanceDetailFilterModel, p *abstraction.Pagination) (*[]model.TrialBalanceDetailFmtEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)
	var data []model.TrialBalanceDetailFmtEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.FormatterBridgesEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	//sort
	// if p.Sort == nil {
	// 	sort := "desc"
	// 	p.Sort = &sort
	// }
	// if p.SortBy == nil {
	// 	sortBy := "id"
	// 	p.SortBy = &sortBy
	// }

	// sort := fmt.Sprintf("%s %s", *p.SortBy, *p.Sort)
	// query = query.Order(sort)

	//pagination
	if p.Page == nil {
		page := 1
		p.Page = &page
	}
	if p.PageSize == nil {
		pageSize := 10
		p.PageSize = &pageSize
	}
	info = abstraction.PaginationInfo{
		Pagination: p,
	}
	limit := *p.PageSize
	offset := limit * (*p.Page - 1)
	query = query.Joins("INNER JOIN m_formatter f ON formatter_bridges.formatter_id = f.id").Joins("INNER JOIN m_formatter_detail fd ON f.id = fd.formatter_id").
		Joins("RIGHT JOIN trial_balance_detail tb ON tb.formatter_bridges_id = formatter_bridges.ID AND (tb.code = fd.code OR SUBSTR(tb.code, 0, 7) = fd.code)").Where("trx_ref_id = ? AND source = ?", m.TrialBalanceID, "TRIAL-BALANCE").
		Order("fd.sort_id ASC").Order("tb.id ASC").Select("tb.*, auto_summary, is_total, is_control, is_label, control_formula")

	querycount := query
	var totalData int64
	if err := querycount.Count(&totalData).Error; err != nil {
		return &data, &info, err
	}
	query = query.Count(&totalData).Limit(limit).Offset(offset)
	if err := query.Find(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, &info, err
	}

	info.Count = int(totalData)
	info.MoreRecords = false
	if len(data) > *p.PageSize {
		info.MoreRecords = true
		info.Count -= 1
		data = data[:len(data)-1]
	}

	return &data, &info, nil
}

func (r *trialbalancedetail) FindSummary(ctx *abstraction.Context, codeCoa *string, formatterBridgesID *int, isCoa *bool) (*model.TrialBalanceDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data model.TrialBalanceDetailEntityModel

	tmpStr := *codeCoa
	if _, err := strconv.Atoi(*codeCoa); err == nil {
		tmpStr += "%"
	}

	query := fmt.Sprintf("SELECT SUM( amount_before_aje ) amount_before_aje, SUM( amount_aje_dr ) amount_aje_dr, SUM( amount_aje_cr ) amount_aje_cr, SUM( amount_after_aje ) amount_after_aje FROM trial_balance_detail WHERE formatter_bridges_id = ? AND LOWER(description) != 'sub total' AND code LIKE '%s'", tmpStr)
	err := conn.Raw(query, &formatterBridgesID).Find(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *trialbalancedetail) FindDetail(ctx *abstraction.Context, trialBalanceID *int, parentID *int) (*[]model.TrialBalanceDetailFmtEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var tb model.TrialBalanceEntityModel
	var fmtBridgeData model.FormatterBridgesEntityModel

	if err := conn.Model(&model.TrialBalanceEntityModel{}).Where("id = ?", trialBalanceID).First(&tb).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(&model.FormatterBridgesEntityModel{}).Where("trx_ref_id = ? AND source = 'TRIAL-BALANCE'", tb.ID).First(&fmtBridgeData).Error; err != nil {
		return nil, err
	}
	query := conn.Model(&model.TrialBalanceDetailEntityModel{}).Joins("INNER JOIN m_formatter_detail fd ON (trial_balance_detail.code = fd.code OR SUBSTR(trial_balance_detail.code, 0, 7) = fd.code) AND fd.formatter_id = 3")
	query = query.Select("trial_balance_detail.id, fd.id as formatter_detail_id, trial_balance_detail.code, trial_balance_detail.description, amount_before_aje, amount_aje_cr, amount_aje_dr, amount_after_aje, parent_id, auto_summary, is_total, is_control, is_parent, is_label, control_formula")
	if parentID != nil && *parentID != 0 {
		query = query.Where("parent_id = ?", parentID)
	} else {
		query = query.Where("parent_id IS NULL")
	}
	query = query.Where("formatter_bridges_id = ?", fmtBridgeData.ID).Order("fd.sort_id ASC")
	var data []model.TrialBalanceDetailFmtEntityModel
	if err := query.Find(&data).Error; err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *trialbalancedetail) FindAllDetail(ctx *abstraction.Context, trialBalanceID *int) (*[]model.TrialBalanceDetailFmtEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var tb model.TrialBalanceEntityModel
	var fmtBridgeData model.FormatterBridgesEntityModel

	if err := conn.Model(&model.TrialBalanceEntityModel{}).Where("id = ?", trialBalanceID).First(&tb).Error; err != nil {
		return nil, err
	}

	if err := conn.Model(&model.FormatterBridgesEntityModel{}).Where("trx_ref_id = ? AND source = 'TRIAL-BALANCE'", tb.ID).First(&fmtBridgeData).Error; err != nil {
		return nil, err
	}

	// query untuk mengambil data trial balance detail dengan formatternya berdasarkan formatter bridge, datanya digabung menggunakan union karena untuk formatter detail 950 dan 960 itu coa nya berbeda2 jadi dihardcode, juga karena perbedaan substring dari code coa nya. code yang di comment dibawah dari code ini adalah kode yang sudah benar tetapi karena tidak lengkap (tidak ada data 950 dan 960) maka dibuatlah union untuk mengambil data dari formatter detail 950 dan 960
	tmpCon1 := "SUBSTR( trial_balance_detail.code, 0, 7 ) = fd.code"
	tmpCon2 := "SUBSTR( trial_balance_detail.code, 0, 4 ) = fd.code"

	querySelect := "SELECT trial_balance_detail.ID, fd.ID AS formatter_detail_id, trial_balance_detail.code, trial_balance_detail.description, amount_before_aje, amount_aje_cr, amount_aje_dr, amount_after_aje, fd.sort_id, fd.parent_id, fd.auto_summary, fd.is_total, fd.is_control, fd.is_parent, fd.is_label, fd.control_formula, fd.is_coa, fd.code formatter_detail_code, fd.description formatter_detail_description, fd.show_group_coa"

	tmpTB := " FROM \"trial_balance_detail\" INNER JOIN m_formatter_detail fd ON fd.formatter_id = ? AND fd.is_show_view = true AND "

	queryStr := fmt.Sprintf("WITH tbxfmt AS (%s %s (trial_balance_detail.code = fd.code OR "+tmpCon1+") WHERE formatter_bridges_id = ?), tbxfmtc AS (%s %s (trial_balance_detail.code = fd.code OR "+tmpCon2+") WHERE formatter_bridges_id = ? AND fd.code IN ('950', '960'))", querySelect, tmpTB, querySelect, tmpTB)

	query := conn.Raw(fmt.Sprintf("%s SELECT * FROM tbxfmt UNION ALL SELECT * FROM tbxfmtc ORDER BY sort_id, id", queryStr), fmtBridgeData.FormatterID, fmtBridgeData.ID, fmtBridgeData.FormatterID, fmtBridgeData.ID)

	// query := conn.Model(&model.TrialBalanceDetailEntityModel{}).Joins("INNER JOIN m_formatter_detail fd ON (trial_balance_detail.code = fd.code OR (CASE WHEN fd.code = '950' OR fd.code = '960' THEN SUBSTR( trial_balance_detail.code, 0, 4) ELSE SUBSTR( trial_balance_detail.code, 0, 7 ) END) = fd.code ) AND fd.formatter_id = ?", fmtBridgeData.FormatterID)
	// query = query.Select("trial_balance_detail.id, fd.id as formatter_detail_id, trial_balance_detail.code, trial_balance_detail.description, amount_before_aje, amount_aje_cr, amount_aje_dr, amount_after_aje, fd.sort_id, parent_id, auto_summary, is_total, is_control, is_parent, is_label, control_formula, is_coa, fd.code formatter_detail_code, fd.description formatter_detail_description")
	// query = query.Where("formatter_bridges_id = ?", fmtBridgeData.ID).Order("fd.sort_id ASC, trial_balance_detail.id ASC")
	var data []model.TrialBalanceDetailFmtEntityModel
	if err := query.Find(&data).Error; err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *trialbalancedetail) FindByCode(ctx *abstraction.Context, e *model.TrialBalanceDetailFilterModel) (*model.TrialBalanceDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data model.TrialBalanceDetailEntityModel
	query := conn.Model(&model.TrialBalanceDetailEntityModel{})
	query = r.Filter(ctx, query, *e)
	if err := query.First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *trialbalancedetail) SummaryByCodes(ctx *abstraction.Context, bridgeID *int, codes []string) (*model.TrialBalanceDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data model.TrialBalanceDetailEntityModel
	query := conn.Model(&model.TrialBalanceDetailEntityModel{}).Where("formatter_bridges_id = ?", bridgeID).Where("code IN (?)", codes).Select("SUM(amount_before_aje) amount_before_aje, SUM(amount_aje_dr) amount_aje_dr, SUM(amount_aje_cr) amount_aje_cr, SUM(amount_after_aje) amount_after_aje")
	if err := query.Take(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *trialbalancedetail) FindByExactCode(ctx *abstraction.Context, e *model.TrialBalanceDetailFilterModel) (*model.TrialBalanceDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data model.TrialBalanceDetailEntityModel
	query := conn.Model(&model.TrialBalanceDetailEntityModel{}).Where("code = ?", *e.Code).Where("formatter_bridges_id = ?", *e.FormatterBridgesID)
	if err := query.First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *trialbalancedetail) FindControlWbs1(ctx *abstraction.Context, trialBalanceID *int) (*model.TrialBalanceDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var fmtBridges model.FormatterBridgesEntityModel
	if err := conn.Model(&model.FormatterBridgesEntityModel{}).Where("trx_ref_id = ? AND source = 'TRIAL-BALANCE'", trialBalanceID).First(&fmtBridges).Error; err != nil {
		return nil, err
	}

	var data model.TrialBalanceDetailEntityModel
	queryFrom := "( SELECT formatter_bridges_id, COALESCE(amount_before_aje, 0) amount_before_aje, COALESCE(amount_aje_cr, 0) amount_aje_cr, COALESCE(amount_aje_dr, 0) amount_aje_dr, COALESCE(amount_after_aje, 0) amount_after_aje FROM trial_balance_detail WHERE code = ? AND formatter_bridges_id = ? )"
	query := fmt.Sprintf("SELECT ( aset.amount_before_aje - liaeku.amount_before_aje ) amount_before_aje, ( aset.amount_aje_cr + liaeku.amount_aje_cr ) amount_aje_cr, ( aset.amount_aje_dr + liaeku.amount_aje_dr ) amount_aje_dr, ( aset.amount_after_aje - liaeku.amount_after_aje ) amount_after_aje FROM %s aset JOIN %s liaeku ON aset.formatter_bridges_id = liaeku.formatter_bridges_id", queryFrom, queryFrom)
	if err := conn.Raw(query, "TOTAL_ASET", fmtBridges.ID, "TOTAL_LIABILITAS_DAN_EKUITAS", fmtBridges.ID).Find(&data).Error; err != nil {
		return nil, err
	}

	data.AmountBeforeAje = helper.AssignAmount(data.AmountBeforeAje)
	data.AmountAjeCr = helper.AssignAmount(data.AmountAjeCr)
	data.AmountAjeDr = helper.AssignAmount(data.AmountAjeDr)
	data.AmountAfterAje = helper.AssignAmount(data.AmountAfterAje)

	return &data, nil
}
