package models

import (
	"github.com/goravel/framework/support/carbon"
)

type AppAvatar struct {
	ID         uint            `gorm:"primaryKey" json:"id"`
	AppID      uint            `gorm:"type:bigint;not null;index" json:"app_id"`
	AvatarHash string          `gorm:"type:char(32);not null;index" json:"avatar_hash"`
	CreatedAt  carbon.DateTime `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt  carbon.DateTime `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
}
