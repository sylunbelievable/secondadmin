package service

import (
	"context"
	"slices"
	"strconv"
	"strings"

	"github.com/sylunbelievable/secondadmin/server/internal/dto"
	"github.com/sylunbelievable/secondadmin/server/internal/entity"
)

const superAdminRole = "admin"

func (s *Services) Menus(ctx context.Context) ([]dto.Menu, error) {
	items, err := s.Repos.Menus(ctx)
	return buildMenuTree(items), err
}

func (s *Services) CreateMenu(ctx context.Context, req dto.MenuRequest) (dto.Menu, error) {
	item, err := s.validateMenu(ctx, 0, req)
	if err != nil {
		return dto.Menu{}, err
	}
	if err := s.Repos.CreateMenu(ctx, &item); err != nil {
		return dto.Menu{}, err
	}
	return menuDTO(item), nil
}

func (s *Services) UpdateMenu(ctx context.Context, id uint64, req dto.MenuRequest) error {
	item, err := s.validateMenu(ctx, id, req)
	if err != nil {
		return err
	}
	return s.Repos.UpdateMenu(ctx, id, map[string]any{
		"parent_id": item.ParentID, "type": item.Type, "name": item.Name, "path": item.Path,
		"component": item.Component, "icon": item.Icon, "sort": item.Sort, "visible": item.Visible,
		"permission": item.Permission, "status": item.Status,
	})
}

func (s *Services) DeleteMenu(ctx context.Context, id uint64) error {
	return s.Repos.DeleteMenu(ctx, id)
}

func (s *Services) SetRoleMenus(ctx context.Context, roleID uint64, menuIDs []uint64) error {
	if _, err := s.Repos.RoleByID(ctx, roleID); err != nil {
		return err
	}
	expanded := make(map[uint64]struct{}, len(menuIDs))
	for _, id := range menuIDs {
		for id != 0 {
			if _, exists := expanded[id]; exists {
				break
			}
			menu, err := s.Repos.MenuByID(ctx, id)
			if err != nil {
				return ErrInvalidInput
			}
			expanded[id] = struct{}{}
			id = menu.ParentID
		}
	}
	menuIDs = menuIDs[:0]
	for id := range expanded {
		menuIDs = append(menuIDs, id)
	}
	slices.Sort(menuIDs)
	return s.Repos.ReplaceRoleMenus(ctx, roleID, menuIDs)
}

func (s *Services) CurrentMenus(ctx context.Context, userID uint64) (dto.CurrentMenus, error) {
	roles, err := s.Casbin.GetRolesForUser(strconv.FormatUint(userID, 10))
	if err != nil {
		return dto.CurrentMenus{}, err
	}
	var items []entity.Menu
	if slices.Contains(roles, superAdminRole) {
		items, err = s.Repos.Menus(ctx)
	} else {
		items, err = s.Repos.MenusByRoleCodes(ctx, roles)
	}
	if err != nil {
		return dto.CurrentMenus{}, err
	}
	byID := make(map[uint64]entity.Menu, len(items))
	for _, item := range items {
		byID[item.ID] = item
	}
	var routes []entity.Menu
	permissions := make([]string, 0)
	for _, item := range items {
		if item.Type == "button" {
			if item.Permission != nil && menuReachable(item, byID, false) {
				permissions = append(permissions, *item.Permission)
			}
			continue
		}
		if menuReachable(item, byID, true) {
			routes = append(routes, item)
		}
	}
	slices.Sort(permissions)
	return dto.CurrentMenus{Menus: pruneEmptyDirectories(buildMenuTree(routes)), Permissions: slices.Compact(permissions)}, nil
}

func menuReachable(item entity.Menu, byID map[uint64]entity.Menu, requireVisible bool) bool {
	seen := map[uint64]bool{}
	for {
		if seen[item.ID] || item.Status != 1 || requireVisible && !item.Visible {
			return false
		}
		seen[item.ID] = true
		if item.ParentID == 0 {
			return true
		}
		parent, ok := byID[item.ParentID]
		if !ok {
			return false
		}
		item = parent
	}
}

