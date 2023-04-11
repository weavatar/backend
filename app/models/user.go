package models

import (
	"github.com/goravel/framework/database/orm"
)

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	OpenID   uint   `gorm:"type:bigint;not null;unique" json:"open_id"`
	UnionID  uint   `gorm:"type:bigint;not null;unique" json:"union_id"`
	Nickname string `gorm:"type:varchar(255);not null;index" json:"nickname"`
	IsAdmin  bool   `gorm:"type:boolean;default:0" json:"is_admin"`
	RealName bool   `gorm:"type:boolean;default:0" json:"real_name"`
	orm.Timestamps
	orm.SoftDeletes

	// 关联
	App []*App `gorm:"foreignKey:UserID;references:ID"`
}
