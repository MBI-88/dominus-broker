package repos


type RestContextInt interface {
	BodyParser(obj any) error
	Param(key string) string
}
