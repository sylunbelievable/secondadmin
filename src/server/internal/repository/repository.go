package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/sylunbelievable/secondadmin/server/internal/entity"
	"gorm.io/gorm"
)

var (
	ErrNotFound = errors.New("not found")
	ErrConflict = errors.New("conflict")
)

type Repositories struct{ DB *gorm.DB }

func New(db *gorm.DB) *Repositories { return &Repositories{DB: db} }

func (r *Repositories) WithinTransaction(ctx context.Context, fn func(*Repositories) error) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error { return fn(New(tx)) })
}

func (r *Repositories) UserByUsername(ctx context.Context, username string) (entity.User, error) {
	var user entity.User
	err := r.DB.WithContext(ctx).Where("username = ?", strings.ToLower(strings.TrimSpace(username))).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = ErrNotFound
	}
	return user, err
}

func (r *Repositories) UserByID(ctx context.Context, id uint64) (entity.User, error) {
	var user entity.User
	err := r.DB.WithContext(ctx).First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = ErrNotFound
	}
	return user, err
}

func (r *Repositories) CreateUser(ctx context.Context, user *entity.User) error {
	user.Username = strings.ToLower(strings.TrimSpace(user.Username))
	return dbError(r.DB.WithContext(ctx).Create(user).Error)
}

func (r *Repositories) UpdateUser(ctx context.Context, id uint64, values map[string]any) error {
	result := r.DB.WithContext(ctx).Model(&entity.User{}).Where("id = ?", id).Updates(values)
	if result.Error == nil && result.RowsAffected == 0 {
		return ErrNotFound
	}
	return dbError(result.Error)
}

func (r *Repositories) ListUsers(ctx context.Context, page, size int) ([]entity.User, int64, error) {
	var users []entity.User
	var total int64
	db := r.DB.WithContext(ctx).Model(&entity.User{})
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return users, total, db.Order("id").Offset((page - 1) * size).Limit(size).Find(&users).Error
}

func (r *Repositories) CreateRole(ctx context.Context, role *entity.Role) error {
	role.Code = strings.ToLower(strings.TrimSpace(role.Code))
	return dbError(r.DB.WithContext(ctx).Create(role).Error)
}

func (r *Repositories) RoleByID(ctx context.Context, id uint64) (entity.Role, error) {
	var role entity.Role
	err := r.DB.WithContext(ctx).First(&role, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = ErrNotFound
	}
	return role, err
}

func (r *Repositories) RoleByCode(ctx context.Context, code string) (entity.Role, error) {
	var role entity.Role
	err := r.DB.WithContext(ctx).Where("code = ?", strings.ToLower(strings.TrimSpace(code))).First(&role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = ErrNotFound
	}
	return role, err
}

func (r *Repositories) ListRoles(ctx context.Context) ([]entity.Role, error) {
	var roles []entity.Role
	return roles, r.DB.WithContext(ctx).Order("id").Find(&roles).Error
}

func (r *Repositories) DeleteRole(ctx context.Context, id uint64) error {
	result := r.DB.WithContext(ctx).Delete(&entity.Role{}, id)
	if result.Error == nil && result.RowsAffected == 0 {
		return ErrNotFound
	}
	return dbError(result.Error)
}

func (r *Repositories) UpdateRole(ctx context.Context, id uint64, values map[string]any) error {
	result := r.DB.WithContext(ctx).Model(&entity.Role{}).Where("id = ?", id).Updates(values)
	if result.Error == nil && result.RowsAffected == 0 {
		return ErrNotFound
	}
	return result.Error
}

func (r *Repositories) CreateAPI(ctx context.Context, api *entity.API) error {
	api.Method = strings.ToUpper(strings.TrimSpace(api.Method))
	return dbError(r.DB.WithContext(ctx).Create(api).Error)
}

func (r *Repositories) APIByID(ctx context.Context, id uint64) (entity.API, error) {
	var api entity.API
	err := r.DB.WithContext(ctx).First(&api, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = ErrNotFound
	}
	return api, err
}

func (r *Repositories) ListAPIs(ctx context.Context) ([]entity.API, error) {
	var apis []entity.API
	return apis, r.DB.WithContext(ctx).Order("id").Find(&apis).Error
}

func (r *Repositories) UpdateAPI(ctx context.Context, id uint64, values map[string]any) error {
	result := r.DB.WithContext(ctx).Model(&entity.API{}).Where("id = ?", id).Updates(values)
	if result.Error == nil && result.RowsAffected == 0 {
		return ErrNotFound
	}
	return dbError(result.Error)
}

func (r *Repositories) DeleteAPI(ctx context.Context, id uint64) error {
	result := r.DB.WithContext(ctx).Delete(&entity.API{}, id)
	if result.Error == nil && result.RowsAffected == 0 {
		return ErrNotFound
	}
	return result.Error
}

func (r *Repositories) AddLoginLog(ctx context.Context, item *entity.LoginLog) error {
	return r.DB.WithContext(ctx).Create(item).Error
}

func (r *Repositories) LoginLogs(ctx context.Context, page, size int) ([]entity.LoginLog, int64, error) {
	var items []entity.LoginLog
	var total int64
	db := r.DB.WithContext(ctx).Model(&entity.LoginLog{})
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return items, total, db.Order("id DESC").Offset((page - 1) * size).Limit(size).Find(&items).Error
}

func dbError(err error) error {
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return ErrConflict
	}
	return err
}
