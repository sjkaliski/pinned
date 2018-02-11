package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/beme/abide"
)

func TestGetUsersHandler(t *testing.T) {
	getRoute := func(v string) string {
		route := "/users"
		if v != "" {
			route += fmt.Sprintf("?v=%s", v)
		}
		return route
	}

	versions := vm.Versions()
	versions = append([]string{""}, versions...)

	for _, v := range versions {
		req := httptest.NewRequest(http.MethodGet, getRoute(v), nil)
		w := httptest.NewRecorder()
		getUsersHandler(w, req)
		abide.AssertHTTPResponse(t, "TestGetUsersHandler-"+getRoute(v), w.Result())
	}
}
