package models

import (
	"github.com/goravel/framework/support/carbon"
)

type App struct {
	ID        string          `gorm:"primaryKey" json:"id"`
	UserID    string          `json:"user_id"`
	Name      string          `json:"name"`
	Secret    string          `json:"-"`
	CreatedAt carbon.DateTime `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt carbon.DateTime `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`

	// 关联
	AppAvatar []*AppAvatar `gorm:"many2many:app_avatars;joinForeignKey:AppID;joinReferences:AvatarHash"`
	// 反向关联
	User *User `gorm:"foreignKey:UserID;references:ID"`
}
