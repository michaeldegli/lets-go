package forms

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

// NewSnippet ...
type NewSnippet struct {
	Title    string
	Content  string
	Expires  string
	Failures map[string]string
}

// SignupUser ...
type SignupUser struct {
	Name     string
	Email    string
	Password string
	Failures map[string]string
}

// LoginUser ...
type LoginUser struct {
	Email    string
	Password string
	Failures map[string]string
}

// Valid ...
func (f *NewSnippet) Valid() bool {

	f.Failures = make(map[string]string)

	if strings.TrimSpace(f.Title) == "" {
		f.Failures["Title"] = "Title is required."
	} else if utf8.RuneCountInString(f.Title) > 100 {
		f.Failures["Title"] = "Title cannot be longer than 100 characters."
	}

	if strings.TrimSpace(f.Content) == "" {
		f.Failures["Content"] = "Content is required."
	}

	permitted := map[string]bool{"3600": true, "86400": true, "31536000": true}

	if strings.TrimSpace(f.Expires) == "" {
		f.Failures["Expires"] = "Expiry time is required."
	} else if !permitted[f.Expires] {
		f.Failures["Expires"] = "Expiry time must conform to default form values."
	}

	// Multiple Form Values
	// if len(r.PostForm["items"]) == 0 {
	// 	failures["items"] = "At least one item must be checked"
	// }
	// for i, item := range r.PostForm["items"] {
	// 	fmt.Fprintf(w, "%d: Item %s\n", i, item)
	// }

	return len(f.Failures) == 0
}

var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Valid ...
func (f *SignupUser) Valid() bool {

	f.Failures = make(map[string]string)

	if strings.TrimSpace(f.Name) == "" {
		f.Failures["Name"] = "Name is required."
	}

	if strings.TrimSpace(f.Email) == "" {
		f.Failures["Email"] = "Email is required."
	} else if len(f.Email) > 254 || !rxEmail.MatchString(f.Email) {
		f.Failures["Email"] = "Email is not valid."
	}

	passLen := 2
	if utf8.RuneCountInString(f.Password) < passLen {
		f.Failures["Password"] = fmt.Sprintf("Password is too short. Should be at least %d chars", passLen)
	}

	return len(f.Failures) == 0
}

// Valid ...
func (f *LoginUser) Valid() bool {

	f.Failures = make(map[string]string)

	if strings.TrimSpace(f.Email) == "" {
		f.Failures["Email"] = "Email is required"
	}
	if strings.TrimSpace(f.Password) == "" {
		f.Failures["Password"] = "Password is required"
	}

	return len(f.Failures) == 0
}
