package bootstrap

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/sylunbelievable/secondadmin/server/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

type DatabaseHandle struct {
	primary  *sql.DB
	resolver *dbresolver.DBResolver
}

func OpenDatabase(ctx context.Context, cfg config.Database) (*gorm.DB, *DatabaseHandle, error) {
	dialector, err := newDialector(cfg.Driver, cfg.WriterDSN())
	if err != nil {
		return nil, nil, err
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger:         logger.Default.LogMode(logger.Error),
		TranslateError: true,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("open database: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, fmt.Errorf("database handle: %w", err)
	}

	handle := &DatabaseHandle{primary: sqlDB}
	if cfg.ReadWrite.Enabled {
		readers := make([]gorm.Dialector, 0, len(cfg.ReaderDSNs()))
		for _, dsn := range cfg.ReaderDSNs() {
			dialector, err := newDialector(cfg.Driver, dsn)
			if err != nil {
				_ = handle.Close()
				return nil, nil, err
			}
			readers = append(readers, dialector)
		}
		resolver := dbresolver.Register(dbresolver.Config{
			Sources:  []gorm.Dialector{dialector},
			Replicas: readers,
			Policy:   dbresolver.RoundRobinPolicy(),
		})
		if err := db.Use(resolver); err != nil {
			_ = handle.Close()
			return nil, nil, fmt.Errorf("register database resolver: %w", err)
		}
		handle.resolver = resolver
		resolver.SetMaxOpenConns(cfg.MaxOpenConnections).
			SetMaxIdleConns(cfg.MaxIdleConnections).
			SetConnMaxLifetime(cfg.ConnectionMaxLifetime)
	} else {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConnections)
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConnections)
		sqlDB.SetConnMaxLifetime(cfg.ConnectionMaxLifetime)
	}

	if err := handle.PingContext(ctx); err != nil {
		_ = handle.Close()
		return nil, nil, fmt.Errorf("ping database: %w", err)
	}
	return db, handle, nil
}

func newDialector(driver, dsn string) (gorm.Dialector, error) {
	switch driver {
	case "postgres":
		return postgres.Open(dsn), nil
	case "mysql":
		return mysql.Open(dsn), nil
	default:
		return nil, fmt.Errorf("unsupported database driver %q", driver)
	}
}

func WriterDatabase(db *gorm.DB) *gorm.DB {
	return db.Clauses(dbresolver.Write)
}

func (h *DatabaseHandle) PingContext(ctx context.Context) error {
	return h.each(func(db *sql.DB) error { return db.PingContext(ctx) })
}

func (h *DatabaseHandle) Close() error {
	return h.each(func(db *sql.DB) error { return db.Close() })
}

func (h *DatabaseHandle) each(fn func(*sql.DB) error) error {
	if h == nil || h.primary == nil {
		return nil
	}
	if h.resolver == nil {
		return fn(h.primary)
	}
	seen := map[*sql.DB]struct{}{}
	var result error
	_ = h.resolver.Call(func(conn gorm.ConnPool) error {
		db, ok := conn.(*sql.DB)
		if !ok {
			return nil
		}
		if _, ok := seen[db]; ok {
			return nil
		}
		seen[db] = struct{}{}
		result = errors.Join(result, fn(db))
		return nil
	})
	if _, ok := seen[h.primary]; !ok {
		result = errors.Join(result, fn(h.primary))
	}
	return result
}
