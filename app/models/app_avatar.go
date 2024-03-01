package models

import (
	"github.com/goravel/framework/support/carbon"
)

type AppAvatar struct {
	ID           uint            `gorm:"primaryKey" json:"id"`
	AppID        uint            `json:"app_id"`
	AvatarSHA256 string          `json:"avatar_sha256"`
	CreatedAt    carbon.DateTime `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt    carbon.DateTime `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
}
