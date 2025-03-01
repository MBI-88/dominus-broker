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

func (i interactor) NewGrpcService() GrpcServiceInt {
	return &grpcService{
		rls:    i.rls,
		lg:     i.log,
		client: i.gclient,
	}
}

func (i interactor) NewSystemService() SystemServiceInt {
	return &systemService{
		lg: i.log,
	}
}

func (i interactor) Set(c repos.GrpClientInt) InteractorInt {
	i.gclient = c
	return i
}

type InteractorInt interface {
	NewGrpcService() GrpcServiceInt
	NewSystemService() SystemServiceInt
	Set(cl repos.GrpClientInt) InteractorInt
}

// Create a new interactor instance
func NewInteractor(lg repos.LogsInt, rls rules.RulesInt) InteractorInt {
	return &interactor{
		log: lg,
		rls: rls,
	}
}
