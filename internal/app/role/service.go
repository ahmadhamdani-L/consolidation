package role

import (
	"errors"
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/dto"
	"mcash-finance-console-core/internal/factory"
	"mcash-finance-console-core/internal/model"
	"mcash-finance-console-core/internal/repository"
	"mcash-finance-console-core/pkg/redis"
	"mcash-finance-console-core/pkg/util/helper"
	"mcash-finance-console-core/pkg/util/response"
	"mcash-finance-console-core/pkg/util/trxmanager"
	"strings"

	"gorm.io/gorm"
)

type service struct {
	Repository        repository.Role
	UserRepository    repository.User
	PermissionDef     repository.PermissionDef
	RolePermission    repository.RolePermission
	RolePermissionApi repository.RolePermissionApi
	Db                *gorm.DB
}

type Service interface {
	Find(ctx *abstraction.Context, payload *dto.RoleGetRequest) (*dto.RoleGetResponse, error)
	FindByID(ctx *abstraction.Context, payload *dto.RoleGetByIDRequest) (*dto.RoleGetByIDResponse, error)
	Create(ctx *abstraction.Context, payload *dto.RoleCreateRequest) (*dto.RoleCreateResponse, error)
	Update(ctx *abstraction.Context, payload *dto.RoleUpdateRequest) (*dto.RoleUpdateResponse, error)
	Delete(ctx *abstraction.Context, payload *dto.RoleDeleteRequest) (*dto.RoleDeleteResponse, error)
}

func NewService(f *factory.Factory) *service {
	repository := f.RoleRepository
	userRepo := f.UserRepository
	rolePermissionRepo := f.RolePermissionRepository
	rolePermissionApiRepo := f.RolePermissionApiRepository
	permissionDefRepo := f.PermissionDefRepository
	db := f.Db
	return &service{
		Repository:        repository,
		Db:                db,
		UserRepository:    userRepo,
		PermissionDef:     permissionDefRepo,
		RolePermission:    rolePermissionRepo,
		RolePermissionApi: rolePermissionApiRepo,
	}
}

func (s *service) Find(ctx *abstraction.Context, payload *dto.RoleGetRequest) (*dto.RoleGetResponse, error) {
	data, info, err := s.Repository.Find(ctx, &payload.RoleFilterModel, &payload.Pagination)
	if err != nil {
		return &dto.RoleGetResponse{}, helper.ErrorHandler(err)
	}

	result := &dto.RoleGetResponse{
		Datas:          *data,
		PaginationInfo: *info,
	}
	return result, nil
}

func (s *service) FindByID(ctx *abstraction.Context, payload *dto.RoleGetByIDRequest) (*dto.RoleGetByIDResponse, error) {
	data, err := s.Repository.FindByID(ctx, &payload.ID)
	if err != nil {
		return &dto.RoleGetByIDResponse{}, helper.ErrorHandler(err)
	}
	result := &dto.RoleGetByIDResponse{
		RoleEntityModel: *data,
	}
	return result, nil
}

