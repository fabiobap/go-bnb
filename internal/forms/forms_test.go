package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestFormValid(t *testing.T) {
	r := httptest.NewRequest("GET", "/whtever", nil)
	form := New(r.PostForm)

	isValid := form.Valid()

	if !isValid {
		t.Error("Form is not valid!")
	}
}

func TestFormRequired(t *testing.T) {
	r := httptest.NewRequest("GET", "/whtever", nil)
	form := New(r.PostForm)

	form.Required("field_1")

	isValid := form.Valid()

	if isValid {
		t.Error("Form is valuating missing field as required")
	}

	postData := url.Values{}
	postData.Add("field_1", "123")

	rn := httptest.NewRequest("GET", "/whtever", nil)
	rn.PostForm = postData
	newForm := New(rn.PostForm)

	newForm.Required("field_1")

	if !newForm.Valid() {
		t.Error("Form has the required field but it's missing somehow")
	}
}

func TestFormEmailIsValid(t *testing.T) {
	postData := url.Values{}
	postData.Add("email", "email@test.org")

	r := httptest.NewRequest("GET", "/whtever", nil)
	r.PostForm = postData
	form := New(r.PostForm)

	form.IsEmail("email")

	if !form.Valid() {
		t.Error("Form has correct email field but somehow is invalid")
	}
}

func TestFormEmailIsNotValid(t *testing.T) {
	postData := url.Values{}
	postData.Add("email", "notvalidemail")

	r := httptest.NewRequest("GET", "/whtever", nil)
	r.PostForm = postData
	form := New(r.PostForm)

	form.IsEmail("email")

	if form.Valid() {
		t.Error("Form has the incorrect email field but somehow is valid")
	}
}

func TestFormEmtpyEmail(t *testing.T) {
	r := httptest.NewRequest("GET", "/whtever", nil)
	form := New(r.PostForm)

	form.IsEmail("email")

	if form.Valid() {
		t.Error("Form has inexistent email field but somehow is valid")
	}
}

func TestFormHasField(t *testing.T) {
	postData := url.Values{}
	postData.Add("field", "123")

	r := httptest.NewRequest("GET", "/whtever", nil)
	r.PostForm = postData
	form := New(r.PostForm)

	if !form.Has("field") {
		t.Error("Form has the correct field but somehow is invalid")
	}
}

func TestFormHasNotField(t *testing.T) {
	r := httptest.NewRequest("GET", "/whtever", nil)
	form := New(r.PostForm)

	if form.Has("field") {
		t.Error("Form hasn't the correct field but somehow is valid")
	}
}

func TestFormMinLength(t *testing.T) {
	postData := url.Values{}
	postData.Add("field", "123")

	r := httptest.NewRequest("GET", "/whtever", nil)
	r.PostForm = postData
	form := New(r.PostForm)

	form.Minlength("field", 3)

	if !form.Valid() {
		t.Error("Form has the correct field length but somehow is invalid")
	}

	isError := form.Errors.Get("field")
	if isError != "" {
		t.Error("Should not have an error but got one")
	}
}

func TestFormEmptyMinLength(t *testing.T) {
	r := httptest.NewRequest("GET", "/whtever", nil)
	form := New(r.PostForm)

	form.Minlength("field", 3)

	if form.Valid() {
		t.Error("Minlength is valuating a inexistent field")
	}

	isError := form.Errors.Get("field")
	if isError == "" {
		t.Error("Should have an error but got none")
	}

}

func TestFormInvalidMinLength(t *testing.T) {
	postData := url.Values{}
	postData.Add("field", "1")

	r := httptest.NewRequest("GET", "/whtever", nil)
	r.PostForm = postData
	form := New(r.PostForm)

	form.Minlength("field", 3)

	if form.Valid() {
		t.Error("Form has the incorrect field length but somehow is valid")
	}
}
