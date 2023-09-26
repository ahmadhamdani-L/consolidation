package repository

import (
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"

	"gorm.io/gorm"
)

type MutasiFaDetail interface {
	Find(ctx *abstraction.Context, m *model.MutasiFaDetailFilterModel, p *abstraction.Pagination) (*[]model.MutasiFaDetailEntityModel, *abstraction.PaginationInfo, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.MutasiFaDetailEntityModel, error)
	Create(ctx *abstraction.Context, e *model.MutasiFaDetailEntityModel) (*model.MutasiFaDetailEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.MutasiFaDetailEntityModel) (*model.MutasiFaDetailEntityModel, error)
	Delete(ctx *abstraction.Context, id *int, e *model.MutasiFaDetailEntityModel) (*model.MutasiFaDetailEntityModel, error)
	FindWithCode(ctx *abstraction.Context, code *string) (*[]model.MutasiFaDetailEntityModel, error)
	FindWithFormatter(ctx *abstraction.Context, m *model.MutasiFaDetailFilterModel) (*[]model.MutasiFaDetailFmtEntityModel, error)
	FindSummary(ctx *abstraction.Context, formatterID, formatterBridgesID *int, sort_id *float64) (*model.MutasiFaDetailEntityModel, error)
	FindTotal(ctx *abstraction.Context, e *model.MutasiFaDetailFilterModel) (*model.MutasiFaDetailEntityModel, error)
	FindByCriteria(ctx *abstraction.Context, filter *model.MutasiFaDetailFilterModel) (data *model.MutasiFaDetailEntityModel, err error)
}

type mutasifadetail struct {
	abstraction.Repository
}

func NewMutasiFaDetail(db *gorm.DB) *mutasifadetail {
	return &mutasifadetail{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *mutasifadetail) FindByCriteria(ctx *abstraction.Context, filter *model.MutasiFaDetailFilterModel) (data *model.MutasiFaDetailEntityModel, err error) {
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.MutasiFaDetailEntityModel{})
	query = r.Filter(ctx, query, *filter)
	if err = query.First(&data).Error; err != nil {
		return
	}
	// if data.ID == 0 {
	// 	err = errors.New("Data Not Found")
	// }
	return
}

func (r *mutasifadetail) Find(ctx *abstraction.Context, m *model.MutasiFaDetailFilterModel, p *abstraction.Pagination) (*[]model.MutasiFaDetailEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.MutasiFaDetailEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.MutasiFaDetailEntityModel{})
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
	if m.MutasiFaID != nil {
		query = query.Where("formatter_bridges_id IN (?)", conn.Model(&model.FormatterBridgesEntityModel{}).Select("id").Where("source = ?", "MUTASI-FA").Where("trx_ref_id = ?", m.MutasiFaID))
	}
	var totalData int64
	query = query.Count(&totalData).Limit(limit).Offset(offset)
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

func (r *mutasifadetail) FindByID(ctx *abstraction.Context, id *int) (*model.MutasiFaDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.MutasiFaDetailEntityModel
	if err := conn.Where("id = ?", &id).First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *mutasifadetail) Create(ctx *abstraction.Context, e *model.MutasiFaDetailEntityModel) (*model.MutasiFaDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}

	return e, nil
}

func (r *mutasifadetail) Update(ctx *abstraction.Context, id *int, e *model.MutasiFaDetailEntityModel) (*model.MutasiFaDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Where("id = ?", &id).First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil

}

func (r *mutasifadetail) Delete(ctx *abstraction.Context, id *int, e *model.MutasiFaDetailEntityModel) (*model.MutasiFaDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id =?", id).Delete(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil
}

