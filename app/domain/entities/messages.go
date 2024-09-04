package entities


type headers struct {
	Remote bool `json:"remote" validate:"required"`
	Topic string `json:"topic" validate:"required"`

}


type MessageRest struct {
	Headers *headers `json:"headers" validate:"required"`
	Payload []byte `json:"payload" validate:"required,max=1000"`
}
