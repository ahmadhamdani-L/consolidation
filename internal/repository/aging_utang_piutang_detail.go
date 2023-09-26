package repository

import (
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"

	"gorm.io/gorm"
)

type AgingUtangPiutangDetail interface {
	Find(ctx *abstraction.Context, m *model.AgingUtangPiutangDetailFilterModel, p *abstraction.Pagination) (*[]model.AgingUtangPiutangDetailEntityModel, *abstraction.PaginationInfo, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.AgingUtangPiutangDetailEntityModel, error)
	Create(ctx *abstraction.Context, e *model.AgingUtangPiutangDetailEntityModel) (*model.AgingUtangPiutangDetailEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.AgingUtangPiutangDetailEntityModel) (*model.AgingUtangPiutangDetailEntityModel, error)
	Delete(ctx *abstraction.Context, id *int, e *model.AgingUtangPiutangDetailEntityModel) (*model.AgingUtangPiutangDetailEntityModel, error)
	FindByCode(ctx *abstraction.Context, fmtBridgesID *int, code *string) (*model.AgingUtangPiutangDetailEntityModel, error)
	FindWithFormatter(ctx *abstraction.Context, m *model.AgingUtangPiutangDetailFilterModel) (*[]model.AgingUtangPiutangDetailFmtEntityModel, error)
	FindSummary(ctx *abstraction.Context, formatterID, formatterBridgesID *int, sort_id *float64) (*model.AgingUtangPiutangDetailEntityModel, error)
	FindByCriteria(ctx *abstraction.Context, filter *model.AgingUtangPiutangDetailFilterModel) (data *model.AgingUtangPiutangDetailEntityModel, err error)
}

type agingutangpiutangdetail struct {
	abstraction.Repository
}

func NewAgingUtangPiutangDetail(db *gorm.DB) *agingutangpiutangdetail {
	return &agingutangpiutangdetail{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *agingutangpiutangdetail) FindByCriteria(ctx *abstraction.Context, filter *model.AgingUtangPiutangDetailFilterModel) (data *model.AgingUtangPiutangDetailEntityModel, err error) {
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.AgingUtangPiutangDetailEntityModel{})
	query = r.Filter(ctx, query, *filter)
	if err = query.First(&data).Error; err != nil {
		return
	}
	// if data.ID == 0 {
	// 	err = errors.New("Data Not Found")
	// }
	return
}

func (r *agingutangpiutangdetail) Find(ctx *abstraction.Context, m *model.AgingUtangPiutangDetailFilterModel, p *abstraction.Pagination) (*[]model.AgingUtangPiutangDetailEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.AgingUtangPiutangDetailEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.AgingUtangPiutangDetailEntityModel{})
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
	query = query.Where("formatter_bridges_id IN (?)", conn.Model(&model.FormatterBridgesEntityModel{}).Select("id").Where("source = ?", "AGING-UTANG-PIUTANG").Where("trx_ref_id = ?", m.AgingUtangPiutangID))
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

func (r *agingutangpiutangdetail) FindByID(ctx *abstraction.Context, id *int) (*model.AgingUtangPiutangDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.AgingUtangPiutangDetailEntityModel
	if err := conn.Where("id = ?", &id).First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *agingutangpiutangdetail) Create(ctx *abstraction.Context, e *model.AgingUtangPiutangDetailEntityModel) (*model.AgingUtangPiutangDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}

	return e, nil
}

func (r *agingutangpiutangdetail) Update(ctx *abstraction.Context, id *int, e *model.AgingUtangPiutangDetailEntityModel) (*model.AgingUtangPiutangDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Where("id = ?", &id).First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil

}

func (r *agingutangpiutangdetail) Delete(ctx *abstraction.Context, id *int, e *model.AgingUtangPiutangDetailEntityModel) (*model.AgingUtangPiutangDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id =?", id).Delete(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil
}

