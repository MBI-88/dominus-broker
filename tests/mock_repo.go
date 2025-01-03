package tests

import (
	"dominus/app/interactors"
	
)

type repoMock struct {
}

func (*repoMock) Migrations() {
}

func (*repoMock) FindObject(f any, collection string, object any) error {
	return nil 
}

func (*repoMock) FindObjects(collection string, objects any, filter any, page, sizze int) error {
	return nil 
}

func (*repoMock) InsertObject(object any, collection string) error {
	return nil 
}

func (*repoMock) UpdateObject(filter any, update any, collection string) error {
	return nil
}

func (*repoMock) DeleteObject(f any, collection string) error {
	return nil
}

func (*repoMock) DeleteObjects(f any, collection string) error {
	return nil
}

func (*repoMock) CountPages(collection string) (int64, error) {
	return 0, nil
}

func (*repoMock) Filter(filter any, object any, collection string) error {
	return nil
}

func (*repoMock) Stats() (any, error) {
	return nil, nil
}




func NewRepoMock() interactors.RepositoryInt {
	return new(repoMock)
}