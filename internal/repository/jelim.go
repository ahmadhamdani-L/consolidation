package repository

import (
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"

	"gorm.io/gorm"
)

type Jelim interface {
	Find(ctx *abstraction.Context, m *model.JelimFilterModel, p *abstraction.Pagination) (*[]model.JelimEntityModel, *abstraction.PaginationInfo, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.JelimEntityModel, error)
	Create(ctx *abstraction.Context, e *model.JelimEntityModel) (*model.JelimEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.JelimEntityModel) (*model.JelimEntityModel, error)
	Delete(ctx *abstraction.Context, id *int, e *model.JelimEntityModel) (*model.JelimEntityModel, error)
	Get(ctx *abstraction.Context, m *model.JelimFilterModel) (*[]model.JelimEntityModel, error)
	GetCount(ctx *abstraction.Context, m *model.JelimFilterModel) (*int64, error)
	GetVersion(ctx *abstraction.Context, m *model.JelimFilterModel) (*model.GetVersionModel, error)
	UpdateTbd(ctx *abstraction.Context, id *int, e *model.ConsolidationDetailEntityModel) (*model.ConsolidationDetailEntityModel, error)
	FindByTbd(ctx *abstraction.Context, id *int, c *string) (*model.ConsolidationDetailEntityModel, error)
	FindByFormatter(ctx *abstraction.Context, id *int) (*model.FormatterBridgesEntityModel, error)
	Export(ctx *abstraction.Context, id *int) (*model.JelimEntityModel, error)
	FindByCoa(ctx *abstraction.Context, c *string) (*model.CoaEntityModel, error)
	FindByCoaType(ctx *abstraction.Context, idType *int) (*model.CoaTypeEntityModel, error)
	FindByCoaGroup(ctx *abstraction.Context, id *int) (*model.CoaGroupEntityModel, error)
	FindByConsoleBridge(ctx *abstraction.Context, id *int) (*[]model.ConsolidationBridgeEntityModel, error)
	FindByConsoleBridgeDetail(ctx *abstraction.Context, id *int, code *string) (*model.ConsolidationBridgeDetailEntityModel, error)
	FindByCoaConsoleDetail(ctx *abstraction.Context, id int, code *string) (*[]model.ConsolidationDetailEntityModel, error)
}

type jelim struct {
	abstraction.Repository
}

