package restdto_test

import (
	"dominus-project/internal/interfaces/fasthttp/dto"
	"encoding/json"
	"testing"

	"github.com/valyala/fasthttp"
)

func TestRestDto(t *testing.T) {
	t.Run("BodyParser ok", func(t *testing.T) {
		ctx := new(fasthttp.RequestCtx)
		ctx.Request.Header.SetMethod(fasthttp.MethodPost)

		data := struct {
			Name        string   `json:"name"`
			Subscribers []string `json:"subscribers"`
		}{
			Name:        "test",
			Subscribers: []string{"server1.api.com", "server2.api.com"},
		}

		payload, err := json.Marshal(data)
		if err != nil {
			t.Fatal(err)
		}

		ctx.Request.SetBody(payload)
		dto := dto.NewRestContext(ctx)

		var obj struct {
			Name        string   `json:"name"`
			Subscribers []string `json:"subscribers"`
		}

		if err := dto.BodyParser(&obj); err != nil {
			t.Fatalf("Expected nil received %s", err)
		}

	})

	t.Run("BodyParser error", func(t *testing.T) {
		ctx := new(fasthttp.RequestCtx)
		ctx.Request.Header.SetMethod(fasthttp.MethodPost)

		data := struct {
			Name        string `json:"name"`
			Subscribers int    `json:"subscribers"`
		}{
			Name:        "test",
			Subscribers: 10,
		}

		payload, err := json.Marshal(data)
		if err != nil {
			t.Fatal(err)
		}

		ctx.Request.SetBody(payload)
		dto := dto.NewRestContext(ctx)

		var obj struct {
			Name        string   `json:"name"`
			Subscribers []string `json:"subscribers"`
		}

		if err := dto.BodyParser(&obj); err == nil {
			t.Fatalf("Expected error received %s", err)
		}
	})

	t.Run("Param ok", func(t *testing.T) {
		ctx := new(fasthttp.RequestCtx)
		ctx.SetUserValue("name", "mysubscriber")
		ctx.Request.Header.SetMethod(fasthttp.MethodPatch)
		dto := dto.NewRestContext(ctx)

		param := dto.Param("name")
		if param != "mysubscriber" {
			t.Fatalf("Expected param mysubscriber received %s", param)
		}
	})

	t.Run("Param error", func(t *testing.T) {
		ctx := new(fasthttp.RequestCtx)
		ctx.SetUserValue("name", "")
		ctx.Request.Header.SetMethod(fasthttp.MethodPatch)
		ctx.Request.SetBody([]byte{})
		dto := dto.NewRestContext(ctx)

		param := dto.Param("name")
		if param != "" {
			t.Fatalf("Expected param empty received %s", param)
		}
	})
}
