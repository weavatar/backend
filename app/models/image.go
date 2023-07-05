package models

import "github.com/goravel/framework/support/carbon"

type Image struct {
	Hash      string          `gorm:"type:char(32);not null;primaryKey" json:"hash"`
	Ban       bool            `gorm:"type:boolean;default:0" json:"ban"`
	CreatedAt carbon.DateTime `gorm:"column:created_at" json:"created_at"`
	UpdatedAt carbon.DateTime `gorm:"column:updated_at" json:"updated_at"`
}
