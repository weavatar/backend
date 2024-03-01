package models

import (
	"github.com/goravel/framework/support/carbon"
)

type Avatar struct {
	SHA256    string          `gorm:"primaryKey" json:"sha256"`
	MD5       string          `json:"md5"`
	Raw       string          `json:"raw"`
	UserID    uint            `json:"user_id"`
	CreatedAt carbon.DateTime `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt carbon.DateTime `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`

	// 反向关联
	App []*App `gorm:"many2many:app_avatars;joinForeignKey:AvatarHash;joinReferences:AppID"`
}
