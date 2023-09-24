package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type testcase struct {
	name               string // test name
	url                string
	method             string // GET, POST, PUT, DELETE, ...
	expectedStatusCode int
}

func TestRepository_Home(t *testing.T) {
	routes := getRoutes()
	testServer := httptest.NewTLSServer(routes) // server to test
	defer testServer.Close()

	test := testcase{name: "home", url: "/", method: "GET", expectedStatusCode: http.StatusOK}
	t.Run(test.name, func(t *testing.T) {
		response, err := testServer.Client().Get(testServer.URL + test.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}
		if test.expectedStatusCode != response.StatusCode {
			t.Errorf("for %s expected status %d but got %d", test.name, test.expectedStatusCode, response.StatusCode)
		}
	})
}
