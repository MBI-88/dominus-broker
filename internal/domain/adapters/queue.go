package adapters


type ProviderDto interface {
	GetPayload() []byte
}

type ConsumerDto interface {
	GetId() string
}

type ProviderClient interface {

}

type ConsumerClient interface {
	
}