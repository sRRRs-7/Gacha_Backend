// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0

package db

import (
	"time"
)

type Account struct {
	ID        int64     `json:"id"`
	Owner     string    `json:"owner"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

type Category struct {
	ID        int64     `json:"id"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
}

type Gacha struct {
	ID        int64     `json:"id"`
	AccountID int64     `json:"account_id"`
	ItemID    int64     `json:"item_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Gallery struct {
	ID        int64     `json:"id"`
	OwnerID   int64     `json:"owner_id"`
	ItemID    int64     `json:"item_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Item struct {
	ID         int64     `json:"id"`
	ItemName   string    `json:"item_name"`
	Rating     int32     `json:"rating"`
	ItemUrl    string    `json:"item_url"`
	CategoryID int32     `json:"category_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type Session struct {
	ID        int64     `json:"id"`
	UserName  string    `json:"user_name"`
	UserAgent string    `json:"user_agent"`
	ClientIp  string    `json:"client_ip"`
	IsBlocked bool      `json:"is_blocked"`
	ExpiredAt time.Time `json:"expired_at"`
}

type User struct {
	ID           int64     `json:"id"`
	UserName     string    `json:"user_name"`
	HashPassword string    `json:"hash_password"`
	FullName     string    `json:"full_name"`
	Email        string    `json:"email"`
	CreatedAt    time.Time `json:"created_at"`
}
