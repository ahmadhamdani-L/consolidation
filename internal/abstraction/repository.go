package abstraction

import (
	"fmt"
	"mcash-finance-console-core/pkg/redis"
	"reflect"
	"strings"
	"time"

	"gorm.io/gorm"
)

type IRepository interface {
	CheckTrx(ctx *Context) *gorm.DB
	Filter(ctx *Context, query *gorm.DB, payload interface{}) *gorm.DB
}

type Repository struct {
	Connection *gorm.DB
	Db         *gorm.DB
	Tx         *gorm.DB
}

func (r *Repository) CheckTrx(ctx *Context) *gorm.DB {
	if ctx.Trx != nil {
		return ctx.Trx.Db
	}
	return r.Db
}

func (r *Repository) Filter(ctx *Context, query *gorm.DB, payload interface{}) *gorm.DB {
	mVal := reflect.ValueOf(payload)
	mType := reflect.TypeOf(payload)

	for i := 0; i < mVal.NumField(); i++ {
		mValChild := mVal.Field(i)
		mTypeChild := mType.Field(i)

		for j := 0; j < mValChild.NumField(); j++ {
			val := mValChild.Field(j)

			if !val.IsNil() {
				if val.Kind() == reflect.Ptr {
					val = mValChild.Field(j).Elem()
				}

				key := mTypeChild.Type.Field(j).Tag.Get("query")
				filter := mTypeChild.Type.Field(j).Tag.Get("filter")
				// TODO need to custom field
				// filterColumn := mTypeChild.Type.Field(j).Tag.Get("filterColumn")

				switch filter {
				case "LIKE":
					query = query.Where(fmt.Sprintf("%s LIKE ?", key), "%"+val.String()+"%")
				case "ILIKE":
					query = query.Where(fmt.Sprintf("%s ILIKE ?", key), "%"+val.String()+"%")
				case "DATE":
					tmpDate, err := time.Parse("2006-01-02", val.String())
					if err != nil {
						continue
					}
					tmpStr := tmpDate.Format("2006-01-02")
					query = query.Where(fmt.Sprintf("DATE(%s) = ?", key), tmpStr)
				case "DATESTRING":
					query = query.Where(fmt.Sprintf("DATE(%s) = ?", key), val.String())

				case "CUSTOM":
					continue
				default:
					query = query.Where(fmt.Sprintf("%s = ?", key), val.Interface())
				}
			}
		}
	}

	return query
}

func (r *Repository) FilterTable(ctx *Context, query *gorm.DB, payload interface{}, tableName string) *gorm.DB {
	mVal := reflect.ValueOf(payload)
	mType := reflect.TypeOf(payload)

	for i := 0; i < mVal.NumField(); i++ {
		mValChild := mVal.Field(i)
		mTypeChild := mType.Field(i)

		for j := 0; j < mValChild.NumField(); j++ {
			val := mValChild.Field(j)

			if !val.IsNil() {
				if val.Kind() == reflect.Ptr {
					val = mValChild.Field(j).Elem()
				}

				key := mTypeChild.Type.Field(j).Tag.Get("query")
				filter := mTypeChild.Type.Field(j).Tag.Get("filter")
				// TODO need to custom field
				// filterColumn := mTypeChild.Type.Field(j).Tag.Get("filterColumn")

				switch filter {
				case "LIKE":
					query = query.Where(fmt.Sprintf("%s.%s LIKE ?", tableName, key), "%"+val.String()+"%")
				case "ILIKE":
					query = query.Where(fmt.Sprintf("%s.%s ILIKE ?", tableName, key), "%"+val.String()+"%")
				case "DATE":
					tmpDate, err := time.Parse("2006-01-02", val.String())
					if err != nil {
						continue
					}
					tmpStr := tmpDate.Format("2006-01-02")
					query = query.Where(fmt.Sprintf("DATE(%s.%s) = ?", tableName, key), tmpStr)
				case "DATESTRING":
					query = query.Where(fmt.Sprintf("DATE(%s.%s) = ?", tableName, key), val.String())

				case "CUSTOM":
					continue
				default:
					query = query.Where(fmt.Sprintf("%s.%s = ?", tableName, key), val.Interface())
				}
			}
		}
	}

	return query
}

