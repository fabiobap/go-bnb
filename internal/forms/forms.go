package forms

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

type Form struct {
	url.Values
	Errors errors
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, fmt.Sprintf("Field %s cannot be blank", field))
		}
	}
}

func (f *Form) Has(field string) bool {
	ff := f.Get(field)

	if ff == "" {
		return false
	}

	return true
}

func (f *Form) Minlength(field string, length int) bool {
	ff := f.Get(field)
	if len(ff) < length {
		f.Errors.Add(field, fmt.Sprintf("The %s field must be at least %d characters long", field, length))
		return false
	}
	return true
}

func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid Email Address")
	}
}
