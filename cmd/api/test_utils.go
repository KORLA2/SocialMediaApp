package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KORLA2/SocialMedia/internal/store"
)

func NewTestApplication(t *testing.T) *application {
	t.Helper()
	return &application{

		store: store.NewTestStorage(),
	}
}

func executeRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)
	return rr

}
