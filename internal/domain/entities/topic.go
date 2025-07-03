package entities

import (
	"fmt"
	"net"
	"sync"

	"github.com/go-playground/validator/v10"
)

type topic struct {
	Name        string   `json:"name" validate:"omitempty,alpha,lowercase"`
	Subscribers []string `json:"subscribers" validate:"required,dive,hostname_port"`
	Message     []byte
	Lck         *sync.RWMutex
	validate    *validator.Validate
	parser      paserBodyInt
}

func (t *topic) SetMessage(data []byte)  {
	t.Lck.Lock()
	defer t.Lck.Unlock()
	t.Message = data
}

func (t *topic) GetMessage() []byte {
	t.Lck.RLock()
	defer t.Lck.RUnlock()
	message := t.Message
	t.Message = []byte{}
	return  message
}

func (t *topic) GetSubscribers() []string {
	t.Lck.RLock()
	defer t.Lck.RUnlock()
	return  t.Subscribers
}

func (t *topic) SetSubscribers(sb []string) {
	t.Lck.Lock()
	defer t.Lck.Unlock() 
	t.Subscribers = sb
}

func (t *topic) GetName() string {
	t.Lck.RLock()
	defer t.Lck.RUnlock()
	return  t.Name
}

func (t *topic) SetName(name string) {
	t.Lck.Lock()
	defer t.Lck.Unlock() 
	t.Name = name
}


func (t *topic) ValidateTopic() error {
	t.Lck.RLock()
	defer t.Lck.RUnlock()
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

func (t *topic) FillTopic() error {
	return  t.parser.BodyParser(t)
}



type TopicInt interface {
	SetMessage(data []byte)
	GetMessage() []byte
	SetSubscribers(sb []string)
	GetSubscribers() []string
	GetName() string
	SetName(name string)
	FillTopic() error
	ValidateTopic() error
}

type paserBodyInt interface {
	BodyParser(obj any) error
}

func NewTopic(dto paserBodyInt) TopicInt {
	customValidate := validator.New()
	customValidate.RegisterValidation("hostname_port", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		
		host, port, _ := net.SplitHostPort(value)

		isValidPort := func(port string) bool {
			if p, err := net.LookupPort("tcp", port); err == nil {
				return p > 0 && p <= 65535
			}
			return false
		}

		if net.ParseIP(host) != nil && isValidPort(port) {
			return  true
		}

		if err := validator.New().Var(value, "hostname"); err == nil {
			return  true
		}

		return false

	})
	return  &topic{
		validate: customValidate,
		Lck: new(sync.RWMutex),
		parser: dto,
	}
}