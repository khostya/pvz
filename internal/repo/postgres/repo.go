package postgres

import "github.com/khostya/pvz/pkg/postgres/transactor"

type Repositories struct {
	ProductRepo   *ProductRepo
	PvzRepo       *PvzRepo
	ReceptionRepo *ReceptionRepo
	UserRepo      *UserRepo
}

func NewRepositories(provider transactor.QueryEngineProvider) Repositories {
	return Repositories{
		UserRepo:      NewUserRepo(provider),
		ProductRepo:   NewProductRepo(provider),
		PvzRepo:       NewPvzRepo(provider),
		ReceptionRepo: NewReceptionRepo(provider),
	}
}
