package tests

import (
	"net/http"
	"net/url"
	"testing"
	"url-shortener/internal/http-server/handlers/model"
	"url-shortener/internal/lib/random"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gavv/httpexpect/v2"
)

const (
	host = "localhost:8082"
)

func TestURLShortenerCreated(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}

	e := httpexpect.Default(t, u.String())

	e.POST("/url/save").WithJSON(model.Request{
		URL:   gofakeit.URL(),
		Alias: random.NewRandomString(10),
	}).Expect().Status(http.StatusCreated).JSON().Object().ContainsKey("alias")
}
