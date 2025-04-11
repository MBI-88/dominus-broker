package entities


type Topic struct {
	Name       string   `json:"name" validate:"required,alpha,lowercase"`
	Partitions []string `json:"partitions" validate:"required,dive,uri"`
	Queue      QueueInt
}