func NewJelim(db *gorm.DB) *jelim {
	return &jelim{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *jelim) FindByCoaConsoleDetail(ctx *abstraction.Context, id int, code *string) (*[]model.ConsolidationDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.ConsolidationDetailEntityModel
	if err := conn.Where("consolidation_id = ? AND code ILIKE ?", id, *code+"%").Find(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *jelim) FindByConsoleBridge(ctx *abstraction.Context, id *int) (*[]model.ConsolidationBridgeEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.ConsolidationBridgeEntityModel
	if err := conn.Where("consolidation_id = ?", &id).First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *jelim) FindByConsoleBridgeDetail(ctx *abstraction.Context, id *int, code *string) (*model.ConsolidationBridgeDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.ConsolidationBridgeDetailEntityModel
	if err := conn.Where("consolidation_bridge_id = ? AND code = ?", &id, &code).First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *jelim) FindByCoaGroup(ctx *abstraction.Context, id *int) (*model.CoaGroupEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.CoaGroupEntityModel
	if err := conn.Where("id = ?", &id).First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *jelim) FindByCoa(ctx *abstraction.Context, c *string) (*model.CoaEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.CoaEntityModel
	if err := conn.Where("code = ? ", &c).First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *jelim) FindByCoaType(ctx *abstraction.Context, idType *int) (*model.CoaTypeEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.CoaTypeEntityModel
	if err := conn.Where("id = ?", &idType).First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *jelim) Find(ctx *abstraction.Context, m *model.JelimFilterModel, p *abstraction.Pagination) (*[]model.JelimEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.JelimEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.JelimEntityModel{})
	//filter
	tableName := model.JelimEntityModel{}.TableName()
	query = r.FilterTable(ctx, query, *m, tableName)
	query = r.AllowedCompany(ctx, query, tableName)

	//filter custom
	tmp1 := m.CompanyCustomFilter
	queryCompany := conn.Model(&model.CompanyEntityModel{}).Select("id")
	query = r.FilterMultiCompany(ctx, query, queryCompany, tmp1, "")
	query = r.FilterMultiVersion(ctx, query, m.JelimFilter)
	queryUser := conn.Model(&model.UserEntityModel{}).Select("id")
	query = r.FilterUser(ctx, query, queryUser, m.Filter, "")

	if m.Search != nil {
		query = query.Where("created_by IN (?)", conn.Model(&model.UserEntityModel{}).Select("id").Where("name ILIKE ?", "%"+*m.Search+"%"))
	}

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
	}
	if p.SortBy != nil && (tmpSortBy != nil && *tmpSortBy != "company" && *tmpSortBy != "user") {
		sortBy := fmt.Sprintf("jelim.%s", *p.SortBy)
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

	if err := query.Preload("Consolidation").Preload("JelimDetail").Joins("Company").Joins("UserCreated", func(db *gorm.DB) *gorm.DB {
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
	if int(totalData) > *p.PageSize {
		info.MoreRecords = true
		// info.Count -= 1
		// datas = datas[:len(datas)-1]
	}

	return &datas, &info, nil
}

func (r *jelim) FindByID(ctx *abstraction.Context, id *int) (*model.JelimEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.JelimEntityModel
	if err := conn.Where("id = ?", &id).Preload("Company").Preload("Consolidation").Preload("JelimDetail").Preload("UserCreated").Preload("UserModified").First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	data.UserCreatedString = data.UserCreated.Name
	data.UserModifiedString = &data.UserModified.Name
	return &data, nil
}

func (r *jelim) Create(ctx *abstraction.Context, e *model.JelimEntityModel) (*model.JelimEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Preload("UserCreated").Preload("UserModified").First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	e.UserCreatedString = e.UserCreated.Name
	e.UserModifiedString = &e.UserModified.Name
	return e, nil
}

func (r *jelim) Update(ctx *abstraction.Context, id *int, e *model.JelimEntityModel) (*model.JelimEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).Preload("Company").WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Where("id = ?", &id).Preload("Company").Preload("UserCreated").Preload("UserModified").First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	e.UserCreatedString = e.UserCreated.Name
	e.UserModifiedString = &e.UserModified.Name
	return e, nil

}
func (r *jelim) UpdateTbd(ctx *abstraction.Context, id *int, e *model.ConsolidationDetailEntityModel) (*model.ConsolidationDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil

}

func (r *jelim) FindByTbd(ctx *abstraction.Context, id *int, c *string) (*model.ConsolidationDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.ConsolidationDetailEntityModel
	if err := conn.Where("consolidation_id = ? and code = ?", &id, &c).First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *jelim) FindByFormatter(ctx *abstraction.Context, id *int) (*model.FormatterBridgesEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.FormatterBridgesEntityModel
	if err := conn.Where("trx_ref_id = ? ", &id).First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *jelim) Destroy(ctx *abstraction.Context, id *int, e *model.JelimEntityModel) (*model.JelimEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id =?", id).Delete(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil
}

func (r *jelim) Delete(ctx *abstraction.Context, id *int, e *model.JelimEntityModel) (*model.JelimEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id =?", id).Update("status", 4).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	e.UserCreatedString = e.UserCreated.Name
	e.UserModifiedString = &e.UserModified.Name
	return e, nil
}

func (r *jelim) Get(ctx *abstraction.Context, m *model.JelimFilterModel) (*[]model.JelimEntityModel, error) {
	var datas []model.JelimEntityModel

	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.JelimEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
		return &datas, err
	}

	// var formatterBridges []model.FormatterBridgesEntityModel

	// query = conn.Model(&model.FormatterBridgesEntityModel{}).Where("trx_ref_id = ?", datas.ID).Where("source = ?", "TRIAL-BALANCE")
	// if err := query.Find(&formatterBridges).Error; err != nil {
	// 	return &datas, err
	// }
	// datas.FormatterBridges = formatterBridges

	return &datas, nil
}

func (r *jelim) GetCount(ctx *abstraction.Context, m *model.JelimFilterModel) (*int64, error) {
	var jmlData int64
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.JelimEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Count(&jmlData).WithContext(ctx.Request().Context()).Error; err != nil {
		return &jmlData, err
	}

	return &jmlData, nil
}

func (r *jelim) GetVersion(ctx *abstraction.Context, m *model.JelimFilterModel) (*model.GetVersionModel, error) {
	var data []model.JelimEntityModel
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.JelimEntityModel{})
	query = r.Filter(ctx, query, *m)

	//filter custom
	tmp1 := m.CompanyCustomFilter
	queryCompany := conn.Model(&model.CompanyEntityModel{}).Select("id")
	query = r.FilterMultiCompany(ctx, query, queryCompany, tmp1, "")

	query = query.Select("consolidation_id").Group("consolidation_id").Order("consolidation_id ASC")

	if err := query.Find(&data).Error; err != nil {
		return &model.GetVersionModel{}, err
	}

	var result model.GetVersionModel
	tmp := []map[string]string{}
	for _, v := range data {
		tmp1 := map[string]string{
			"value": fmt.Sprintf("%d", v.ConsolidationID),
			"label": fmt.Sprintf("Version %d", v.ConsolidationID),
		}
		tmp = append(tmp, tmp1)
	}
	result.Version = tmp
	return &result, nil
}

func (r *jelim) Export(ctx *abstraction.Context, id *int) (*model.JelimEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data model.JelimEntityModel
	query := conn.Model(&model.JelimEntityModel{}).Preload("Company").Preload("JelimDetail").Where("id = ?", &id).Find(&data)
	if err := query.Error; err != nil {
		return nil, err
	}
	return &data, nil
}
