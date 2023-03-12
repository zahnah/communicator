package handlers

import (
	"fmt"
	"net/http/httptest"
	"net/url"
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

	{
		"search-availability-post", "/search-availability", "POST", []postData{
			{
				key:   "start",
				value: "2020-01-01",
			},
			{
				key:   "end",
				value: "2020-01-02",
			},
		}, 200,
	},
	{
		"search-availability-json", "/search-availability-json", "POST", []postData{
			{
				key:   "start",
				value: "2020-01-01",
			},
			{
				key:   "end",
				value: "2020-01-02",
			},
		}, 200,
	},
	{
		"make-reservation", "/make-reservation", "POST", []postData{
			{
				key:   "first_name",
				value: "John",
			},
			{
				key:   "last_name",
				value: "Smith",
			},
			{
				key:   "email",
				value: "me@smith.com",
			},
			{
				key:   "phone",
				value: "11111",
			},
		}, 200,
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
		} else {
			values := url.Values{}
			for _, x := range e.params {
				values.Add(x.key, x.value)
			}

			resp, err := ts.Client().PostForm(ts.URL+e.url, values)
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
