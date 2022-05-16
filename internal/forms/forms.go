package forms

//server side form validation,
import (
	"fmt"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Form creates a custom form struct, embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}

// Valid returns true if there are no errors, otherwise false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

//creates a new form,New initializes a form struct
func New(data url.Values) *Form { //returns a pointer to Form
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Required checks for required fields
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

//define certain checks to see if the form data received is valid
//Has checks if form field is in post and not empty
func (f *Form) Has(field string) bool {
	x := f.Get(field)
	//if the field has any data
	if x == "" { //x== nothing
		return false
	} else {
		return true
	}
	//this has field will still display the information the user entered initially ,when the user had entered the wrong detials, without having the details in the form cleared again,so he wont have to start filling the form again
}

// MinLength checks for string minimum length
func (f *Form) MinLength(field string, length int) bool {
	x := f.Get(field)
	if len(x) < length { //if len is < the actual lenght of the field
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}
	return true
}

// IsEmail checks for valid email address
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email address")
	}
}
