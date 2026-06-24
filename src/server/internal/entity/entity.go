package entity

import (
	"time"

	"github.com/sylunbelievable/secondadmin/server/internal/snowflake"
	"gorm.io/gorm"
)

type User struct {
	ID                uint64 `gorm:"primaryKey;autoIncrement:false"`
	Username          string `gorm:"size:64;uniqueIndex"`
	PasswordHash      string `gorm:"size:255"`
	Nickname          string `gorm:"size:100"`
	Status            int16
	PasswordChangedAt time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (User) TableName() string { return "sys_users" }

func (u *User) BeforeCreate(*gorm.DB) error { return assignID(&u.ID) }

type Role struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement:false"`
	Code      string `gorm:"size:64;uniqueIndex"`
	Name      string `gorm:"size:100"`
	Status    int16
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Role) TableName() string { return "sys_roles" }

func (r *Role) BeforeCreate(*gorm.DB) error { return assignID(&r.ID) }

type API struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement:false"`
	Group     string `gorm:"column:group;size:100"`
	Name      string `gorm:"size:100"`
	Path      string `gorm:"size:255"`
	Method    string `gorm:"size:10"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (API) TableName() string { return "sys_apis" }

func (a *API) BeforeCreate(*gorm.DB) error { return assignID(&a.ID) }

type LoginLog struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement:false"`
	UserID    *uint64
	Username  string
	Event     string
	Success   bool
	IP        string
	UserAgent string
	DeviceID  string
	CreatedAt time.Time
}

func (LoginLog) TableName() string { return "sys_login_logs" }

func (l *LoginLog) BeforeCreate(*gorm.DB) error { return assignID(&l.ID) }

func assignID(id *uint64) error {
	if *id != 0 {
		return nil
	}
	next, err := snowflake.Next()
	if err != nil {
		return err
	}
	*id = next
	return nil
}
