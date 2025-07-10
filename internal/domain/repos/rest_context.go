package repos


type IRestContext interface {
	BodyParser(obj any) error
	Param(key string) string
}
