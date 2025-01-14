package tests

import (
	"dominus/app/domain/repos"
)

type repoMock struct {
}

func (*repoMock) Migrations() {
}


func (*repoMock) InsertObject(object any, collection string) error {
	return nil 
}

func (*repoMock) DeleteObjects(collection string) error {
	return nil
}

func (*repoMock) CountPages(collection string) (int64, error) {
	return 0, nil
}

func (*repoMock) Filter(filter map[string]string, object any, collection string, page, size int) error {
	return nil
}

func (*repoMock) Backup(filter map[string]string, object any, collection string) error  {
	return nil 
}

func (*repoMock) Stats() (any, error) {
	return nil, nil
}




func NewRepoMock() repos.RepositoryInt {
	return new(repoMock)
}