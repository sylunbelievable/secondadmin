package dto

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	AuthMode string `json:"authMode"`
	DeviceID string `json:"deviceId"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type Tokens struct {
	AccessToken  string `json:"accessToken,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
	CSRFToken    string `json:"csrfToken,omitempty"`
	ExpiresIn    int64  `json:"expiresIn"`
}

type User struct {
	ID       uint64 `json:"id,string"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Status   int16  `json:"status"`
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
}

type UpdateUserRequest struct {
	Nickname *string `json:"nickname"`
	Status   *int16  `json:"status"`
	Password *string `json:"password"`
}

type Role struct {
	ID     uint64 `json:"id,string"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	Status int16  `json:"status"`
}

type CreateRoleRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type UpdateRoleRequest struct {
	Name   *string `json:"name"`
	Status *int16  `json:"status"`
}

type API struct {
	ID     uint64 `json:"id,string"`
	Group  string `json:"group"`
	Name   string `json:"name"`
	Path   string `json:"path"`
	Method string `json:"method"`
}

type CreateAPIRequest struct {
	Group  string `json:"group"`
	Name   string `json:"name"`
	Path   string `json:"path"`
	Method string `json:"method"`
}

type IDList []uint64

func (ids *IDList) UnmarshalJSON(data []byte) error {
	var text []string
	if err := json.Unmarshal(data, &text); err == nil {
		out := make([]uint64, 0, len(text))
		for _, item := range text {
			id, err := strconv.ParseUint(strings.TrimSpace(item), 10, 64)
			if err != nil {
				return err
			}
			out = append(out, id)
		}
		*ids = out
		return nil
	}

	var numbers []uint64
	if err := json.Unmarshal(data, &numbers); err != nil {
		return err
	}
	*ids = numbers
	return nil
}

type IDsRequest struct {
	IDs IDList `json:"ids"`
}

type Session struct {
	ID        string    `json:"id"`
	DeviceID  string    `json:"deviceId"`
	AuthMode  string    `json:"authMode"`
	CreatedAt time.Time `json:"createdAt"`
}