func pruneEmptyDirectories(items []dto.Menu) []dto.Menu {
	result := make([]dto.Menu, 0, len(items))
	for _, item := range items {
		item.Children = pruneEmptyDirectories(item.Children)
		if item.Type != "directory" || len(item.Children) > 0 {
			result = append(result, item)
		}
	}
	return result
}

func (s *Services) validateMenu(ctx context.Context, id uint64, req dto.MenuRequest) (entity.Menu, error) {
	if !slices.Contains([]string{"directory", "menu", "button"}, req.Type) || strings.TrimSpace(req.Name) == "" {
		return entity.Menu{}, ErrInvalidInput
	}
	if req.ParentID != 0 {
		parent, err := s.Repos.MenuByID(ctx, req.ParentID)
		if err != nil || parent.Type == "button" {
			return entity.Menu{}, ErrInvalidInput
		}
		if id != 0 {
			cyclic, err := s.Repos.MenuHasDescendant(ctx, id, req.ParentID)
			if err != nil || cyclic {
				return entity.Menu{}, ErrInvalidInput
			}
		}
	}
	visible, status := true, int16(1)
	if req.Visible != nil {
		visible = *req.Visible
	}
	if req.Status != nil {
		status = *req.Status
	}
	permission := req.Permission
	if permission != nil {
		value := strings.TrimSpace(*permission)
		if value == "" {
			permission = nil
		} else {
			permission = &value
		}
	}
	if req.Type == "button" {
		req.Path, req.Component, req.Icon, visible = "", "", "", false
		if permission == nil {
			return entity.Menu{}, ErrInvalidInput
		}
	}
	return entity.Menu{
		ParentID: req.ParentID, Type: req.Type, Name: strings.TrimSpace(req.Name), Path: req.Path,
		Component: req.Component, Icon: req.Icon, Sort: req.Sort, Visible: visible,
		Permission: permission, Status: status,
	}, nil
}

func buildMenuTree(items []entity.Menu) []dto.Menu {
	children := make(map[uint64][]entity.Menu)
	for _, item := range items {
		children[item.ParentID] = append(children[item.ParentID], item)
	}
	var build func(uint64) []dto.Menu
	build = func(parent uint64) []dto.Menu {
		result := make([]dto.Menu, 0, len(children[parent]))
		for _, item := range children[parent] {
			node := menuDTO(item)
			node.Children = build(item.ID)
			result = append(result, node)
		}
		return result
	}
	return build(0)
}

func (s *Services) Dictionaries(ctx context.Context) ([]dto.Dictionary, error) {
	items, err := s.Repos.Dictionaries(ctx)
	out := make([]dto.Dictionary, len(items))
	for i := range items {
		out[i] = dictionaryDTO(items[i])
	}
	return out, err
}

func (s *Services) CreateDictionary(ctx context.Context, req dto.DictionaryRequest) (dto.Dictionary, error) {
	if req.Code == "" || req.Name == "" {
		return dto.Dictionary{}, ErrInvalidInput
	}
	item := entity.Dictionary{Code: req.Code, Name: strings.TrimSpace(req.Name), Status: 1}
	if err := s.Repos.CreateDictionary(ctx, &item); err != nil {
		return dto.Dictionary{}, err
	}
	return dictionaryDTO(item), nil
}

func (s *Services) UpdateDictionary(ctx context.Context, id uint64, req dto.DictionaryRequest) error {
	values := map[string]any{"name": strings.TrimSpace(req.Name)}
	if req.Code != "" {
		values["code"] = strings.ToLower(strings.TrimSpace(req.Code))
	}
	if req.Status != nil {
		values["status"] = *req.Status
	}
	return s.Repos.UpdateDictionary(ctx, id, values)
}

func (s *Services) DeleteDictionary(ctx context.Context, id uint64) error {
	return s.Repos.DeleteDictionary(ctx, id)
}

