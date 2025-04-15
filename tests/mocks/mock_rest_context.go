package mocks

import (
	"dominus-project/app/domain/repos"
	"encoding/json"
	"fmt"
	"os"
)

type restContext struct {
	status  bool
	payload string
	param   string
	queries map[string]string
}

func (r *restContext) BodyParser(obj any) error {
	if r.status {
		if err := json.Unmarshal(r.readJson(), obj); err != nil {
			return fmt.Errorf("Error parse real")
		}
		return nil
	}
	return fmt.Errorf("Error parse")
}

func (r *restContext) Queries() map[string]string {
	if r.status {
		return r.queries
	}
	return make(map[string]string, 0)
}

func (r *restContext) Param(key string) string {
	if r.status {
		return r.param
	}
	return ""
}

func (r *restContext) readJson() []byte {
	var file []byte
	if _, err := os.Stat(r.payload); os.IsNotExist(err) {
		panic(err)
	}
	file, _ = os.ReadFile(r.payload)
	return file
}

func NewRestConext(status bool, path, param string, queries map[string]string) repos.RestContextInt {
	return &restContext{
		status:  status,
		payload: path,
		param:  param,
	}
}
