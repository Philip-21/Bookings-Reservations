package main

import (
	"fmt"
	"testing"

	"github.com/Philip-21/bookings/internal/config"
	"github.com/go-chi/chi"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	router := routes(&app)

	switch v := router.(type) {
	case *chi.Mux:
		// do nothing; test passed
	default:
		t.Error(fmt.Sprintf("type is not *chi.router, type is %T", v))
	}
}
