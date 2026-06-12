package producer

type ProducerDto interface {
	GetPayload() []byte
}
