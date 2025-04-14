package interactors

import (
	"dominus-project/app/domain/entities"
	"dominus-project/app/domain/repos"
	"dominus-project/app/domain/rules"
)

type interactor struct {
	log     repos.LogsInt
	gclient repos.GrpClientInt
	topics  entities.TopicsInt
	limit   int64
}

func (i interactor) NewGrpcService() GrpcServiceInt {
	return &grpcService{
		lg:     i.log,
		client: i.gclient,
		topics: i.topics,
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

func (i interactor) NewManagerService() ManagerInt {
	return &managerService{
		topics: i.topics,
		rls:    rules.NewValidator(),
		limit:  i.limit,
	}
}

type InteractorInt interface {
	NewGrpcService() GrpcServiceInt
	NewSystemService() SystemServiceInt
	Set(cl repos.GrpClientInt) InteractorInt
}

// Create a new interactor instance
func NewInteractor(lg repos.LogsInt, t entities.TopicsInt, l int64) InteractorInt {
	return &interactor{
		log:    lg,
		topics: t,
		limit:  l,
	}
}
