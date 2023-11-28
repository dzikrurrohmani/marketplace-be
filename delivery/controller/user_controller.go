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

type UserController struct {
	ucUser usecase.UserUsecase
	api.BaseApi
}

func (uc *UserController) createNewUser(c *gin.Context) {
	var newUser model.User
	err := uc.ParseRequestBody(c, &newUser)
	if err != nil {
		uc.Failed(c, utils.RequiredError())
		return
	}
	err = utils.Retry(3, 0, uc.ucUser.CreateUser(&newUser))
	if err != nil {
		uc.Failed(c, err)
		return
	}
	uc.Success(c, newUser)
}

func (uc *UserController) userLogin(c *gin.Context) {
	var user model.User
	err := uc.ParseRequestBody(c, &user)
	if err != nil {
		uc.Failed(c, utils.RequiredError())
		return
	}
	var resp gin.H
	err = uc.ucUser.UserLogin(&resp, &user)
	if err != nil {
		uc.Failed(c, err)
		return
	}
	uc.Success(c, resp)
}

func (uc *UserController) readAllUser(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	order := c.Query("order")
	var users []model.User
	err := uc.ucUser.ReadAllUser(&users, page, limit, order)
	if err != nil {
		uc.Failed(c, err)
		return
	}
	uc.Success(c, users)
}

func (uc *UserController) updateUser(c *gin.Context) {
	var editedUser model.User
	if err := uc.ParseRequestBody(c, &editedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	} else {
		err := uc.ucUser.UpdateUser(&editedUser)
		if err != nil {
			uc.Failed(c, err)
			return
		}
		uc.Success(c, nil)
	}
}

func (uc *UserController) deleteUser(c *gin.Context) {
	var deletedUser model.User
	if err := uc.ParseRequestBody(c, &deletedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	} else {
		err := uc.ucUser.DeleteUser(&deletedUser)
		if err != nil {
			uc.Failed(c, err)
			return
		}
		uc.Success(c, nil)
	}
}

func StartUserController(router *gin.RouterGroup, ucUser usecase.UserUsecase) {
	controller := UserController{
		ucUser: ucUser,
	}
	router.POST("user", controller.createNewUser)
	router.POST("login", controller.userLogin)
	router.GET("user", ucUser.UserVerify(), controller.readAllUser)
	router.PUT("user", ucUser.UserVerify(), controller.updateUser)
	router.DELETE("user", ucUser.UserVerify(), controller.deleteUser)

}
