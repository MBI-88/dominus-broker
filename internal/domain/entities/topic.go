package entities

import (
	"fmt"
	"net"
	"sync"

	"github.com/go-playground/validator/v10"
)

type Topic interface {
	SetMessage(data []byte) error
	GetMessage() []byte
	SetSubscribers(sb []string)
	GetSubscribers() []string
	GetName() string
	SetName(name string)
	ParseTopic() error
	ValidateTopic() error
}

type parserBody interface {
	BodyParser(obj any) error
}

type topic struct {
	queueLimit  int
	Name        string   `json:"name" validate:"omitempty,alpha,lowercase"`
	Subscribers []string `json:"subscribers" validate:"required,dive,hostname_port"`
	queue       Queue
	Lck         *sync.RWMutex
	validate    *validator.Validate
	parser      parserBody
}

func NewTopic(dto parserBody, limit int) Topic {
	customValidate := validator.New()
	if err := customValidate.RegisterValidation("hostname_port", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()

		host, port, _ := net.SplitHostPort(value)

		isValidPort := func(port string) bool {
			if p, err := net.LookupPort("tcp", port); err == nil {
				return p > 0 && p <= 65535
			}
			return false
		}

		if net.ParseIP(host) != nil && isValidPort(port) {
			return true
		}

		if err := validator.New().Var(value, "hostname"); err == nil {
			return true
		}

		return false

	}); err != nil {
		return nil
	}

	return &topic{
		validate:   customValidate,
		Lck:        new(sync.RWMutex),
		parser:     dto,
		queue:      NewQueue(),
		queueLimit: limit,
	}
}

func (t *topic) SetMessage(data []byte) error {
	t.Lck.Lock()
	defer t.Lck.Unlock()
	if t.queue.Length() < t.queueLimit {
		t.queue.Enqueue(data)
		return nil
	}
	return fmt.Errorf("queue full")
}

func (t *topic) GetMessage() []byte {
	t.Lck.RLock()
	defer t.Lck.RUnlock()
	msg, ok := t.queue.Dequeue()
	if !ok {
		return nil
	}
	return msg
}

func (t *topic) GetSubscribers() []string {
	t.Lck.RLock()
	defer t.Lck.RUnlock()
	return t.Subscribers
}

func (t *topic) SetSubscribers(sb []string) {
	t.Lck.Lock()
	defer t.Lck.Unlock()
	t.Subscribers = sb
}

func (t *topic) GetName() string {
	t.Lck.RLock()
	defer t.Lck.RUnlock()
	return t.Name
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

func (t *topic) ParseTopic() error {
	return t.parser.BodyParser(t)
}
