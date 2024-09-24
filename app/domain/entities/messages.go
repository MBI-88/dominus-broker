package entities


// headers define topic for local y receiver for remote
type headers struct {
	Topic    string `json:"topic" validate:"required"`
	Receiver string `json:"receiver,omitempty"`
}

// Meesage cotains headers and body 
type Message struct {
	Headers headers `json:"headers" validate:"required"`
	Payload []byte   `json:"payload" validate:"required,max=1000"`
}
