package model

import "gorm.io/gorm"

type BillDetail struct {
	gorm.Model
	Code   string `json:"productCode"`
	Name   string `json:"productName"`
	Price  int    `json:"productPrice"`
	Qty    int    `json:"productQty"`
	BillID uint   `json:"billId"`
}
