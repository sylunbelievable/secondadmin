package entity

import "time"

type Menu struct {
	ID         uint64 `gorm:"primaryKey"`
	ParentID   uint64
	Type       string
	Name       string
	Path       string
	Component  string
	Icon       string
	Sort       int
	Visible    bool
	Permission *string
	Status     int16
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (Menu) TableName() string { return "sys_menus" }

type RoleMenu struct {
	RoleID uint64 `gorm:"primaryKey"`
	MenuID uint64 `gorm:"primaryKey"`
}

func (RoleMenu) TableName() string { return "sys_role_menus" }

type Dictionary struct {
	ID        uint64 `gorm:"primaryKey"`
	Code      string
	Name      string
	Status    int16
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Dictionary) TableName() string { return "sys_dictionaries" }

type DictionaryItem struct {
	ID           uint64 `gorm:"primaryKey"`
	DictionaryID uint64
	Label        string
	Value        string
	Sort         int
	Status       int16
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (DictionaryItem) TableName() string { return "sys_dictionary_items" }

type OperationLog struct {
	ID         uint64 `gorm:"primaryKey"`
	UserID     uint64
	RequestID  string
	Method     string
	Path       string
	StatusCode int
	DurationMS int64
	IP         string
	UserAgent  string
	CreatedAt  time.Time
}

func (OperationLog) TableName() string { return "sys_operation_logs" }
