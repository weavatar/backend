package models

import "github.com/goravel/framework/database/orm"

type App struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	UserID uint   `gorm:"type:int;not null;index" json:"user_id"`
	Name   string `gorm:"type:varchar(255);not null" json:"name"`
	Secret string `gorm:"type:varchar(255);not null" json:"-"`
	orm.Timestamps

	// 关联
	AppAvatar []*AppAvatar `gorm:"many2many:app_avatars;joinForeignKey:AppID;joinReferences:AvatarHash"`
	// 反向关联
	User *User `gorm:"foreignKey:UserID;references:ID"`
}
