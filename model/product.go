package model

import "gorm.io/gorm"

type Product struct {
	gorm.Model           // tidak diberi pointer karena ketika dicari berdasarkan id nanti panic
	Code       string    `json:"productCode" gorm:"unique"`
	Name       *string   `json:"productName"`
	CashPrice  *int      `json:"productCashPrice"`
	DebtPrice  *int      `json:"productDebtPrice"`
	Stock      *int      `json:"productStock"`
	CategoryID *uint     `json:"categoryId,omitempty"`
	Category   *Category `json:"productCategory,omitempty"`
}
