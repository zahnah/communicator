package forms

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestErrors_Get(t *testing.T) {
	form := New(url.Values{})
	form.Required("a")
	err := form.Errors.Get("a")
	if err != "This field can't be blank" {
		t.Error(fmt.Sprintf("form has to return an error: %s", err))
	}

	err = form.Errors.Get("b")
	if err != "" {
		t.Error(fmt.Sprintf("form has not to return an error: %s", err))
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")

	r, _ = http.NewRequest("POST", "/whatever", nil)

	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("shows does not have required fields when it does")
	}
}

func TestForm_Has(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)

	postedData := url.Values{}
	postedData.Add("a", "a")

	r.PostForm = postedData
	form := New(r.PostForm)

	form.Has("a")
	if !form.Valid() {
		t.Error("form has to have a field a")
	}

	form.Has("b")
	if form.Valid() {
		t.Error("form has not to have a field b")
	}
}

func TestForm_MinLength(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)

	postedData := url.Values{}
	postedData.Add("a", "123456")

	r.PostForm = postedData
	form := New(r.PostForm)

	form.MinLength("a", 6)
	if !form.Valid() {
		t.Error("form hasn't to have an error")
	}

	form.MinLength("a", 7)
	if form.Valid() {
		t.Error("form has to have an error")
	}
}

func TestForm_IsEmail(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)

	postedData := url.Values{}
	postedData.Add("email", "e@mail.me")

	r.PostForm = postedData
	form := New(r.PostForm)

	form.IsEmail("email")
	if !form.Valid() {
		t.Error("email should be valid email")
	}

	postedData = url.Values{}
	postedData.Add("a", "email")

	r.PostForm = postedData
	form = New(r.PostForm)

	form.IsEmail("email")
	if form.Valid() {
		t.Error("email shouldn't be a valid email")
	}
}