func (s *service) Create(ctx *abstraction.Context, payload *dto.RoleCreateRequest) (*dto.RoleCreateResponse, error) {
	var data model.RoleEntityModel

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		data.Context = ctx
		payload.Code = strings.ReplaceAll(strings.ToUpper(payload.Name), " ", "-")
		data.RoleEntity = payload.RoleEntity

		result, err := s.Repository.Create(ctx, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		tmpFunctionalID := make(map[string]bool)
		for _, v := range payload.RolePermission {
			if tmpFunctionalID[v.FunctionalID] {
				continue
			}
			tmpFunctionalID[v.FunctionalID] = true
			rolePermissionData := model.RolePermissionEntityModel{}
			rolePermissionData.RoleID = result.ID
			tmpFalse := false
			rolePermissionData.Create = &tmpFalse
			rolePermissionData.Read = &tmpFalse
			rolePermissionData.Update = &tmpFalse
			rolePermissionData.Delete = &tmpFalse

			permissionDef, err := s.PermissionDef.FindByFunctionalID(ctx, v.FunctionalID)
			if err != nil {
				return helper.ErrorHandler(err)
			}

			rolePermissionData.FunctionalID = v.FunctionalID
			rolePermissionData.Create = &tmpFalse
			rolePermissionData.Read = &tmpFalse
			rolePermissionData.Update = &tmpFalse
			rolePermissionData.Delete = &tmpFalse
			if v.Create != nil && *v.Create && permissionDef.AllowCreate != nil && *permissionDef.AllowCreate {
				rolePermissionData.Create = v.Create
			}
			if v.Read != nil && *v.Read && permissionDef.AllowRead != nil && *permissionDef.AllowRead {
				rolePermissionData.Read = v.Read
			}
			if v.Update != nil && *v.Update && permissionDef.AllowUpdate != nil && *permissionDef.AllowUpdate {
				rolePermissionData.Update = v.Update
			}
			if v.Delete != nil && *v.Delete && permissionDef.AllowDelete != nil && *permissionDef.AllowDelete {
				rolePermissionData.Delete = v.Delete
			}

			_, err = s.RolePermission.Create(ctx, &rolePermissionData)
			if err != nil {
				return err
			}
			canAct := []string{}
			if permissionDef.AllowCreate != nil && *permissionDef.AllowCreate && v.Create != nil && *v.Create {
				canAct = append(canAct, "POST")
			}
			if permissionDef.AllowRead != nil && *permissionDef.AllowRead && v.Read != nil && *v.Read {
				canAct = append(canAct, "GET")
			}
			if permissionDef.AllowUpdate != nil && *permissionDef.AllowUpdate && v.Update != nil && *v.Update {
				canAct = append(canAct, "PATCH")
			}
			if permissionDef.AllowDelete != nil && *permissionDef.AllowDelete && v.Delete != nil && *v.Delete {
				canAct = append(canAct, "DELETE")
			}

			for _, v := range canAct {
				rolePermissionApiData := model.RolePermissionApiEntityModel{}
				rolePermissionApiData.RoleID = result.ID
				rolePermissionApiData.ApiMethod = v
				switch v {
				case "GET":
					rolePermissionApiData.ApiPath = permissionDef.ApiPathRead
				case "POST":
					rolePermissionApiData.ApiPath = permissionDef.ApiPathCreate
				case "PATCH":
					rolePermissionApiData.ApiPath = permissionDef.ApiPathUpdate
				case "DELETE":
					rolePermissionApiData.ApiPath = permissionDef.ApiPathDelete
				}
				rolePermissionApi, err := s.RolePermissionApi.Create(ctx, &rolePermissionApiData)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				err = redis.RedisClient.Set(ctx.Request().Context(), fmt.Sprintf("role=%d:allow=%s_%s", rolePermissionApi.RoleID, rolePermissionApi.ApiMethod, rolePermissionApi.ApiPath), true, 0).Err()
				if err != nil {
					panic(err)
				}
			}
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.RoleCreateResponse{}, err
	}

	result := &dto.RoleCreateResponse{
		RoleEntityModel: data,
	}
	return result, nil
}

func (s *service) Updates(ctx *abstraction.Context, payload *dto.RoleUpdateRequest) (*dto.RoleUpdateResponse, error) { // tidak dipakai. update by data role permission, not recreating data
	var data model.RoleEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		role, err := s.Repository.FindByID(ctx, &payload.ID)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data.Context = ctx
		data.RoleEntity = payload.RoleEntity

		result, err := s.Repository.Update(ctx, &role.ID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		// err = s.Repository.DeleteRelatedRole(ctx, &result.ID)
		// if err != nil {
		// 	return err
		// }

		// CacheRolePA, err := redis.RedisClient.Keys(ctx.Request().Context(), fmt.Sprintf("role=%d*", role.ID)).Result()
		// if err != nil && redis.IsNil(err) != true {
		// 	return err
		// }

		// for _, keyRolePA := range CacheRolePA {
		// 	err = redis.RedisClient.Del(ctx.Request().Context(), keyRolePA).Err()
		// 	if err != nil {
		// 		return err
		// 	}
		// }

		tmpFunctionalID := make(map[string]bool)
		tmpFalse := false
		for _, v := range payload.RolePermission {
			if tmpFunctionalID[v.FunctionalID] {
				continue
			}
			tmpFunctionalID[v.FunctionalID] = true
			rolePermissionData := model.RolePermissionEntityModel{}
			rolePermissionData.RoleID = role.ID

			rolePermissionApiData := model.RolePermissionApiEntityModel{}
			rolePermissionApiData.RoleID = role.ID

			permissionDef, err := s.PermissionDef.FindByFunctionalID(ctx, v.FunctionalID)
			if err != nil {
				return helper.ErrorHandler(err)
			}

			rolePermission, err := s.RolePermission.FindByFunctionalID(ctx, &permissionDef.FunctionalID)
			if err != nil && err.Error() != "record not found" {
				return err
			}
			rolePermissionData.FunctionalID = permissionDef.FunctionalID
			rolePermissionData.Create = &tmpFalse
			rolePermissionData.Read = &tmpFalse
			rolePermissionData.Update = &tmpFalse
			rolePermissionData.Delete = &tmpFalse
			if v.Create != nil && *v.Create && permissionDef.AllowCreate != nil && *permissionDef.AllowCreate {
				rolePermissionData.Create = v.Create
			}
			if v.Read != nil && *v.Read && permissionDef.AllowRead != nil && *permissionDef.AllowRead {
				rolePermissionData.Read = v.Read
			}
			if v.Update != nil && *v.Update && permissionDef.AllowUpdate != nil && *permissionDef.AllowUpdate {
				rolePermissionData.Update = v.Update
			}
			if v.Delete != nil && *v.Delete && permissionDef.AllowDelete != nil && *permissionDef.AllowDelete {
				rolePermissionData.Delete = v.Delete
			}

			if v.Create != nil && !*v.Create && v.Read != nil && !*v.Read && v.Update != nil && !*v.Update && v.Delete != nil && !*v.Delete {
				_, err = s.RolePermission.DeleteByRoleFunctionalID(ctx, &rolePermissionData)
				if err != nil {
					return helper.ErrorHandler(err)
				}
			} else {
				if rolePermission.ID == 0 {
					_, err = s.RolePermission.Create(ctx, &rolePermissionData)
					if err != nil {
						return helper.ErrorHandler(err)
					}
				} else {
					_, err = s.RolePermission.UpdateByRoleFunctionalID(ctx, &role.ID, &permissionDef.FunctionalID, &rolePermissionData)
					if err != nil {
						return err
					}
				}
			}

			canAct := []string{"POST", "GET", "PATCH", "DELETE"}

			for _, vAct := range canAct {
				rolePermissionApiData := model.RolePermissionApiEntityModel{}
				rolePermissionApiData.RoleID = role.ID
				rolePermissionApiData.ApiMethod = vAct
				switch vAct {
				case "GET":
					rolePermissionApiData.ApiPath = permissionDef.ApiPathRead
				case "POST":
					rolePermissionApiData.ApiPath = permissionDef.ApiPathCreate
				case "PATCH":
					rolePermissionApiData.ApiPath = permissionDef.ApiPathUpdate
				case "DELETE":
					rolePermissionApiData.ApiPath = permissionDef.ApiPathDelete
				}

				if (permissionDef.AllowCreate != nil && *permissionDef.AllowCreate && v.Create != nil && *v.Create) || (permissionDef.AllowRead != nil && *permissionDef.AllowRead && v.Read != nil && *v.Read) || (permissionDef.AllowUpdate != nil && *permissionDef.AllowUpdate && v.Update != nil && *v.Update) || (permissionDef.AllowDelete != nil && *permissionDef.AllowDelete && v.Delete != nil && *v.Delete) {
					rolePermissionApi, err := s.RolePermissionApi.FirstOrCreate(ctx, &rolePermissionApiData)
					if err != nil {
						return helper.ErrorHandler(err)
					}

					err = redis.RedisClient.Set(ctx.Request().Context(), fmt.Sprintf("role=%d:functional_id=%s:allow=%s_%s", rolePermissionApi.RoleID, rolePermissionData.FunctionalID, rolePermissionApi.ApiMethod, rolePermissionApi.ApiPath), true, 0).Err()
					if err != nil {
						panic(err)
					}
				} else {
					_, err := s.RolePermissionApi.DeleteByCriteria(ctx, &rolePermissionApiData)
					if err != nil {
						return helper.ErrorHandler(err)
					}

					err = redis.RedisClient.Del(ctx.Request().Context(), fmt.Sprintf("role=%d:functional_id=%s:allow=%s_%s", role.ID, v.FunctionalID, rolePermissionApiData.ApiMethod, rolePermissionApiData.ApiPath)).Err()
					if err != nil {
						return err
					}
				}
			}
		}

		data = *result
		return nil
	}); err != nil {
		return &dto.RoleUpdateResponse{}, err
	}
	result := &dto.RoleUpdateResponse{
		RoleEntityModel: data,
	}
	return result, nil
}

