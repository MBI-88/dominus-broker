package adapters

type RestDto interface {
	BodyParser(obj any) error
	Param(key string) string
}
