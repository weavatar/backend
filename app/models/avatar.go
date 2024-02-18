package models

import (
	"github.com/goravel/framework/support/carbon"
)

type Avatar struct {
	Hash      string          `gorm:"type:char(32);not null;primaryKey" json:"hash"`
	Raw       string          `gorm:"type:varchar(255);not null;unique" json:"raw"`
	UserID    uint            `gorm:"type:bigint(20);default:null;index" json:"user_id"`
	CreatedAt carbon.DateTime `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt carbon.DateTime `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`

	// 反向关联
	App []*App `gorm:"many2many:app_avatars;joinForeignKey:AvatarHash;joinReferences:AppID"`
}
