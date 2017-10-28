package model

import "github.com/jinzhu/gorm"

//IndoorInfoJSON is users indoor information for JSON
type IndoorInfoJSON struct {
	Temperature  float32 `json:"temperature"`
	Humidity     int     `json:"humidity"`
	Illumination int     `json:"illumination"`
}

//IndoorInfo is users indoor information for DB
type IndoorInfo struct {
	gorm.Model
	Temperature  float32 `gorm:"column:temperature"`
	Humidity     int     `gorm:"column:humidity"`
	Illumination int     `gorm:"column:illumination"`
}
