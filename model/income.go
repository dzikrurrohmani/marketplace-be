package model

import "gorm.io/gorm"

type Income struct {
	gorm.Model        // tidak diberi pointer karena ketika dicari berdasarkan id nanti panic
	Date       string `json:"incomeDate" gorm:"type:date"`
	Amount     *int   `json:"incomeAmount"`
	BillID     uint   `json:"billId"`
}

func (i *Income) AfterCreate(tx *gorm.DB) (err error) {
	var bill Bill
	if tx.Where(map[string]interface{}{"id": i.BillID}).First(&bill).Error != nil {
		tx.Rollback()
	}
	moneyDebited := *bill.Debited + *i.Amount
	bill.Debited = &moneyDebited
	tx.Updates(&bill)
	return
}
