package handlers

import (
	"fmt"
	"net/http/httptest"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{
		"home", "/", "GET", []postData{}, 200,
	},
	{
		"about", "/about", "GET", []postData{}, 200,
	},
	{
		"gq", "/generals", "GET", []postData{}, 200,
	},
	{
		"sa", "/search-availability", "GET", []postData{}, 200,
	},
	{
		"contact", "/contact", "GET", []postData{}, 200,
	},
	{
		"mr", "/make-reservation", "GET", []postData{}, 200,
	},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Error(fmt.Sprintf("for %s expected %d but got %d", e.url, e.expectedStatusCode, resp.StatusCode))
			}
		}
	}
}
