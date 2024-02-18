package models

import (
	"github.com/goravel/framework/support/carbon"

	"gorm.io/gorm"
)

type User struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	OpenID    string          `json:"open_id"`
	UnionID   string          `json:"union_id"`
	Nickname  string          `json:"nickname"`
	Avatar    string          `json:"avatar"`
	IsAdmin   bool            `json:"is_admin"`
	RealName  bool            `json:"real_name"`
	CreatedAt carbon.DateTime `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt carbon.DateTime `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt  `gorm:"column:deleted_at" json:"-"`

	// 关联
	App []*App `gorm:"foreignKey:UserID;references:ID"`
}