func (s *Services) CreateDictionaryItem(ctx context.Context, dictionaryID uint64, req dto.DictionaryItemRequest) (dto.DictionaryItem, error) {
	if _, err := s.Repos.DictionaryByID(ctx, dictionaryID); err != nil || req.Label == "" || req.Value == "" {
		return dto.DictionaryItem{}, ErrInvalidInput
	}
	item := entity.DictionaryItem{DictionaryID: dictionaryID, Label: strings.TrimSpace(req.Label), Value: req.Value, Sort: req.Sort, Status: 1}
	if err := s.Repos.CreateDictionaryItem(ctx, &item); err != nil {
		return dto.DictionaryItem{}, err
	}
	return dictionaryItemDTO(item), nil
}

func (s *Services) DictionaryItems(ctx context.Context, code string) ([]dto.DictionaryItem, error) {
	items, err := s.Repos.EnabledDictionaryItems(ctx, code)
	out := make([]dto.DictionaryItem, len(items))
	for i := range items {
		out[i] = dictionaryItemDTO(items[i])
	}
	return out, err
}

func (s *Services) UpdateDictionaryItem(ctx context.Context, id uint64, req dto.DictionaryItemRequest) error {
	if req.Label == "" || req.Value == "" {
		return ErrInvalidInput
	}
	values := map[string]any{"label": strings.TrimSpace(req.Label), "value": req.Value, "sort": req.Sort}
	if req.Status != nil {
		values["status"] = *req.Status
	}
	return s.Repos.UpdateDictionaryItem(ctx, id, values)
}

func (s *Services) DeleteDictionaryItem(ctx context.Context, id uint64) error {
	return s.Repos.DeleteDictionaryItem(ctx, id)
}

func (s *Services) AddOperationLog(ctx context.Context, item *entity.OperationLog) error {
	return s.Repos.AddOperationLog(ctx, item)
}

func (s *Services) OperationLogs(ctx context.Context, page, size int) ([]dto.OperationLog, int64, error) {
	items, total, err := s.Repos.OperationLogs(ctx, page, size)
	out := make([]dto.OperationLog, len(items))
	for i := range items {
		item := items[i]
		out[i] = dto.OperationLog{
			ID: item.ID, UserID: item.UserID, RequestID: item.RequestID, Method: item.Method,
			Path: item.Path, StatusCode: item.StatusCode, DurationMS: item.DurationMS,
			IP: item.IP, UserAgent: item.UserAgent, CreatedAt: item.CreatedAt,
		}
	}
	return out, total, err
}

func (s *Services) LoginLogs(ctx context.Context, page, size int) ([]dto.LoginLog, int64, error) {
	items, total, err := s.Repos.LoginLogs(ctx, page, size)
	out := make([]dto.LoginLog, len(items))
	for i := range items {
		item := items[i]
		out[i] = dto.LoginLog{
			ID: item.ID, UserID: item.UserID, Username: item.Username, Event: item.Event, Success: item.Success,
			IP: item.IP, UserAgent: item.UserAgent, DeviceID: item.DeviceID, CreatedAt: item.CreatedAt,
		}
	}
	return out, total, err
}

func menuDTO(item entity.Menu) dto.Menu {
	return dto.Menu{
		ID: item.ID, ParentID: item.ParentID, Type: item.Type, Name: item.Name, Path: item.Path,
		Component: item.Component, Icon: item.Icon, Sort: item.Sort, Visible: item.Visible,
		Permission: item.Permission, Status: item.Status,
	}
}

func dictionaryDTO(item entity.Dictionary) dto.Dictionary {
	return dto.Dictionary{ID: item.ID, Code: item.Code, Name: item.Name, Status: item.Status}
}

func dictionaryItemDTO(item entity.DictionaryItem) dto.DictionaryItem {
	return dto.DictionaryItem{
		ID: item.ID, DictionaryID: item.DictionaryID, Label: item.Label,
		Value: item.Value, Sort: item.Sort, Status: item.Status,
	}
}