func (r *mutasifadetail) FindWithCode(ctx *abstraction.Context, code *string) (*[]model.MutasiFaDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.MutasiFaDetailEntityModel
	if err := conn.Where("code ILIKE ?", *code+"%").Find(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *mutasifadetail) FindWithFormatter(ctx *abstraction.Context, m *model.MutasiFaDetailFilterModel) (*[]model.MutasiFaDetailFmtEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data []model.MutasiFaDetailFmtEntityModel
	// err := conn.Model(&model.FormatterBridgesEntityModel{}).
	// 	Joins("INNER JOIN m_formatter f ON formatter_bridges.formatter_id = f.id").
	// 	Joins("INNER JOIN m_formatter_detail fd ON f.id = fd.formatter_id").
	// 	Joins("INNER JOIN mutasi_fa_detail tb ON fd.code = tb.code AND tb.formatter_bridges_id = formatter_bridges.id").
	// 	Where("formatter_bridges_id = ?", m.FormatterBridgesID).Order("fd.sort_id ASC").
	// 	Select("tb.*, auto_summary, is_total, is_control, is_label, control_formula").Find(&data).Error
	// if err != nil {
	// 	return nil, err
	// }
	query := conn.Model(&model.FormatterBridgesEntityModel{}).
		Joins("INNER JOIN m_formatter f ON formatter_bridges.formatter_id = f.id").
		Joins("INNER JOIN m_formatter_detail fd ON f.id = fd.formatter_id").
		Joins("INNER JOIN mutasi_fa_detail tb ON fd.code = tb.code AND tb.formatter_bridges_id = formatter_bridges.id").
		Select("tb.*, auto_summary, is_total, is_control, is_label, control_formula").Order("fd.sort_id ASC")

	if m.FormatterBridgesID != nil && *m.FormatterBridgesID != 0 {
		query = query.Where("formatter_bridges_id = ?", m.FormatterBridgesID)
	}

	if m.MutasiFaID != nil && *m.MutasiFaID != 0 {
		query = query.Where("trx_ref_id = ? AND source = ?", m.MutasiFaID, "MUTASI-FA")
	}

	err := query.Find(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *mutasifadetail) FindSummary(ctx *abstraction.Context, formatterID, formatterBridgesID *int, sort_id *float64) (*model.MutasiFaDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data model.MutasiFaDetailEntityModel
	var tmp model.FormatterDetailEntityModel
	err := conn.Model(&model.FormatterDetailEntityModel{}).Where("formatter_id = ? AND auto_summary IS NOT NULL AND auto_summary = true AND sort_id < ?", formatterID, sort_id).Order("sort_id DESC").Limit(1).Find(&tmp).Error
	if err != nil {
		return nil, err
	}

	lastID := 0.0
	if tmp.SortID != 0 {
		lastID = tmp.SortID
	}

	err = conn.Raw("SELECT SUM( beginning_balance ) beginning_balance, SUM( acquisition_of_subsidiary ) acquisition_of_subsidiary, SUM( additions ) additions, SUM( deductions ) deductions, SUM( reclassification ) reclassification, SUM( revaluation ) revaluation, SUM( ending_balance ) ending_balance FROM (SELECT tb.* FROM formatter_bridges INNER JOIN m_formatter f ON formatter_bridges.formatter_id = f.ID INNER JOIN m_formatter_detail fd ON f.ID = fd.formatter_id INNER JOIN mutasi_fa_detail tb ON fd.code = tb.code AND tb.formatter_bridges_id = formatter_bridges.ID WHERE ( formatter_bridges_id = ? ) AND fd.sort_id < ? AND fd.sort_id > ? ORDER BY fd.sort_id ASC) tmp ", formatterBridgesID, sort_id, lastID).Find(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *mutasifadetail) FindTotal(ctx *abstraction.Context, e *model.MutasiFaDetailFilterModel) (*model.MutasiFaDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data model.MutasiFaDetailEntityModel
	query := conn.Model(&model.MutasiFaDetailEntityModel{})
	query = query.Where("source = ? AND trx_ref_id = ?", "MUTASI-FA", e.MutasiFaID)
	query = query.Joins("INNER JOIN formatter_bridges ON formatter_bridges.id = formatter_bridges_id AND code = ?", e.Code)

	if err := query.First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return &data, nil
}