func (s *service) Update(ctx *abstraction.Context, payload *dto.RoleUpdateRequest) (*dto.RoleUpdateResponse, error) { // update by re creating data role permission
	var data model.RoleEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		role, err := s.Repository.FindByID(ctx, &payload.ID)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data.Context = ctx
		// payload.Code = strings.ReplaceAll(strings.ToUpper(payload.Name), " ", "-")
		payload.Code = role.Code
		data.RoleEntity = payload.RoleEntity

		result, err := s.Repository.Update(ctx, &role.ID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		err = s.Repository.DeleteRelatedRole(ctx, &result.ID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		CacheRolePA, err := redis.RedisClient.Keys(ctx.Request().Context(), fmt.Sprintf("role=%d:*", role.ID)).Result()
		if err != nil && redis.IsNil(err) != true {
			return err
		}

		for _, keyRolePA := range CacheRolePA {
			err = redis.RedisClient.Del(ctx.Request().Context(), keyRolePA).Err()
			if err != nil {
				return err
			}
		}

		tmpFunctionalID := make(map[string]bool)
		for _, v := range payload.RolePermission {
			if tmpFunctionalID[v.FunctionalID] {
				continue
			}
			tmpFunctionalID[v.FunctionalID] = true
			rolePermissionData := model.RolePermissionEntityModel{}
			rolePermissionData.RoleID = role.ID

			rolePermissionApiData := model.RolePermissionApiEntityModel{}
			rolePermissionApiData.RoleID = role.ID

			permissionDef, err := s.PermissionDef.FindByFunctionalID(ctx, v.FunctionalID)
			if err != nil {
				return helper.ErrorHandler(err)
			}
			tmpFalse := false
			rolePermissionData.FunctionalID = permissionDef.FunctionalID
			rolePermissionData.Create = &tmpFalse
			rolePermissionData.Read = &tmpFalse
			rolePermissionData.Update = &tmpFalse
			rolePermissionData.Delete = &tmpFalse

			canAct := []string{}
			if permissionDef.AllowCreate != nil && *permissionDef.AllowCreate && v.Create != nil && *v.Create {
				rolePermissionData.Create = v.Create
				canAct = append(canAct, "POST")
			}
			if permissionDef.AllowRead != nil && *permissionDef.AllowRead && v.Read != nil && *v.Read {
				rolePermissionData.Read = v.Read
				canAct = append(canAct, "GET")
			}
			if permissionDef.AllowUpdate != nil && *permissionDef.AllowUpdate && v.Update != nil && *v.Update {
				rolePermissionData.Update = v.Update
				canAct = append(canAct, "PATCH")
			}
			if permissionDef.AllowDelete != nil && *permissionDef.AllowDelete && v.Delete != nil && *v.Delete {
				canAct = append(canAct, "DELETE")
				rolePermissionData.Delete = v.Delete
			}

			_, err = s.RolePermission.Create(ctx, &rolePermissionData)
			if err != nil {
				return helper.ErrorHandler(err)
			}

			for _, vAct := range canAct {
				rolePermissionApiData := model.RolePermissionApiEntityModel{}
				rolePermissionApiData.RoleID = role.ID
				rolePermissionApiData.ApiMethod = vAct
				switch vAct {
				case "GET":
					rolePermissionApiData.ApiPath = permissionDef.ApiPathRead
				case "POST":
					rolePermissionApiData.ApiPath = permissionDef.ApiPathCreate
				case "PATCH":
					rolePermissionApiData.ApiPath = permissionDef.ApiPathUpdate
				case "DELETE":
					rolePermissionApiData.ApiPath = permissionDef.ApiPathDelete
				}

				rolePermissionApi, err := s.RolePermissionApi.Create(ctx, &rolePermissionApiData)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				err = redis.RedisClient.Set(ctx.Request().Context(), fmt.Sprintf("role=%d:allow=%s_%s", rolePermissionApi.RoleID, rolePermissionApi.ApiMethod, rolePermissionApi.ApiPath), true, 0).Err()
				if err != nil {
					return err
				}
			}
		}

		data = *result
		return nil
	}); err != nil {
		return &dto.RoleUpdateResponse{}, err
	}
	result := &dto.RoleUpdateResponse{
		RoleEntityModel: data,
	}
	return result, nil
}

