package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/sylunbelievable/secondadmin/server/internal/entity"
	"gorm.io/gorm"
)

func (r *Repositories) Menus(ctx context.Context) ([]entity.Menu, error) {
	var items []entity.Menu
	return items, r.DB.WithContext(ctx).Order("sort, id").Find(&items).Error
}

func (r *Repositories) MenuByID(ctx context.Context, id uint64) (entity.Menu, error) {
	var item entity.Menu
	err := r.DB.WithContext(ctx).First(&item, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = ErrNotFound
	}
	return item, err
}

func (r *Repositories) CreateMenu(ctx context.Context, item *entity.Menu) error {
	return dbError(r.DB.WithContext(ctx).Create(item).Error)
}

func (r *Repositories) UpdateMenu(ctx context.Context, id uint64, values map[string]any) error {
	result := r.DB.WithContext(ctx).Model(&entity.Menu{}).Where("id = ?", id).Updates(values)
	if result.Error == nil && result.RowsAffected == 0 {
		return ErrNotFound
	}
	return dbError(result.Error)
}

func (r *Repositories) DeleteMenu(ctx context.Context, id uint64) error {
	var count int64
	if err := r.DB.WithContext(ctx).Model(&entity.Menu{}).Where("parent_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return ErrConflict
	}
	if err := r.DB.WithContext(ctx).Model(&entity.RoleMenu{}).Where("menu_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return ErrConflict
	}
	result := r.DB.WithContext(ctx).Delete(&entity.Menu{}, id)
	if result.Error == nil && result.RowsAffected == 0 {
		return ErrNotFound
	}
	return result.Error
}

func (r *Repositories) MenuHasDescendant(ctx context.Context, id, candidateParent uint64) (bool, error) {
	for candidateParent != 0 {
		if candidateParent == id {
			return true, nil
		}
		parent, err := r.MenuByID(ctx, candidateParent)
		if err != nil {
			return false, err
		}
		candidateParent = parent.ParentID
	}
	return false, nil
}

func (r *Repositories) ReplaceRoleMenus(ctx context.Context, roleID uint64, menuIDs []uint64) error {
	return r.WithinTransaction(ctx, func(tx *Repositories) error {
		if err := tx.DB.WithContext(ctx).Where("role_id = ?", roleID).Delete(&entity.RoleMenu{}).Error; err != nil {
			return err
		}
		for _, menuID := range menuIDs {
			if err := tx.DB.WithContext(ctx).Create(&entity.RoleMenu{RoleID: roleID, MenuID: menuID}).Error; err != nil {
				return dbError(err)
			}
		}
		return nil
	})
}

func (r *Repositories) MenusByRoleCodes(ctx context.Context, codes []string) ([]entity.Menu, error) {
	var items []entity.Menu
	if len(codes) == 0 {
		return items, nil
	}
	err := r.DB.WithContext(ctx).Distinct("sys_menus.*").
		Joins("JOIN sys_role_menus ON sys_role_menus.menu_id = sys_menus.id").
		Joins("JOIN sys_roles ON sys_roles.id = sys_role_menus.role_id").
		Where("sys_roles.code IN ? AND sys_roles.status = 1 AND sys_menus.status = 1", codes).
		Order("sys_menus.sort, sys_menus.id").Find(&items).Error
	return items, err
}

func (r *Repositories) Dictionaries(ctx context.Context) ([]entity.Dictionary, error) {
	var items []entity.Dictionary
	return items, r.DB.WithContext(ctx).Order("id").Find(&items).Error
}

func (r *Repositories) DictionaryByID(ctx context.Context, id uint64) (entity.Dictionary, error) {
	var item entity.Dictionary
	err := r.DB.WithContext(ctx).First(&item, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = ErrNotFound
	}
	return item, err
}

func (r *Repositories) CreateDictionary(ctx context.Context, item *entity.Dictionary) error {
	item.Code = strings.ToLower(strings.TrimSpace(item.Code))
	return dbError(r.DB.WithContext(ctx).Create(item).Error)
}

func (r *Repositories) UpdateDictionary(ctx context.Context, id uint64, values map[string]any) error {
	result := r.DB.WithContext(ctx).Model(&entity.Dictionary{}).Where("id = ?", id).Updates(values)
	if result.Error == nil && result.RowsAffected == 0 {
		return ErrNotFound
	}
	return dbError(result.Error)
}

func (r *Repositories) DeleteDictionary(ctx context.Context, id uint64) error {
	var count int64
	if err := r.DB.WithContext(ctx).Model(&entity.DictionaryItem{}).Where("dictionary_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return ErrConflict
	}
	result := r.DB.WithContext(ctx).Delete(&entity.Dictionary{}, id)
	if result.Error == nil && result.RowsAffected == 0 {
		return ErrNotFound
	}
	return result.Error
}

func (r *Repositories) CreateDictionaryItem(ctx context.Context, item *entity.DictionaryItem) error {
	return dbError(r.DB.WithContext(ctx).Create(item).Error)
}

func (r *Repositories) UpdateDictionaryItem(ctx context.Context, id uint64, values map[string]any) error {
	result := r.DB.WithContext(ctx).Model(&entity.DictionaryItem{}).Where("id = ?", id).Updates(values)
	if result.Error == nil && result.RowsAffected == 0 {
		return ErrNotFound
	}
	return dbError(result.Error)
}

func (r *Repositories) DeleteDictionaryItem(ctx context.Context, id uint64) error {
	result := r.DB.WithContext(ctx).Delete(&entity.DictionaryItem{}, id)
	if result.Error == nil && result.RowsAffected == 0 {
		return ErrNotFound
	}
	return result.Error
}

func (r *Repositories) EnabledDictionaryItems(ctx context.Context, code string) ([]entity.DictionaryItem, error) {
	var items []entity.DictionaryItem
	err := r.DB.WithContext(ctx).Model(&entity.DictionaryItem{}).
		Joins("JOIN sys_dictionaries ON sys_dictionaries.id = sys_dictionary_items.dictionary_id").
		Where("sys_dictionaries.code = ? AND sys_dictionaries.status = 1 AND sys_dictionary_items.status = 1", strings.ToLower(code)).
		Order("sys_dictionary_items.sort, sys_dictionary_items.id").Find(&items).Error
	return items, err
}

func (r *Repositories) AddOperationLog(ctx context.Context, item *entity.OperationLog) error {
	return r.DB.WithContext(ctx).Create(item).Error
}

func (r *Repositories) OperationLogs(ctx context.Context, page, size int) ([]entity.OperationLog, int64, error) {
	var items []entity.OperationLog
	var total int64
	db := r.DB.WithContext(ctx).Model(&entity.OperationLog{})
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return items, total, db.Order("id DESC").Offset((page - 1) * size).Limit(size).Find(&items).Error
}
