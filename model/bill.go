package model

import (
	"gorm.io/gorm"
)

type Bill struct {
	gorm.Model
	Code        string       `json:"transactionCode" gorm:"unique"`
	Date        string       `json:"transactionDate" gorm:"type:date"`
	CustName    *string      `json:"customerName"`
	CustPhone   *string      `json:"customerPhone"`
	IsCash      *bool        `json:"isCash"`
	IsPaid      *bool        `json:"isPaid"`
	Debited     *int         `json:"moneyDebited"`
	GrandTotal  int          `json:"grandTotal"`
	BillDetails []BillDetail `json:"billDetails"`
	Incomes     []Income     `json:"billIncomes"`
}

func (b *Bill) BeforeCreate(tx *gorm.DB) (err error) {
	if b.Debited == nil {
		zero := 0
		b.Debited = &zero
	}
	if *b.Debited >= b.GrandTotal {
		isPaid := true
		b.IsPaid = &isPaid
	} else {
		isPaid := false
		b.IsPaid = &isPaid
	}
	if *b.Debited > 0 {
		b.Incomes = []Income{{Date: b.Date, Amount: b.Debited}}
	}
	return
}

func (b *Bill) AfterCreate(tx *gorm.DB) (err error) {
	for _, billDetail := range b.BillDetails {
		var product Product
		if tx.Where(map[string]interface{}{"code": billDetail.Code}).First(&product).Error == nil {
			remainingProduct := *product.Stock - billDetail.Qty
			if remainingProduct <= 0 {
				zero := 0
				product.Stock = &zero
			} else {
				product.Stock = &remainingProduct
			}
			tx.Updates(&product)
		}
	}
	return
}

func (b *Bill) BeforeUpdate(tx *gorm.DB) (err error) {
	var total int
	if b.GrandTotal > 0 {
		total = b.GrandTotal
	} else {
		var bill Bill
		if tx.Where(map[string]interface{}{"id": b.ID}).First(&bill).Error == nil {
			total = bill.GrandTotal
		}
	}
	if *b.Debited >= total {
		isPaid := true
		b.IsPaid = &isPaid
	} else {
		isPaid := false
		b.IsPaid = &isPaid
	}
	return
}
