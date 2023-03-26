package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zahnah/study-app/internal/models"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
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

	//{
	//	"search-availability-post", "/search-availability", "POST", []postData{
	//		{
	//			key:   "start",
	//			value: "2020-01-01",
	//		},
	//		{
	//			key:   "end",
	//			value: "2020-01-02",
	//		},
	//	}, 200,
	//},
	//{
	//	"search-availability-json", "/search-availability-json", "POST", []postData{
	//		{
	//			key:   "start",
	//			value: "2020-01-01",
	//		},
	//		{
	//			key:   "end",
	//			value: "2020-01-02",
	//		},
	//	}, 200,
	//},
	//{
	//	"make-reservation", "/make-reservation", "POST", []postData{
	//		{
	//			key:   "first_name",
	//			value: "John",
	//		},
	//		{
	//			key:   "last_name",
	//			value: "Smith",
	//		},
	//		{
	//			key:   "email",
	//			value: "me@smith.com",
	//		},
	//		{
	//			key:   "phone",
	//			value: "11111",
	//		},
	//	}, 307,
	//},
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

func TestRepository_MakeReservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	req, _ := http.NewRequest("GET", "/make-reservations", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.MakeReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// reservation is not in the session
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// non existent room
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	reservation = models.Reservation{
		RoomID: 100,
	}
	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_PostMakeReservation(t *testing.T) {
	// reservation is not in the session
	reservation := models.Reservation{
		RoomID:    1,
		StartDate: time.Now(),
		EndDate:   time.Now(),
	}

	postedData := url.Values{}
	postedData.Add("first_name", "John")
	postedData.Add("last_name", "Smith")
	postedData.Add("email", "smith@email.local")
	postedData.Add("phone", "111-111")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx := getCtx(req)
	session.Put(ctx, "reservation", reservation)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostMakeReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// missed post body
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	session.Put(ctx, "reservation", reservation)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostMakeReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// missed session reservation
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostMakeReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// validation error
	postedData = url.Values{}
	postedData.Add("first_name", "")
	postedData.Add("last_name", "")
	postedData.Add("email", "smith")
	postedData.Add("phone", "111-111")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	session.Put(ctx, "reservation", reservation)

	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostMakeReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// cannot save reservation
	reservation = models.Reservation{
		RoomID:    3,
		StartDate: time.Now(),
		EndDate:   time.Now(),
	}

	postedData = url.Values{}
	postedData.Add("first_name", "John")
	postedData.Add("last_name", "Smith")
	postedData.Add("email", "smith@email.local")
	postedData.Add("phone", "111-111")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	session.Put(ctx, "reservation", reservation)

	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostMakeReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// cannot save room restriction
	reservation = models.Reservation{
		RoomID:    2,
		StartDate: time.Now(),
		EndDate:   time.Now(),
	}
	postedData = url.Values{}
	postedData.Add("first_name", "John")
	postedData.Add("last_name", "Smith")
	postedData.Add("email", "smith@email.local")
	postedData.Add("phone", "111-111")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	session.Put(ctx, "reservation", reservation)

	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostMakeReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_PostAvailabilityJSON(t *testing.T) {
	// cannot save reservation
	postedData := url.Values{}
	postedData.Add("start", "2050-01-02")
	postedData.Add("end", "2050-01-03")
	postedData.Add("room_id", "1")

	req, _ := http.NewRequest("POST", "/search-availability-json", strings.NewReader(postedData.Encode()))
	ctx := getCtx(req)

	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostAvailabilityJSON)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	var j jsonResponse
	err := json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("Failed parse request")
	}
	if j.Message != "Not Available" {
		t.Errorf("Wrong response, expected '%s', reseived '%s'", "Not Available", j.Message)
	}

	// search error
	postedData = url.Values{}
	postedData.Add("start", "2050-01-02")
	postedData.Add("end", "2050-01-03")
	postedData.Add("room_id", "3")

	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)

	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailabilityJSON)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	err = json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("Failed parse request")
	}
	if j.Message != "Couldn't find data" {
		t.Errorf("Wrong response, expected '%s', reseived '%s'", "Couldn't find data", j.Message)
	}

	// search error

	req, _ = http.NewRequest("POST", "/search-availability-json", nil)
	ctx = getCtx(req)

	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailabilityJSON)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusInternalServerError)
	}

	err = json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("Failed parse request")
	}
	if j.Message != "Internal server error" {
		t.Errorf("Wrong response, expected '%s', reseived '%s'", "Internal server error", j.Message)
	}
}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println("Error")
	}

	return ctx
}
