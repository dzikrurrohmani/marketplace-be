package controller

import (
	"net/http"
	"store/delivery/api"
	"store/model"
	"store/usecase"
	"store/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	ucTransaction usecase.TransactionUsecase
	api.BaseApi
}

func (tc *TransactionController) createNewTransaction(c *gin.Context) {
	var newTransaction model.Bill
	err := tc.ParseRequestBody(c, &newTransaction)
	if err != nil {
		tc.Failed(c, utils.RequiredError())
		return
	}
	err = utils.Retry(3, 0, tc.ucTransaction.CreateTransaction(&newTransaction))
	if err != nil {
		tc.Failed(c, err)
		return
	}
	tc.Success(c, newTransaction)
}

func (tc *TransactionController) readAllTransaction(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	order := c.Query("order")
	var transactions []model.Bill
	err := tc.ucTransaction.ReadAllTransaction(&transactions, page, limit, order)
	if err != nil {
		tc.Failed(c, err)
		return
	}
	tc.Success(c, transactions)
}

func (tc *TransactionController) readTransactionById(c *gin.Context) {
	transactionId, _ := strconv.Atoi(c.Param("id"))
	var transactions model.Bill
	err := tc.ucTransaction.ReadTransactionById(&transactions, transactionId)
	if err != nil {
		tc.Failed(c, err)
		return
	}
	tc.Success(c, transactions)

}

func (tc *TransactionController) updateTransaction(c *gin.Context) {
	var editedTransaction model.Bill
	if err := tc.ParseRequestBody(c, &editedTransaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	} else {
		err := tc.ucTransaction.UpdateTransaction(&editedTransaction)
		if err != nil {
			tc.Failed(c, err)
			return
		}
		tc.Success(c, nil)
	}
}

func (tc *TransactionController) deleteTransaction(c *gin.Context) {
	var deletedTransaction model.Bill
	if err := tc.ParseRequestBody(c, &deletedTransaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	} else {
		err := tc.ucTransaction.DeleteTransaction(&deletedTransaction)
		if err != nil {
			tc.Failed(c, err)
			return
		}
		tc.Success(c, nil)
	}
}

func (tc *TransactionController) addTransactionIncome(c *gin.Context) {
	var newIncome model.Income
	err := tc.ParseRequestBody(c, &newIncome)
	if err != nil {
		tc.Failed(c, utils.RequiredError())
		return
	}
	err = tc.ucTransaction.CreateIncome(&newIncome)
	if err != nil {
		tc.Failed(c, err)
		return
	}
	tc.Success(c, newIncome)
}

func StartTransactionController(router *gin.RouterGroup, ucTransaction usecase.TransactionUsecase, ucUser usecase.UserUsecase) {
	controller := TransactionController{
		ucTransaction: ucTransaction,
	}
	router.Use(ucUser.UserVerify())
	router.POST("transaction", controller.createNewTransaction)
	router.POST("transaction/income", controller.addTransactionIncome)
	router.GET("transaction", controller.readAllTransaction)
	router.GET("transaction/detail/:id", controller.readTransactionById)
	router.PUT("transaction", controller.updateTransaction)
	router.DELETE("transaction", controller.deleteTransaction)

}
