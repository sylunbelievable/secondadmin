package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/casbin/casbin/v3"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/sylunbelievable/secondadmin/server/internal/bootstrap"
	"github.com/sylunbelievable/secondadmin/server/internal/config"
	"github.com/sylunbelievable/secondadmin/server/internal/entity"
	"github.com/sylunbelievable/secondadmin/server/internal/repository"
	"github.com/sylunbelievable/secondadmin/server/internal/service"
)

func run() error {
	username := strings.ToLower(strings.TrimSpace(os.Getenv("ADMIN_USERNAME")))
	password := os.Getenv("ADMIN_PASSWORD")
	if username == "" || password == "" {
		return errors.New("ADMIN_USERNAME and ADMIN_PASSWORD are required")
	}
	hash, err := service.HashPassword(password)
	if err != nil {
		return err
	}
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	db, sqlDB, err := bootstrap.OpenDatabase(context.Background(), cfg.Database)
	if err != nil {
		return err
	}
	defer sqlDB.Close()
	writeDB := bootstrap.WriterDatabase(db)
	repos := repository.New(writeDB)

	var user entity.User
	var role entity.Role
	err = repos.WithinTransaction(context.Background(), func(tx *repository.Repositories) error {
		user, err = tx.UserByUsername(context.Background(), username)
		if errors.Is(err, repository.ErrNotFound) {
			user = entity.User{Username: username, PasswordHash: hash, Nickname: "Administrator", Status: 1, PasswordChangedAt: time.Now()}
			err = tx.CreateUser(context.Background(), &user)
		}
		if err != nil {
			return err
		}
		role, err = tx.RoleByCode(context.Background(), "admin")
		if errors.Is(err, repository.ErrNotFound) {
			role = entity.Role{Code: "admin", Name: "Administrator", Status: 1}
			err = tx.CreateRole(context.Background(), &role)
		}
		return err
	})
	if err != nil {
		return err
	}
	if err = ensureAdminMenus(context.Background(), repos, role.ID); err != nil {
		return err
	}

	gormadapter.TurnOffAutoMigrate(writeDB)
	adapter, err := gormadapter.NewAdapterByDB(writeDB)
	if err != nil {
		return err
	}
	enforcer, err := casbin.NewEnforcer("configs/casbin-model.conf", adapter)
	if err != nil {
		return err
	}
	subject := strconv.FormatUint(user.ID, 10)
	if _, err = enforcer.AddRoleForUser(subject, role.Code); err != nil {
		return err
	}
	if _, err = enforcer.AddPolicy(role.Code, "/api/v1/*", ".*"); err != nil {
		return err
	}
	fmt.Printf("administrator ready: %s\n", username)
	return nil
}
