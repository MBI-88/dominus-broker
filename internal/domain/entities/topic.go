package entities

import (
	"fmt"
	"sync"

	"github.com/go-playground/validator/v10"
)

type Topic struct {
	Name        string   `json:"name" validate:"omitempty,alpha,lowercase"`
	Subscribers []string `json:"subscribers" validate:"required,dive,hostname_port"`
	Message     []byte
	Lck         *sync.Mutex
	validate    *validator.Validate
}

func (t *Topic) SetMessage(data []byte)  {
	t.Lck.Lock()
	defer t.Lck.Unlock()
	t.Message = data
}

func (t *Topic) GetMessage() []byte {
	t.Lck.Lock()
	defer t.Lck.Unlock()
	return  t.Message
}

func (t *Topic) GetSubscribers() []string {
	t.Lck.Lock()
	defer t.Lck.Unlock()
	return  t.Subscribers
}

func (t *Topic) SetSubscribers(sb []string) {
	t.Lck.Lock()
	defer t.Lck.Unlock() 
	t.Subscribers = sb
}

func (t *Topic) GetName() string {
	t.Lck.Lock()
	defer t.Lck.Unlock()
	return  t.Name
}

func (t *Topic) SetName(name string) {
	t.Lck.Lock()
	defer t.Lck.Unlock() 
	t.Name = name
}


func (t *Topic) ValidateStruct() error {
	var errors []errorResponse
	if err := t.validate.Struct(t); err != nil {
		for _, er := range err.(validator.ValidationErrors) {
			var element errorResponse
			element.FailedField = er.Field()
			element.Tag = er.Tag()
			element.Value = er.Value()
			element.Error = true
			errors = append(errors, element)
		}
	}
	if errors != nil {
		errMsgs := make([]string, 0, 50)
		for _, er := range errors {
			errMsgs = append(errMsgs, fmt.Sprintf("%s: %v disagreement with %s", er.FailedField, er.Value, er.Tag))
		}
		return fmt.Errorf("%v", errMsgs)
	}
	return nil
}