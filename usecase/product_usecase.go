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

type ProductUsecase interface {
	CreateProduct(products *[]model.Product) error
	ReadAllProduct(selected *[]model.Product, page int, limit int, order string) error
	ReadProductById(selected *model.Product, productId int) error
	ReadProductByCategory(selected *[]model.Product, productCat string, page int, limit int, order string) error
	UpdateProduct(updatedProduct *model.Product) error
	DeleteProduct(deletedProduct *model.Product) error
}

type productUsecase struct {
	dbRepo repository.DbRepository
}

func (p *productUsecase) CreateProduct(products *[]model.Product) error {
	for i := 0; i < len(*products); i++ {
		(*products)[i].Code = strings.ToUpper("P" + uuid.New().String()[:7])
	}
	return p.dbRepo.Create(products)
}

func (p *productUsecase) ReadAllProduct(selected *[]model.Product, page int, limit int, order string) error {
	offset := limit * (page - 1)
	err := p.dbRepo.Read(utils.ReadOption{Limit: &limit, Offset: &offset, Many: true, Result: &selected, Order: &order, Preloads: []string{"Category"}})
	if err != nil {
		return err
	}
	return nil
}

func (p *productUsecase) ReadProductById(selected *model.Product, productId int) error {
	err := p.dbRepo.Read(utils.ReadOption{Many: false, Result: &selected, Where: map[string]interface{}{"id": productId}, Preloads: []string{"Category"}})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		} else {
			return err
		}
	}
	return nil
}

func (p *productUsecase) ReadProductByCategory(selected *[]model.Product, productCat string, page int, limit int, order string) error {
	var category model.Category
	err := p.dbRepo.Read(utils.ReadOption{Many: true, Result: &category, Where: map[string]interface{}{"name": productCat}})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		} else {
			return err
		}
	}
	offset := limit * (page - 1)
	err = p.dbRepo.Read(utils.ReadOption{Many: true, Result: &selected, Where: map[string]interface{}{"category_id": category.ID}, Limit: &limit, Offset: &offset, Order: &order})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		} else {
			return err
		}
	}
	return nil
}

func (p *productUsecase) UpdateProduct(updatedProduct *model.Product) error {
	err := p.dbRepo.Update(&updatedProduct)
	if err != nil {
		return err
	}
	return nil
}

func (p *productUsecase) DeleteProduct(deletedProduct *model.Product) error {
	err := p.dbRepo.Delete(&deletedProduct)
	if err != nil {
		return err
	}
	return nil
}

func NewProductUsecase(dbRepo repository.DbRepository) ProductUsecase {
	usecase := new(productUsecase)
	usecase.dbRepo = dbRepo
	return usecase
}
