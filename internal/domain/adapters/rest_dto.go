package adapters

type IRestDto interface {
	BodyParser(obj any) error
	Param(key string) string
}
