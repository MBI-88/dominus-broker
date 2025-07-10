package repos


type IGrpResponse interface {
	Descriptor() ([]byte, []int)
	GetMessage() string
	GetStatus() uint32
	ProtoMessage()
	Reset()
	String() string
	Validate() error
	ValidateAll() error
}