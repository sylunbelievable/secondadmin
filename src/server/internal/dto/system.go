package dto

import "time"

type Menu struct {
	ID         uint64  `json:"id,string"`
	ParentID   uint64  `json:"parentId,string"`
	Type       string  `json:"type"`
	Name       string  `json:"name"`
	Path       string  `json:"path,omitempty"`
	Component  string  `json:"component,omitempty"`
	Icon       string  `json:"icon,omitempty"`
	Sort       int     `json:"sort"`
	Visible    bool    `json:"visible"`
	Permission *string `json:"permission,omitempty"`
	Status     int16   `json:"status"`
	Children   []Menu  `json:"children,omitempty"`
}

type MenuRequest struct {
	ParentID   uint64  `json:"parentId,string"`
	Type       string  `json:"type"`
	Name       string  `json:"name"`
	Path       string  `json:"path"`
	Component  string  `json:"component"`
	Icon       string  `json:"icon"`
	Sort       int     `json:"sort"`
	Visible    *bool   `json:"visible"`
	Permission *string `json:"permission"`
	Status     *int16  `json:"status"`
}

type CurrentMenus struct {
	Menus       []Menu   `json:"menus"`
	Permissions []string `json:"permissions"`
}

type Dictionary struct {
	ID     uint64 `json:"id,string"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	Status int16  `json:"status"`
}

type DictionaryRequest struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Status *int16 `json:"status"`
}

type DictionaryItem struct {
	ID           uint64 `json:"id,string"`
	DictionaryID uint64 `json:"dictionaryId,string"`
	Label        string `json:"label"`
	Value        string `json:"value"`
	Sort         int    `json:"sort"`
	Status       int16  `json:"status"`
}

type DictionaryItemRequest struct {
	Label  string `json:"label"`
	Value  string `json:"value"`
	Sort   int    `json:"sort"`
	Status *int16 `json:"status"`
}

type OperationLog struct {
	ID         uint64    `json:"id,string"`
	UserID     uint64    `json:"userId,string"`
	RequestID  string    `json:"requestId"`
	Method     string    `json:"method"`
	Path       string    `json:"path"`
	StatusCode int       `json:"statusCode"`
	DurationMS int64     `json:"durationMs"`
	IP         string    `json:"ip"`
	UserAgent  string    `json:"userAgent"`
	CreatedAt  time.Time `json:"createdAt"`
}

type LoginLog struct {
	ID        uint64    `json:"id,string"`
	UserID    *uint64   `json:"userId,omitempty,string"`
	Username  string    `json:"username"`
	Event     string    `json:"event"`
	Success   bool      `json:"success"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"userAgent"`
	DeviceID  string    `json:"deviceId"`
	CreatedAt time.Time `json:"createdAt"`
}
