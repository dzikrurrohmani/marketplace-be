package manager

import (
	"log"
	"os"
	"store/config"
	"store/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Infra interface {
	SqlDb() *gorm.DB
	TokenConfig() config.TokenConfig
}

type infra struct {
	dbResource  *gorm.DB
	tokenConfig config.TokenConfig
}

func (i *infra) SqlDb() *gorm.DB {
	return i.dbResource
}

func (i *infra) TokenConfig() config.TokenConfig {
	return i.tokenConfig
}

func NewInfra(config config.Config) Infra {
	resource, err := initDbResource(config.DataSourceName)
	if err != nil {
		log.Fatal(err.Error())
	}
	return &infra{dbResource: resource, tokenConfig: config.TokenConfig}
}

func initDbResource(dataSourceName string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dataSourceName), &gorm.Config{})

	env := os.Getenv("ENV")
	if env == "development" {
		db.Debug()
		db.AutoMigrate(
			&model.User{},
			&model.Product{},
			&model.Category{},
			&model.Bill{},
			&model.BillDetail{},
			&model.Income{})

	} else if env == "production" {
		db.Debug()
	}
	if err != nil {
		return nil, err
	}
	return db, nil
}
