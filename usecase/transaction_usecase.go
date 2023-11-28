package usecase

import (
	"errors"
	"store/model"
	"store/repository"
	"store/utils"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionUsecase interface {
	CreateTransaction(transaction *model.Bill) error
	ReadAllTransaction(selected *[]model.Bill, page int, limit int, order string) error
	ReadTransactionById(selected *model.Bill, transactionId int) error
	UpdateTransaction(updatedTransaction *model.Bill) error
	DeleteTransaction(deletedTransaction *model.Bill) error
	CreateIncome(income *model.Income) error
}

type transactionUsecase struct {
	dbRepo repository.DbRepository
}

func (t *transactionUsecase) CreateTransaction(transaction *model.Bill) error {
	(*transaction).Code = strings.ToUpper("T" + uuid.New().String()[:7])
	return t.dbRepo.Create(transaction)
}

func (t *transactionUsecase) ReadAllTransaction(selected *[]model.Bill, page int, limit int, order string) error {
	offset := limit * (page - 1)
	err := t.dbRepo.Read(utils.ReadOption{Limit: &limit, Offset: &offset, Many: true, Result: &selected, Order: &order})
	if err != nil {
		return err
	}
	return nil
}

func (t *transactionUsecase) ReadTransactionById(selected *model.Bill, transactionId int) error {
	err := t.dbRepo.Read(utils.ReadOption{Many: false, Result: &selected, Where: map[string]interface{}{"id": transactionId}, Preloads: []string{"BillDetails"}})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		} else {
			return err
		}
	}
	return nil
}

func (t *transactionUsecase) UpdateTransaction(updatedTransaction *model.Bill) error {
	err := t.dbRepo.Update(&updatedTransaction)
	if err != nil {
		return err
	}
	return nil
}

func (t *transactionUsecase) DeleteTransaction(deletedTransaction *model.Bill) error {
	err := t.dbRepo.Delete(&deletedTransaction)
	if err != nil {
		return err
	}
	return nil
}

func (t *transactionUsecase) CreateIncome(income *model.Income) error {
	return t.dbRepo.Create(income)
}

func NewTransactionUsecase(dbRepo repository.DbRepository) TransactionUsecase {
	usecase := new(transactionUsecase)
	usecase.dbRepo = dbRepo
	return usecase
}
