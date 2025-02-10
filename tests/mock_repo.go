package tests

import (
	"dominus/app/domain/repos"
)

type repoMock struct {
	happyPath bool
	mockData string
}

func (r *repoMock) Migrations() {
}


func (r *repoMock) InsertObject(object any, collection string) error {
	return nil 
}

func (r *repoMock) DeleteObjects(collection string) error {
	return nil
}

func (r *repoMock) CountPages(collection string) (int64, error) {
	return 0, nil
}

func (r *repoMock) Filter(filter map[string]string, object any, collection string, page, size int) error {
	return nil
}

func (r *repoMock) Backup(filter map[string]string, object any, collection string) error  {
	return nil 
}

func (r *repoMock) Stats() (any, error) {
	return nil, nil
}




func NewRepoMock(happyPath bool, mock string) repos.RepositoryInt {
	return &repoMock{
		happyPath: happyPath,
		mockData: mock,
	}
}