package interactors

import (
	"dominus-project/app/domain/repos"
)

type interactor struct {
	log     repos.LogsInt
	gclient repos.GrpClientInt
}

func (i interactor) NewGrpcService() GrpcServiceInt {
	return &grpcService{
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
func NewInteractor(lg repos.LogsInt) InteractorInt {
	return &interactor{
		log: lg,
	}
}
