package repository

import (
	"errors"
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	"mcash-finance-console-core/pkg/constant"

	"gorm.io/gorm"
)

type ApprovalValidation interface {
	Approve(ctx *abstraction.Context, e *model.ValidationDetailFilterModel) error
	Find(ctx *abstraction.Context, e *model.TrialBalanceFilterModel, p *abstraction.Pagination) (*[]model.TrialBalanceEntityModel, *abstraction.PaginationInfo, error)
}

type approvalvalidation struct {
	abstraction.Repository
}

func NewApprovalValidation(db *gorm.DB) *approvalvalidation {
	return &approvalvalidation{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *approvalvalidation) Approve(ctx *abstraction.Context, e *model.ValidationDetailFilterModel) error {
	conn := r.CheckTrx(ctx)

	var tmp int64
	query := conn.Model(&model.ValidationDetailEntityModel{}).Where("status != ?", constant.VALIDATION_STATUS_BALANCE)
	query = r.Filter(ctx, query, *e)
	if err := query.Count(&tmp).WithContext(ctx.Request().Context()).Error; err != nil {
		return err
	}

	if tmp != 0 {
		return errors.New("not all modul is validated")
	}

	modul := []string{
		model.AgingUtangPiutangEntityModel{}.TableName(),
		model.EmployeeBenefitEntityModel{}.TableName(),
		model.InvestasiNonTbkEntityModel{}.TableName(),
		model.InvestasiTbkEntityModel{}.TableName(),
		model.MutasiDtaEntityModel{}.TableName(),
		model.MutasiFaEntityModel{}.TableName(),
		model.MutasiIaEntityModel{}.TableName(),
		model.TrialBalanceEntityModel{}.TableName(),
		model.MutasiPersediaanEntityModel{}.TableName(),
		model.MutasiRuaEntityModel{}.TableName(),
		model.PembelianPenjualanBerelasiEntityModel{}.TableName(),
	}

	status := constant.MODUL_STATUS_CONFIRMED
	for _, v := range modul {
		queryUpdate := conn.Table(v)
		queryUpdate = r.Filter(ctx, queryUpdate, *e)
		if err := queryUpdate.Update("status", status).WithContext(ctx.Request().Context()).Error; err != nil {
			return err
		}
	}

	dataTB := model.TrialBalanceEntityModel{}
	tb := conn.Model(&dataTB)
	tb = r.Filter(ctx, tb, *e)
	if err := tb.First(&dataTB).WithContext(ctx.Request().Context()).Error; err != nil {
		return err
	}

	aje := conn.Table(model.AdjustmentEntityModel{}.TableName()).Where("tb_id = ?", dataTB.ID)
	if err := aje.Update("status", status).WithContext(ctx.Request().Context()).Error; err != nil {
		return err
	}

	return nil
}

func (r *approvalvalidation) Find(ctx *abstraction.Context, e *model.TrialBalanceFilterModel, p *abstraction.Pagination) (*[]model.TrialBalanceEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.TrialBalanceEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.TrialBalanceEntityModel{})
	//filter
	tableName := model.TrialBalanceEntityModel{}.TableName()
	query = r.FilterTable(ctx, query, *e, tableName)
	query = r.AllowedCompany(ctx, query, tableName)

	//sort
	if p.Sort == nil {
		sort := "desc"
		p.Sort = &sort
	}
	if p.SortBy == nil {
		sortBy := "created_at"
		p.SortBy = &sortBy
	}

	tmpSortBy := p.SortBy
	if p.SortBy != nil && *p.SortBy == "company" {
		sortBy := "\"Company\".name"
		p.SortBy = &sortBy
	}
	if p.SortBy != nil && (tmpSortBy != nil && *tmpSortBy != "company" && *tmpSortBy != "user") {
		sortBy := fmt.Sprintf("trial_balance.%s", *p.SortBy)
		p.SortBy = &sortBy
	}

	sort := fmt.Sprintf("%s %s", *p.SortBy, *p.Sort)
	query = query.Order(sort)
	p.SortBy = tmpSortBy

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
	var totalData int64
	query = query.Joins("Company").Joins("INNER JOIN validation_detail vd ON trial_balance.company_id = vd.company_id AND trial_balance.period = vd.period AND trial_balance.versions = vd.versions").Where("trial_balance.status = 1").Where("trial_balance.validation_note = ?", constant.VALIDATION_NOTE_BALANCE).Group("trial_balance.id, \"Company\".id").Having("SUM ( CASE WHEN vd.status = 2 THEN 0 ELSE 1 END ) = 0").Select("\"trial_balance\".\"id\",\"trial_balance\".\"created_at\",\"trial_balance\".\"created_by\",\"trial_balance\".\"modified_at\",\"trial_balance\".\"modified_by\",\"trial_balance\".\"period\",\"trial_balance\".\"versions\",\"trial_balance\".\"company_id\",\"trial_balance\".\"status\",\"trial_balance\".\"validation_note\",\"Company\".\"id\" AS \"Company__id\",\"Company\".\"created_at\" AS \"Company__created_at\",\"Company\".\"created_by\" AS \"Company__created_by\",\"Company\".\"modified_at\" AS \"Company__modified_at\",\"Company\".\"modified_by\" AS \"Company__modified_by\",\"Company\".\"code\" AS \"Company__code\",\"Company\".\"name\" AS \"Company__name\",\"Company\".\"pic\" AS \"Company__pic\",\"Company\".\"parent_company_id\" AS \"Company__parent_company_id\",\"Company\".\"is_active\" AS \"Company__is_active\"")
	query = query.Count(&totalData).Limit(limit).Offset(offset)

	if err := query.Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, nil, err
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
