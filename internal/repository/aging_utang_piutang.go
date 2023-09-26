package repository

import (
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"

	"gorm.io/gorm"
)

type AgingUtangPiutang interface {
	Find(ctx *abstraction.Context, m *model.AgingUtangPiutangFilterModel, p *abstraction.Pagination) (*[]model.AgingUtangPiutangEntityModel, *abstraction.PaginationInfo, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.AgingUtangPiutangEntityModel, error)
	Create(ctx *abstraction.Context, e *model.AgingUtangPiutangEntityModel) (*model.AgingUtangPiutangEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.AgingUtangPiutangEntityModel) (*model.AgingUtangPiutangEntityModel, error)
	Delete(ctx *abstraction.Context, id *int, e *model.AgingUtangPiutangEntityModel) (*model.AgingUtangPiutangEntityModel, error)
	GetCount(ctx *abstraction.Context, m *model.AgingUtangPiutangFilterModel) (*int64, error)
	GetVersion(ctx *abstraction.Context, m *model.AgingUtangPiutangFilterModel) (*model.GetVersionModel, error)
	FindByCriteria(ctx *abstraction.Context, m *model.AgingUtangPiutangFilterModel) (*model.AgingUtangPiutangEntityModel, error)
}

type agingutangpiutang struct {
	abstraction.Repository
}

func NewAgingUtangPiutang(db *gorm.DB) *agingutangpiutang {
	return &agingutangpiutang{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *agingutangpiutang) Find(ctx *abstraction.Context, m *model.AgingUtangPiutangFilterModel, p *abstraction.Pagination) (*[]model.AgingUtangPiutangEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.AgingUtangPiutangEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.AgingUtangPiutangEntityModel{})
	//filter
	// query = r.Filter(ctx, query, *m)
	tableName := model.AgingUtangPiutangEntityModel{}.TableName()
	query = r.FilterTable(ctx, query, *m, tableName)
	query = r.AllowedCompany(ctx, query, tableName)
	tmp := m.CompanyCustomFilter
	queryCompany := conn.Model(&model.CompanyEntityModel{}).Select("id")
	queryUser := conn.Model(&model.UserEntityModel{}).Select("id")
	query = r.FilterUser(ctx, query, queryUser, m.Filter, tableName)
	query = r.FilterMultiCompany(ctx, query, queryCompany, tmp, tableName)
	query = r.FilterMultiVersion(ctx, query, m.AgingUtangPiutangFilter)
	query = query.Where("aging_utang_piutang.status != 0")
	query = query.Joins(`INNER JOIN "imported_worksheet" "Import" ON "Import"."company_id" = "aging_utang_piutang"."company_id" 
	AND "Import"."period" = "aging_utang_piutang"."period" 
	AND "Import"."versions" = "aging_utang_piutang"."versions" 
	AND "Import"."status" NOT IN (0, 1)`)
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
	} else if p.SortBy != nil && *p.SortBy == "user" {
		sortBy := "\"UserCreated\".name"
		p.SortBy = &sortBy
	} else {
		sortBy := fmt.Sprintf("aging_utang_piutang.%s", *p.SortBy)
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
	query = query.Count(&totalData).Limit(limit).Offset(offset)
	// if err := query.Preload("Company").Preload("UserCreated").Preload("UserModified").Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
	if err := query.Joins("Company").Joins("UserCreated", func(db *gorm.DB) *gorm.DB {
		db = db.Select("id, name")
		return db
	}).Preload("UserModified").Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
		return &datas, &info, err
	}

	for i, v := range datas {
		datas[i].UserCreatedString = v.UserCreated.Name
		datas[i].UserModifiedString = &v.UserModified.Name
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

func (r *agingutangpiutang) FindByID(ctx *abstraction.Context, id *int) (*model.AgingUtangPiutangEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.AgingUtangPiutangEntityModel
	if err := conn.Where("id = ?", &id).Preload("Company").Preload("UserCreated").Preload("UserModified").First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	data.UserCreatedString = data.UserCreated.Name
	data.UserModifiedString = &data.UserModified.Name
	return &data, nil
}

func (r *agingutangpiutang) Create(ctx *abstraction.Context, e *model.AgingUtangPiutangEntityModel) (*model.AgingUtangPiutangEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Preload("Company").Preload("UserCreated").Preload("UserModified").First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	e.UserCreatedString = e.UserCreated.Name
	e.UserModifiedString = &e.UserModified.Name
	return e, nil
}

func (r *agingutangpiutang) Update(ctx *abstraction.Context, id *int, e *model.AgingUtangPiutangEntityModel) (*model.AgingUtangPiutangEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Where("id = ?", &id).Preload("Company").Preload("UserCreated").Preload("UserModified").First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	e.UserCreatedString = e.UserCreated.Name
	e.UserModifiedString = &e.UserModified.Name
	return e, nil
}

func (r *agingutangpiutang) Destroy(ctx *abstraction.Context, id *int, e *model.AgingUtangPiutangEntityModel) (*model.AgingUtangPiutangEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id =?", id).Delete(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil
}

func (r *agingutangpiutang) Delete(ctx *abstraction.Context, id *int, e *model.AgingUtangPiutangEntityModel) (*model.AgingUtangPiutangEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id =?", id).Update("status", 4).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil
}

func (r *agingutangpiutang) GetCount(ctx *abstraction.Context, m *model.AgingUtangPiutangFilterModel) (*int64, error) {
	var jmlData int64
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.AgingUtangPiutangEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Count(&jmlData).WithContext(ctx.Request().Context()).Error; err != nil {
		return &jmlData, err
	}

	return &jmlData, nil
}

func (r *agingutangpiutang) GetVersion(ctx *abstraction.Context, m *model.AgingUtangPiutangFilterModel) (*model.GetVersionModel, error) {
	var data []model.AgingUtangPiutangEntityModel
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.AgingUtangPiutangEntityModel{})
	query = r.Filter(ctx, query, *m)

	tmp1 := m.CompanyCustomFilter
	queryCompany := conn.Model(&model.CompanyEntityModel{}).Select("id")
	query = r.FilterMultiCompany(ctx, query, queryCompany, tmp1, "")

	query = query.Select("versions").Where("status != 0").Group("versions").Order("versions ASC")

	if err := query.Find(&data).Error; err != nil {
		return &model.GetVersionModel{}, err
	}

	var result model.GetVersionModel
	tmp := []map[string]string{}
	for _, v := range data {
		tmp1 := map[string]string{
			"value": fmt.Sprintf("%d", v.Versions),
			"label": fmt.Sprintf("Version %d", v.Versions),
		}
		tmp = append(tmp, tmp1)
	}
	result.Version = tmp
	return &result, nil
}

func (r *agingutangpiutang) FindByCriteria(ctx *abstraction.Context, m *model.AgingUtangPiutangFilterModel) (*model.AgingUtangPiutangEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data model.AgingUtangPiutangEntityModel
	query := conn.Model(&model.AgingUtangPiutangEntityModel{})
	query = r.Filter(ctx, query, *m)
	if err := query.First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
