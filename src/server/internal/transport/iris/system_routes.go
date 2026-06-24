package iris

import (
	"github.com/kataras/iris/v12"
	"github.com/sylunbelievable/secondadmin/server/internal/dto"
	"github.com/sylunbelievable/secondadmin/server/internal/service"
)

func registerSystemRoutes(admin iris.Party, s *service.Services) {
	admin.Get("/menus", listMenus(s))
	admin.Post("/menus", createMenu(s))
	admin.Put("/menus/{id:uint64}", updateMenu(s))
	admin.Delete("/menus/{id:uint64}", deleteMenu(s))
	admin.Put("/roles/{id:uint64}/menus", setRoleMenus(s))

	admin.Get("/dictionaries", listDictionaries(s))
	admin.Post("/dictionaries", createDictionary(s))
	admin.Put("/dictionaries/{id:uint64}", updateDictionary(s))
	admin.Delete("/dictionaries/{id:uint64}", deleteDictionary(s))
	admin.Post("/dictionaries/{id:uint64}/items", createDictionaryItem(s))
	admin.Put("/dictionary-items/{id:uint64}", updateDictionaryItem(s))
	admin.Delete("/dictionary-items/{id:uint64}", deleteDictionaryItem(s))
	admin.Get("/operation-logs", operationLogs(s))
	admin.Get("/login-logs", loginLogs(s))
}

func listMenus(s *service.Services) iris.Handler {
	return func(ctx iris.Context) {
		items, err := s.Menus(ctx.Request().Context())
		if err != nil {
			writeError(ctx, err)
			return
		}
		_ = ctx.JSON(items)
	}
}

func currentMenus(s *service.Services) iris.Handler {
	return func(ctx iris.Context) {
		item, err := s.CurrentMenus(ctx.Request().Context(), principal(ctx).UserID)
		if err != nil {
			writeError(ctx, err)
			return
		}
		_ = ctx.JSON(item)
	}
}

func createMenu(s *service.Services) iris.Handler {
	return func(ctx iris.Context) {
		var req dto.MenuRequest
		if ctx.ReadJSON(&req) != nil {
			writeError(ctx, service.ErrInvalidInput)
			return
		}
		item, err := s.CreateMenu(ctx.Request().Context(), req)
		if err != nil {
			writeError(ctx, err)
			return
		}
		ctx.StatusCode(iris.StatusCreated)
		_ = ctx.JSON(item)
	}
}

func updateMenu(s *service.Services) iris.Handler {
	return func(ctx iris.Context) {
		var req dto.MenuRequest
		if ctx.ReadJSON(&req) != nil {
			writeError(ctx, service.ErrInvalidInput)
			return
		}
		if err := s.UpdateMenu(ctx.Request().Context(), ctx.Params().GetUint64Default("id", 0), req); err != nil {
			writeError(ctx, err)
			return
		}
		ctx.StatusCode(iris.StatusNoContent)
	}
}

func deleteMenu(s *service.Services) iris.Handler {
	return func(ctx iris.Context) {
		if err := s.DeleteMenu(ctx.Request().Context(), ctx.Params().GetUint64Default("id", 0)); err != nil {
			writeError(ctx, err)
			return
		}
		ctx.StatusCode(iris.StatusNoContent)
	}
}

func setRoleMenus(s *service.Services) iris.Handler {
	return func(ctx iris.Context) {
		var req dto.IDsRequest
		if ctx.ReadJSON(&req) != nil {
			writeError(ctx, service.ErrInvalidInput)
			return
		}
		if err := s.SetRoleMenus(ctx.Request().Context(), ctx.Params().GetUint64Default("id", 0), req.IDs); err != nil {
			writeError(ctx, err)
			return
		}
		ctx.StatusCode(iris.StatusNoContent)
	}
}

func listDictionaries(s *service.Services) iris.Handler {
	return func(ctx iris.Context) {
		items, err := s.Dictionaries(ctx.Request().Context())
		if err != nil {
			writeError(ctx, err)
			return
		}
		_ = ctx.JSON(items)
	}
}