func (r *agingutangpiutangdetail) FindByCode(ctx *abstraction.Context, fmtBridgesID *int, code *string) (*model.AgingUtangPiutangDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.AgingUtangPiutangDetailEntityModel
	if err := conn.Where("code ILIKE ?", *code+"%").Where("formatter_bridges_id = ?", fmtBridgesID).First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *agingutangpiutangdetail) FindWithFormatter(ctx *abstraction.Context, m *model.AgingUtangPiutangDetailFilterModel) (*[]model.AgingUtangPiutangDetailFmtEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data []model.AgingUtangPiutangDetailFmtEntityModel
	query := conn.Model(&model.FormatterBridgesEntityModel{}).
		Joins("INNER JOIN m_formatter f ON formatter_bridges.formatter_id = f.id").
		Joins("INNER JOIN m_formatter_detail fd ON f.id = fd.formatter_id").
		Joins("INNER JOIN aging_utang_piutang_detail tb ON fd.code = tb.code AND tb.formatter_bridges_id = formatter_bridges.id").
		Select("tb.*, auto_summary, is_total, is_control, is_label, control_formula").Order("fd.sort_id ASC")

	if m.FormatterBridgesID != nil && *m.FormatterBridgesID != 0 {
		query = query.Where("formatter_bridges_id = ?", m.FormatterBridgesID)
	}

	if m.AgingUtangPiutangID != nil && *m.AgingUtangPiutangID != 0 {
		query = query.Where("trx_ref_id = ? AND source = ?", m.AgingUtangPiutangID, "AGING-UTANG-PIUTANG")
	}

	err := query.Find(&data).Error
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *agingutangpiutangdetail) FindSummary(ctx *abstraction.Context, formatterID, formatterBridgesID *int, sort_id *float64) (*model.AgingUtangPiutangDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data model.AgingUtangPiutangDetailEntityModel
	var tmp model.FormatterDetailEntityModel
	err := conn.Model(&model.FormatterDetailEntityModel{}).Where("formatter_id = ? AND auto_summary IS NOT NULL AND auto_summary = true AND sort_id < ?", formatterID, sort_id).Order("sort_id DESC").Limit(1).Find(&tmp).Error
	if err != nil {
		return nil, err
	}

	lastID := 0.0
	if tmp.SortID != 0 {
		lastID = tmp.SortID
	}

	err = conn.Raw("SELECT SUM( piutangusaha_3rdparty ) piutangusaha_3rdparty, SUM ( piutangusaha_berelasi ) piutangusaha_berelasi, SUM ( piutanglainshortterm_3rdparty ) piutanglainshortterm_3rdparty, SUM ( piutanglainshortterm_berelasi ) piutanglainshortterm_berelasi, SUM ( piutangberelasishortterm ) piutangberelasishortterm, SUM ( piutanglainlongterm_3rdparty ) piutanglainlongterm_3rdparty, SUM ( piutanglainlongterm_berelasi ) piutanglainlongterm_berelasi, SUM ( piutangberelasilongterm ) piutangberelasilongterm, SUM ( utangusaha_3rdparty ) utangusaha_3rdparty, SUM ( utangusaha_berelasi ) utangusaha_berelasi, SUM ( utanglainshortterm_3rdparty ) utanglainshortterm_3rdparty, SUM ( utanglainshortterm_berelasi ) utanglainshortterm_berelasi, SUM ( utangberelasishortterm ) utangberelasishortterm, SUM ( utanglainlongterm_3rdparty ) utanglainlongterm_3rdparty, SUM ( utanglainlongterm_berelasi ) utanglainlongterm_berelasi, SUM ( utangberelasilongterm ) utangberelasilongterm FROM (SELECT tb.* FROM formatter_bridges INNER JOIN m_formatter f ON formatter_bridges.formatter_id = f.ID INNER JOIN m_formatter_detail fd ON f.ID = fd.formatter_id INNER JOIN aging_utang_piutang_detail tb ON fd.code = tb.code AND tb.formatter_bridges_id = formatter_bridges.ID WHERE ( formatter_bridges_id = ? ) AND fd.sort_id < ? AND fd.sort_id > ? ORDER BY fd.sort_id ASC) tmp ", formatterBridgesID, sort_id, lastID).Find(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}
