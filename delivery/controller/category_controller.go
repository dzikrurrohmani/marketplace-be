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

type CategoryController struct {
	ucCategory usecase.CategoryUsecase
	api.BaseApi
}

func (cc *CategoryController) createNewCategory(c *gin.Context) {
	var newCategory model.Category
	err := cc.ParseRequestBody(c, &newCategory)
	if err != nil {
		cc.Failed(c, utils.RequiredError())
		return
	}
	err = utils.Retry(3, 0, cc.ucCategory.CreateCategory(&newCategory))
	if err != nil {
		cc.Failed(c, err)
		return
	}
	cc.Success(c, newCategory)
}

func (cc *CategoryController) readAllCategory(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	var categories []model.Category
	err := cc.ucCategory.ReadAllCategory(&categories, page, limit)
	if err != nil {
		cc.Failed(c, err)
		return
	}
	cc.Success(c, categories)

}

func (cc *CategoryController) readCategoryById(c *gin.Context) {
	paramId := c.Param("id")
	categoryId, _ := strconv.Atoi(paramId)
	var categories model.Category
	err := cc.ucCategory.ReadCategoryById(&categories, categoryId)
	if err != nil {
		cc.Failed(c, err)
		return
	}
	cc.Success(c, categories)

}

func (cc *CategoryController) updateCategory(c *gin.Context) {
	var editedCategory model.Category
	if err := cc.ParseRequestBody(c, &editedCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	} else {
		err := cc.ucCategory.UpdateCategory(&editedCategory)
		if err != nil {
			cc.Failed(c, err)
			return
		}
		cc.Success(c, nil)
	}
}

func (cc *CategoryController) deleteCategory(c *gin.Context) {
	var deletedCategory model.Category
	if err := cc.ParseRequestBody(c, &deletedCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	} else {
		err := cc.ucCategory.DeleteCategory(&deletedCategory)
		if err != nil {
			cc.Failed(c, err)
			return
		}
		cc.Success(c, nil)
	}
}

func StartCategoryController(router *gin.RouterGroup, ucCategory usecase.CategoryUsecase, ucUser usecase.UserUsecase) {
	controller := CategoryController{
		ucCategory: ucCategory,
	}
	router.Use(ucUser.UserVerify())
	router.POST("category", controller.createNewCategory)
	router.GET("category", controller.readAllCategory)
	router.GET("category/:id", controller.readCategoryById)
	router.PUT("category", controller.updateCategory)
	router.DELETE("category", controller.deleteCategory)

}
