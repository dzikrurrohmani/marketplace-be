package repository

import (
	"store/utils"

	"gorm.io/gorm"
)

type DbRepository interface {
	Create(data interface{}) error
	Read(option utils.ReadOption) error
	Update(data interface{}) error
	Delete(data interface{}) error
}

type dbRepository struct {
	db *gorm.DB
}

func (u *dbRepository) Create(data interface{}) error {
	result := u.db.Create(data)
	return result.Error
}

func (b *dbRepository) Read(option utils.ReadOption) error {
	result := option.Read(b.db)
	return result.Error
}

func (u *dbRepository) Update(data interface{}) error {
	result := u.db.Updates(data)
	return result.Error
}

func (u *dbRepository) Delete(data interface{}) error {
	result := u.db.Delete(data)
	return result.Error
}

func NewDbRepository(db *gorm.DB) DbRepository {
	repo := new(dbRepository)
	repo.db = db
	return repo
}
