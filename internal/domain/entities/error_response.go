package entities


type errorResponse struct {
	Error       bool
	FailedField string
	Tag         string
	Value       any
}