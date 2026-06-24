package main

import (
	"context"
	"fmt"

	"github.com/sylunbelievable/secondadmin/server/internal/entity"
	"github.com/sylunbelievable/secondadmin/server/internal/repository"
	"gorm.io/gorm/clause"
)

type menuSeed struct {
	Key, Parent string
	Menu        entity.Menu
}

var adminMenuSeeds = []menuSeed{
	{Key: "system", Menu: entity.Menu{Type: "directory", Name: "系统管理", Sort: 10, Visible: true, Status: 1}},
	{Key: "users", Parent: "system", Menu: entity.Menu{Type: "menu", Name: "用户管理", Path: "/users", Component: "users", Sort: 10, Visible: true, Status: 1}},
	{Key: "users:create", Parent: "users", Menu: button("新建用户", "system:user:create", 10)},
	{Key: "users:update", Parent: "users", Menu: button("编辑用户", "system:user:update", 20)},
	{Key: "users:roles", Parent: "users", Menu: button("分配角色", "system:user:roles", 30)},
	{Key: "roles", Parent: "system", Menu: entity.Menu{Type: "menu", Name: "角色管理", Path: "/roles", Component: "roles", Sort: 20, Visible: true, Status: 1}},
	{Key: "roles:create", Parent: "roles", Menu: button("新建角色", "system:role:create", 10)},
	{Key: "roles:update", Parent: "roles", Menu: button("编辑角色", "system:role:update", 20)},
	{Key: "roles:delete", Parent: "roles", Menu: button("删除角色", "system:role:delete", 30)},
	{Key: "roles:apis", Parent: "roles", Menu: button("分配 API", "system:role:apis", 40)},
	{Key: "roles:menus", Parent: "roles", Menu: button("分配菜单", "system:role:menus", 50)},
	{Key: "apis", Parent: "system", Menu: entity.Menu{Type: "menu", Name: "API 管理", Path: "/apis", Component: "apis", Sort: 30, Visible: true, Status: 1}},
	{Key: "apis:create", Parent: "apis", Menu: button("新建 API", "system:api:create", 10)},
	{Key: "apis:update", Parent: "apis", Menu: button("编辑 API", "system:api:update", 20)},
	{Key: "apis:delete", Parent: "apis", Menu: button("删除 API", "system:api:delete", 30)},
	{Key: "menus", Parent: "system", Menu: entity.Menu{Type: "menu", Name: "菜单管理", Path: "/menus", Component: "menus", Sort: 40, Visible: true, Status: 1}},
	{Key: "menus:create", Parent: "menus", Menu: button("新建菜单", "system:menu:create", 10)},
	{Key: "menus:update", Parent: "menus", Menu: button("编辑菜单", "system:menu:update", 20)},
	{Key: "menus:delete", Parent: "menus", Menu: button("删除菜单", "system:menu:delete", 30)},
	{Key: "dictionaries", Parent: "system", Menu: entity.Menu{Type: "menu", Name: "数据字典", Path: "/dictionaries", Component: "dictionaries", Sort: 50, Visible: true, Status: 1}},
	{Key: "dictionaries:create", Parent: "dictionaries", Menu: button("新建字典", "system:dictionary:create", 10)},
	{Key: "dictionaries:update", Parent: "dictionaries", Menu: button("编辑字典", "system:dictionary:update", 20)},
	{Key: "dictionaries:delete", Parent: "dictionaries", Menu: button("删除字典", "system:dictionary:delete", 30)},
	{Key: "dictionaries:items", Parent: "dictionaries", Menu: button("管理字典项", "system:dictionary:item", 40)},
	{Key: "operation-logs", Parent: "system", Menu: entity.Menu{Type: "menu", Name: "操作日志", Path: "/operation-logs", Component: "operation-logs", Sort: 60, Visible: true, Status: 1}},
	{Key: "login-logs", Parent: "system", Menu: entity.Menu{Type: "menu", Name: "登录日志", Path: "/login-logs", Component: "login-logs", Sort: 70, Visible: true, Status: 1}},
	{Key: "sessions", Parent: "system", Menu: entity.Menu{Type: "menu", Name: "在线设备", Path: "/sessions", Component: "sessions", Sort: 80, Visible: true, Status: 1}},
	{Key: "sessions:delete", Parent: "sessions", Menu: button("下线设备", "system:session:delete", 10)},
}

func button(name, permission string, sort int) entity.Menu {
	return entity.Menu{Type: "button", Name: name, Sort: sort, Permission: &permission, Status: 1}
}

func ensureAdminMenus(ctx context.Context, repos *repository.Repositories, roleID uint64) error {
	if err := validateAdminMenuSeeds(); err != nil {
		return err
	}
	return repos.WithinTransaction(ctx, func(tx *repository.Repositories) error {
		ids := make(map[string]uint64, len(adminMenuSeeds))
		for _, seed := range adminMenuSeeds {
			item := seed.Menu
			item.ParentID = ids[seed.Parent]
			values := map[string]any{
				"type": item.Type, "path": item.Path, "component": item.Component, "icon": item.Icon,
				"sort": item.Sort, "visible": item.Visible, "permission": item.Permission, "status": item.Status,
			}
			if err := tx.DB.WithContext(ctx).Where("parent_id = ? AND name = ?", item.ParentID, item.Name).
				Assign(values).FirstOrCreate(&item).Error; err != nil {
				return err
			}
			ids[seed.Key] = item.ID
			if err := tx.DB.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).
				Create(&entity.RoleMenu{RoleID: roleID, MenuID: item.ID}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func validateAdminMenuSeeds() error {
	keys := map[string]bool{"": true}
	for _, seed := range adminMenuSeeds {
		if seed.Key == "" || keys[seed.Key] || !keys[seed.Parent] {
			return fmt.Errorf("invalid admin menu seed %q", seed.Key)
		}
		keys[seed.Key] = true
	}
	return nil
}
