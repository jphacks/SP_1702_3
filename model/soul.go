package model

import "github.com/jinzhu/gorm"

type SoulJSON struct {
	UserID   int  `json:"user_id"`
	OnDevice bool `json:"on_device"`
}

type Soul struct {
	gorm.Model
	UserID   int  `gorm:"column:user_id"`
	OnDevice bool `gorm:"column:on_device"`
}
