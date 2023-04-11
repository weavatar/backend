package models

import "github.com/goravel/framework/database/orm"

type AppAvatar struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	AppID      uint   `gorm:"type:bigint;not null;index" json:"app_id"`
	AvatarHash string `gorm:"type:char(32);not null;index" json:"avatar_hash"`
	Ban        bool   `gorm:"type:boolean;default:0" json:"ban"`
	Checked    bool   `gorm:"type:boolean;default:0" json:"check"`
	orm.Timestamps
}
