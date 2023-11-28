package manager

import "store/repository"

type RepositoryManager interface {
	DbRepo() repository.DbRepository
	TokenRepo() repository.TokenRepository
}

type repositoryManager struct {
	infra Infra
}

func (r *repositoryManager) DbRepo() repository.DbRepository {
	return repository.NewDbRepository(r.infra.SqlDb())
}

func (r *repositoryManager) TokenRepo() repository.TokenRepository {
	return repository.NewTokenRepository(r.infra.TokenConfig())
}

func NewRepositoryManager(infra Infra) RepositoryManager {
	return &repositoryManager{
		infra: infra,
	}
}
