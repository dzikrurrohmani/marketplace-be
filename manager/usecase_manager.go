package manager

import "store/usecase"

type UsecaseManager interface {
	UserUsecase() usecase.UserUsecase
	ProductUsecase() usecase.ProductUsecase
	CategoryUsecase() usecase.CategoryUsecase
	TransactionUsecase() usecase.TransactionUsecase
}

type usecaseManager struct {
	repoManager RepositoryManager
}

func (u *usecaseManager) UserUsecase() usecase.UserUsecase {
	return usecase.NewUserUsecase(u.repoManager.DbRepo(), u.repoManager.TokenRepo())
}

func (u *usecaseManager) ProductUsecase() usecase.ProductUsecase {
	return usecase.NewProductUsecase(u.repoManager.DbRepo())
}

func (u *usecaseManager) CategoryUsecase() usecase.CategoryUsecase {
	return usecase.NewCategoryUsecase(u.repoManager.DbRepo())
}

func (u *usecaseManager) TransactionUsecase() usecase.TransactionUsecase {
	return usecase.NewTransactionUsecase(u.repoManager.DbRepo())
}

func NewUsecaseManager(repoManager RepositoryManager) UsecaseManager {
	return &usecaseManager{
		repoManager: repoManager,
	}
}
