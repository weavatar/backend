package models

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

type User struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	OpenID    uint            `gorm:"type:bigint;not null;unique" json:"open_id"`
	UnionID   uint            `gorm:"type:bigint;not null;unique" json:"union_id"`
	Nickname  string          `gorm:"type:varchar(255);not null;index" json:"nickname"`
	IsAdmin   bool            `gorm:"type:boolean;default:0" json:"is_admin"`
	RealName  bool            `gorm:"type:boolean;default:0" json:"real_name"`
	CreatedAt carbon.DateTime `gorm:"column:created_at" json:"created_at"`
	UpdatedAt carbon.DateTime `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt  `gorm:"column:deleted_at" json:"-"`

	// 关联
	App []*App `gorm:"foreignKey:UserID;references:ID"`
}