func (r *Repository) FilterMultiCompany(ctx *Context, query *gorm.DB, queryCompany *gorm.DB, payload interface{}, tableName string) *gorm.DB {
	mVal := reflect.ValueOf(&payload).Elem()
	if !mVal.Elem().FieldByName("ArrCompanyString").IsNil() {
		company := mVal.Elem().FieldByName("ArrCompanyString").Elem().Interface().([]string)
		operator := mVal.Elem().FieldByName("ArrCompanyOperator").Elem().Interface().([]string)
		for i, v := range company {
			switch strings.ToUpper(operator[i]) {
			case "not_contain":
				queryCompany = queryCompany.Or("name NOT ILIKE ?", "%"+v+"%")
			case "exact":
				queryCompany = queryCompany.Or("name = ?", "%"+v+"%")
			case "not_exact":
				queryCompany = queryCompany.Or("name != ?", "%"+v+"%")
			default:
				queryCompany = queryCompany.Or("name ILIKE ?", "%"+v+"%")
			}
		}
		if tableName != "" {
			query = query.Where(fmt.Sprintf("%s.company_id IN (?)", tableName), queryCompany)
		} else {
			query = query.Where("company_id IN (?)", queryCompany)
		}
	}

	return query
}

func (r *Repository) FilterMultiVersion(ctx *Context, query *gorm.DB, payload interface{}) *gorm.DB {
	mVal := reflect.ValueOf(&payload).Elem()
	if !mVal.Elem().FieldByName("ArrVersions").IsNil() {
		versions := mVal.Elem().FieldByName("ArrVersions").Elem().Interface().([]int)
		if len(versions) > 0 {
			listVersions := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(versions)), ","), "[]")
			query = query.Where(fmt.Sprintf("trial_balance.versions IN (%s)", listVersions))
		}
	}
	return query
}

func (r *Repository) FilterUser(ctx *Context, query, queryUser *gorm.DB, payload interface{}, tableName string) *gorm.DB {
	mVal := reflect.ValueOf(&payload).Elem()
	if !mVal.Elem().FieldByName("UserCreatedString").IsNil() {
		userCreatedString := mVal.Elem().FieldByName("UserCreatedString").Elem().Interface().(string)
		if userCreatedString != "" {
			if tableName != "" {
				query = query.Where(fmt.Sprintf("%s.created_by IN (?)", tableName), queryUser.Where(fmt.Sprintf("name ILIKE ?"), "%"+userCreatedString+"%"))
			} else {
				query = query.Where("created_by IN (?)", queryUser.Where(fmt.Sprintf("name ILIKE ?"), "%"+userCreatedString+"%"))
			}
		}
	}
	return query
}

func (r *Repository) FilterMultiStatus(ctx *Context, query *gorm.DB, payload interface{}) *gorm.DB {
	mVal := reflect.ValueOf(&payload).Elem()
	if !mVal.Elem().FieldByName("ArrStatus").IsNil() {
		status := mVal.Elem().FieldByName("ArrStatus").Elem().Interface().([]int)
		if len(status) > 0 {
			listStatus := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(status)), ","), "[]")
			query = query.Where(fmt.Sprintf("status IN (%s)", listStatus))
		}
	}
	return query
}

func (r *Repository) AllowedCompany(ctx *Context, query *gorm.DB, tableName string) *gorm.DB {
	allCompany, err := redis.RedisClient.Get(ctx.Request().Context(), fmt.Sprintf("access_all_company:user:%d", ctx.Auth.ID)).Result()
	// mau is nil nya true atau false tapi kalau dia errornya is nil dia ignore
	if redis.IsNil(err) || !redis.IsNil(err) {

	} else if err != nil {
		return query
	}

	if allCompany == "1" {
		return query
	}

	CacheCompany, err := redis.RedisClient.LRange(ctx.Request().Context(), fmt.Sprintf("access_company:user:%d", ctx.Auth.ID), 0, -1).Result()
	if err != nil && redis.IsNil(err) == true {
		return query
	}

	if len(CacheCompany) > 0 {
		listCompanyAllowed := strings.Trim(strings.Join(CacheCompany, ","), "[]")
		query = query.Where(fmt.Sprintf("%s.company_id IN (%s)", tableName, listCompanyAllowed))
	}
	return query
}
