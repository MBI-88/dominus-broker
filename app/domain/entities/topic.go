package entities


type Topic struct {
	Name       string   `json:"name" validate:"alpha,lowercase"`
	Partitions []string `json:"partitions" validate:"uri"`
	Queue      QueueInt
}

