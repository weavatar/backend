package models

import "github.com/goravel/framework/support/carbon"

type Image struct {
	Hash      string          `gorm:"primaryKey" json:"hash"`
	Ban       bool            `json:"ban"`
	CreatedAt carbon.DateTime `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt carbon.DateTime `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
}