func createDictionary(s *service.Services) iris.Handler {
	return func(ctx iris.Context) {
		var req dto.DictionaryRequest
		if ctx.ReadJSON(&req) != nil {
			writeError(ctx, service.ErrInvalidInput)
			return
		}
		item, err := s.CreateDictionary(ctx.Request().Context(), req)
		if err != nil {
			writeError(ctx, err)
			return
		}
		ctx.StatusCode(iris.StatusCreated)
		_ = ctx.JSON(item)
	}
}

func updateDictionary(s *service.Services) iris.Handler {
	return func(ctx iris.Context) {
		var req dto.DictionaryRequest
		if ctx.ReadJSON(&req) != nil {
			writeError(ctx, service.ErrInvalidInput)
			return
		}
		if err := s.UpdateDictionary(ctx.Request().Context(), ctx.Params().GetUint64Default("id", 0), req); err != nil {
			writeError(ctx, err)
			return
		}
		ctx.StatusCode(iris.StatusNoContent)
	}
}

func deleteDictionary(s *service.Services) iris.Handler {
	return func(ctx iris.Context) {
		if err := s.DeleteDictionary(ctx.Request().Context(), ctx.Params().GetUint64Default("id", 0)); err != nil {
			writeError(ctx, err)
			return
		}
		ctx.StatusCode(iris.StatusNoContent)
	}
}

func createDictionaryItem(s *service.Services) iris.Handler {
	return func(ctx iris.Context) {
		var req dto.DictionaryItemRequest
		if ctx.ReadJSON(&req) != nil {
			writeError(ctx, service.ErrInvalidInput)
			return
		}
		item, err := s.CreateDictionaryItem(ctx.Request().Context(), ctx.Params().GetUint64Default("id", 0), req)
		if err != nil {
			writeError(ctx, err)
			return
		}
		ctx.StatusCode(iris.StatusCreated)
		_ = ctx.JSON(item)
	}
}

func updateDictionaryItem(s *service.Services) iris.Handler {
	return func(ctx iris.Context) {
		var req dto.DictionaryItemRequest
		if ctx.ReadJSON(&req) != nil {
			writeError(ctx, service.ErrInvalidInput)
			return
		}
		if err := s.UpdateDictionaryItem(ctx.Request().Context(), ctx.Params().GetUint64Default("id", 0), req); err != nil {
			writeError(ctx, err)
			return
		}
		ctx.StatusCode(iris.StatusNoContent)
	}
}

func deleteDictionaryItem(s *service.Services) iris.Handler {
	return func(ctx iris.Context) {
		if err := s.DeleteDictionaryItem(ctx.Request().Context(), ctx.Params().GetUint64Default("id", 0)); err != nil {
			writeError(ctx, err)
			return
		}
		ctx.StatusCode(iris.StatusNoContent)
	}
}

func dictionaryItems(s *service.Services) iris.Handler {
	return func(ctx iris.Context) {
		items, err := s.DictionaryItems(ctx.Request().Context(), ctx.Params().Get("code"))
		if err != nil {
			writeError(ctx, err)
			return
		}
		_ = ctx.JSON(items)
	}
}

func operationLogs(s *service.Services) iris.Handler {
	return func(ctx iris.Context) {
		page, size := positive(ctx.URLParamIntDefault("page", 1), 1), positive(ctx.URLParamIntDefault("pageSize", 20), 20)
		if size > 100 {
			size = 100
		}
		items, total, err := s.OperationLogs(ctx.Request().Context(), page, size)
		if err != nil {
			writeError(ctx, err)
			return
		}
		_ = ctx.JSON(iris.Map{"items": items, "total": total})
	}
}

func loginLogs(s *service.Services) iris.Handler {
	return func(ctx iris.Context) {
		page, size := positive(ctx.URLParamIntDefault("page", 1), 1), positive(ctx.URLParamIntDefault("pageSize", 20), 20)
		if size > 100 {
			size = 100
		}
		items, total, err := s.LoginLogs(ctx.Request().Context(), page, size)
		if err != nil {
			writeError(ctx, err)
			return
		}
		_ = ctx.JSON(iris.Map{"items": items, "total": total})
	}
}
