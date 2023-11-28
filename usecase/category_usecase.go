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

type CategoryUsecase interface {
	CreateCategory(category *model.Category) error
	ReadAllCategory(selected *[]model.Category, page int, limit int) error
	ReadCategoryById(selected *model.Category, categoryId int) error
	UpdateCategory(updatedCategory *model.Category) error
	DeleteCategory(deletedCategory *model.Category) error
}

type categoryUsecase struct {
	dbRepo repository.DbRepository
}

func (c *categoryUsecase) CreateCategory(category *model.Category) error {
	(*category).Code = strings.ToUpper("C" + uuid.New().String()[:7])
	return c.dbRepo.Create(category)
}

func (c *categoryUsecase) ReadAllCategory(selected *[]model.Category, page int, limit int) error {
	offset := limit * (page - 1)
	err := c.dbRepo.Read(utils.ReadOption{Limit: &limit, Offset: &offset, Many: true, Result: &selected})
	if err != nil {
		return err
	}
	return nil
}

func (c *categoryUsecase) ReadCategoryById(selected *model.Category, categoryId int) error {
	err := c.dbRepo.Read(utils.ReadOption{Many: false, Result: &selected, Where: map[string]interface{}{"id": categoryId}})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		} else {
			return err
		}
	}
	return nil
}

func (c *categoryUsecase) UpdateCategory(updatedCategory *model.Category) error {
	err := c.dbRepo.Update(&updatedCategory)
	if err != nil {
		return err
	}
	return nil
}

func (c *categoryUsecase) DeleteCategory(deletedCategory *model.Category) error {
	err := c.dbRepo.Delete(&deletedCategory)
	if err != nil {
		return err
	}
	return nil
}

func NewCategoryUsecase(dbRepo repository.DbRepository) CategoryUsecase {
	usecase := new(categoryUsecase)
	usecase.dbRepo = dbRepo
	return usecase
}
