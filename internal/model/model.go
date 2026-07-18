// Package model got schema defined for GORM and also json
package model

import "time"

const (
	RoleUser uint = 1 << iota
	RoleAdmin
	RoleModerator
)

type User struct {
	ID            int64     `gorm:"primaryKey" json:"id"`
	Username      string    `gorm:"uniqueIndex;not null" json:"username"`
	DisplayName   string    `gorm:"not null" json:"display_name"`
	Email         string    `gorm:"uniqueIndex;not null" json:"email"`
	PhoneNumber   string    `json:"phone_number"`
	Password      string    `gorm:"not null" json:"-"`
	Permission    uint      `gorm:"default:1"`
	EmailVerified bool      `gorm:"default:false" json:"email_verified"`
	PhoneVerified bool      `gorm:"default:true" json:"phone_verified"`
	CreatedAt     time.Time `json:"created_at"`
}
