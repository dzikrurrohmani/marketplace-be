package controller

import (
	"fmt"
	"net/http"
	"store/delivery/api"
	"store/model"
	"store/usecase"
	"store/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	ucProduct usecase.ProductUsecase
	api.BaseApi
}

func (pc *ProductController) createNewProduct(c *gin.Context) {
	var newProducts struct {
		Products []model.Product
	}
	err := pc.ParseRequestBody(c, &newProducts)
	if err != nil {
		pc.Failed(c, utils.RequiredError())
		return
	}
	err = utils.Retry(3, 0, pc.ucProduct.CreateProduct(&newProducts.Products))
	if err != nil {
		pc.Failed(c, err)
		return
	}
	pc.Success(c, newProducts.Products)
}

func (pc *ProductController) readAllProduct(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	order := c.Query("order")
	fmt.Println(page, limit, order)
	var products []model.Product
	err := pc.ucProduct.ReadAllProduct(&products, page, limit, order)
	if err != nil {
		pc.Failed(c, err)
		return
	}
	pc.Success(c, products)
}

func (pc *ProductController) readProductById(c *gin.Context) {
	productId, _ := strconv.Atoi(c.Param("id"))
	var products model.Product
	err := pc.ucProduct.ReadProductById(&products, productId)
	if err != nil {
		pc.Failed(c, err)
		return
	}
	pc.Success(c, products)
}

func (pc *ProductController) readProductByCategory(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	order := c.Query("order")
	productCat := c.Param("cat")
	var products []model.Product
	err := pc.ucProduct.ReadProductByCategory(&products, productCat, page, limit, order)
	if err != nil {
		pc.Failed(c, err)
		return
	}
	pc.Success(c, products)
}

func (pc *ProductController) updateProduct(c *gin.Context) {
	var editedProduct model.Product
	if err := pc.ParseRequestBody(c, &editedProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	} else {
		err := pc.ucProduct.UpdateProduct(&editedProduct)
		if err != nil {
			pc.Failed(c, err)
			return
		}
		pc.Success(c, nil)
	}
}

func (pc *ProductController) deleteProduct(c *gin.Context) {
	var deletedProduct model.Product
	if err := pc.ParseRequestBody(c, &deletedProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	} else {
		err := pc.ucProduct.DeleteProduct(&deletedProduct)
		if err != nil {
			pc.Failed(c, err)
			return
		}
		pc.Success(c, nil)
	}
}

func StartProductController(router *gin.RouterGroup, ucProduct usecase.ProductUsecase, ucUser usecase.UserUsecase) {
	controller := ProductController{
		ucProduct: ucProduct,
	}
	router.Use(ucUser.UserVerify())
	router.POST("product", controller.createNewProduct)
	router.GET("product", controller.readAllProduct)
	router.GET("product/detail/:id", controller.readProductById)
	router.GET("product/:cat", controller.readProductByCategory)
	router.PUT("product", controller.updateProduct)
	router.DELETE("product", controller.deleteProduct)

}