func (s *service) Delete(ctx *abstraction.Context, payload *dto.RoleDeleteRequest) (*dto.RoleDeleteResponse, error) {
	var data model.RoleEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		role, err := s.Repository.FindByID(ctx, &payload.ID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		criteriaUserHasRole := model.UserFilterModel{}
		criteriaUserHasRole.RoleID = &role.ID

		_, totalData, err := s.UserRepository.Find(ctx, &criteriaUserHasRole, &abstraction.Pagination{})
		if err != nil {
			return helper.ErrorHandler(err)
		}

		if totalData.Count > 0 {
			return response.ErrorBuilder(&response.ErrorConstant.HasRelatedData, errors.New("role has related data"))
		}

		data.Context = ctx
		result, err := s.Repository.Delete(ctx, &payload.ID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		CacheRolePA, err := redis.RedisClient.Keys(ctx.Request().Context(), fmt.Sprintf("role=%d:*", role.ID)).Result()
		if err != nil && redis.IsNil(err) != true {
			return err
		}

		for _, keyRolePA := range CacheRolePA {
			err = redis.RedisClient.Del(ctx.Request().Context(), keyRolePA).Err()
			if err != nil {
				return err
			}
		}

		data = *result
		return nil
	}); err != nil {
		return &dto.RoleDeleteResponse{}, err
	}
	result := &dto.RoleDeleteResponse{
		RoleEntityModel: data,
	}
	return result, nil
}

