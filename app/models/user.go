package models

import (
	"github.com/goravel/framework/support/carbon"

	"gorm.io/gorm"
)

type User struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	OpenID    string          `gorm:"type:char(32);not null;unique" json:"open_id"`
	UnionID   string          `gorm:"type:char(32);not null;unique" json:"union_id"`
	Nickname  string          `gorm:"type:varchar(255);not null;index" json:"nickname"`
	Avatar    string          `gorm:"type:varchar(255);not null" json:"avatar"`
	IsAdmin   bool            `gorm:"type:boolean;default:false" json:"is_admin"`
	RealName  bool            `gorm:"type:boolean;default:false" json:"real_name"`
	CreatedAt carbon.DateTime `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt carbon.DateTime `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt  `gorm:"column:deleted_at" json:"-"`

	// 关联
	App []*App `gorm:"foreignKey:UserID;references:ID"`
}
