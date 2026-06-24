package entity

import "time"

type User struct {
	ID                uint64 `gorm:"primaryKey"`
	Username          string `gorm:"size:64;uniqueIndex"`
	PasswordHash      string `gorm:"size:255"`
	Nickname          string `gorm:"size:100"`
	Status            int16
	PasswordChangedAt time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (User) TableName() string { return "sys_users" }

type Role struct {
	ID        uint64 `gorm:"primaryKey"`
	Code      string `gorm:"size:64;uniqueIndex"`
	Name      string `gorm:"size:100"`
	Status    int16
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Role) TableName() string { return "sys_roles" }

type API struct {
	ID        uint64 `gorm:"primaryKey"`
	Group     string `gorm:"column:group;size:100"`
	Name      string `gorm:"size:100"`
	Path      string `gorm:"size:255"`
	Method    string `gorm:"size:10"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (API) TableName() string { return "sys_apis" }

type LoginLog struct {
	ID        uint64 `gorm:"primaryKey"`
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
