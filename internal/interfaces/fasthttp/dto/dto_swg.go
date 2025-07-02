package dto

type SwaggerTopic struct {
	// Name of the topic, e.g., prod.topic, monitoring, etc.
	Name string `json:"name" example:"prod.topic"`

	// List of subscriber URLs
	Subscribers []string `json:"subscribers" example:"[localhost:8080, localhost:8081]"`
}
