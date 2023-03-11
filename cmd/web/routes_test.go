package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/zahnah/study-app/internal/config"
	"testing"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)
	switch v := mux.(type) {
	case *chi.Mux:
	default:
		t.Error(fmt.Sprintf("type is not http.Handler, but is %T", v))
	}
}
