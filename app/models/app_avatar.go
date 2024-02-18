package models

import (
	"github.com/goravel/framework/support/carbon"
)

type AppAvatar struct {
	ID         uint            `gorm:"primaryKey" json:"id"`
	AppID      uint            `json:"app_id"`
	AvatarHash string          `json:"avatar_hash"`
	CreatedAt  carbon.DateTime `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt  carbon.DateTime `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
}
