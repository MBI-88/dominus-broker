package interactors

import (
	"dominus/app/domain/repos"
	"dominus/app/domain/rules"
)

type interactor struct {
	log     repos.LogsInt
	gclient repos.GrpClientInt
	rls     rules.RulesInt
}

func (i interactor) NewConnection() ConnectionInt {
	return &connection{
		rls:        i.rls,
		lg:         i.log,
		client:     i.gclient,
	}
}

func (i interactor) NewManager() ManagerInt {
	return &manager{
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
func NewInteractor(lg repos.LogsInt, rls rules.RulesInt) InteractorInt {
	return &interactor{
		log:  lg,
		rls:  rls,
	}
}
