package model

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Code string  `json:"categoryCode" gorm:"unique"`
	Name *string `json:"categoryName" gorm:"unique"`
}
