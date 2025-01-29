package interactors

import (
	"dominus/app/domain/repos"
	"dominus/app/domain/rules"
)

type interactor struct {
	repo    repos.RepositoryInt
	log     LogsInt
	gclient repos.GrpClientInt
	rls     rules.RulesInt
}

func (i interactor) NewConnection() ConnectionInt {
	return &connection{
		rls:        i.rls,
		lg:         i.log,
		collection: "logs",
		client:     i.gclient,
		repo:       i.repo,
	}
}

func (i interactor) NewManager() ManagerInt {
	return &manager{
		rls:        i.rls,
		repo:       i.repo,
		collection: "logs",
		lg:         i.log,
	}
}

func (i interactor) Set(c repos.GrpClientInt) InteractorInt {
	i.gclient = c
	return i
}

type InteractorInt interface {
	NewConnection() ConnectionInt
	NewManager() ManagerInt
	Set(cl repos.GrpClientInt) InteractorInt
}

// Create a new interactor instance
func NewInteractor(rp repos.RepositoryInt, lg LogsInt, rls rules.RulesInt) InteractorInt {
	return &interactor{
		repo: rp,
		log:  lg,
		rls:  rls,
	}
}