func (s *service) DeleteByPermission(ctx *abstraction.Context, payload *dto.RoleDeletePermissionRequest) (*dto.RoleDeletePermissionResponse, error) {
	var data model.RolePermissionEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		role, err := s.Repository.FindByID(ctx, &payload.ID)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data.FunctionalID = payload.FunctionalID
		data.Context = ctx

		permissionDef, err := s.PermissionDef.FindByFunctionalID(ctx, data.FunctionalID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		for _, vAct := range []string{"GET", "POST", "PATCH", "DELETE"} {
			rolePermissionApiData := model.RolePermissionApiEntityModel{}
			rolePermissionApiData.RoleID = role.ID
			rolePermissionApiData.ApiMethod = vAct
			switch vAct {
			case "GET":
				rolePermissionApiData.ApiPath = permissionDef.ApiPathRead
			case "POST":
				rolePermissionApiData.ApiPath = permissionDef.ApiPathCreate
			case "PATCH":
				rolePermissionApiData.ApiPath = permissionDef.ApiPathUpdate
			case "DELETE":
				rolePermissionApiData.ApiPath = permissionDef.ApiPathDelete
			}

			_, err := s.RolePermissionApi.DeleteByCriteria(ctx, &rolePermissionApiData)
			if err != nil {
				return helper.ErrorHandler(err)
			}
		}

		result, err := s.RolePermission.DeleteByRoleFunctionalID(ctx, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		CacheRolePA, err := redis.RedisClient.Keys(ctx.Request().Context(), fmt.Sprintf("role=%d:allow=*", role.ID)).Result()
		if err != nil && redis.IsNil(err) != true {
			return err
		}

		for _, keyRolePA := range CacheRolePA {
			err = redis.RedisClient.Del(ctx.Request().Context(), keyRolePA).Err()
			if err != nil {
				return err
			}
		}

		data = *result
		return nil
	}); err != nil {
		return nil, err
	}
	result := &dto.RoleDeletePermissionResponse{
		RolePermissionEntityModel: data,
	}
	return result, nil
}
