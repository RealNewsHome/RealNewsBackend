package main

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	got, err := HashPassword("testcat")
	want := "$2a$14$YuaRwxCAtU61IYylZh8hgeQA8J5k7gJSXttrI0N.W/RstM6wFm0Ye"
	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}

	if err != nil {
		t.Errorf("an error occurred in test")
	}
}

func TestAuthenticate(t *testing.T) {
	got, err := Authenticate("myfavoritecat@feet.com", "testcat")
	if len(got) > 0 {
		t.Errorf("no token returned, instead got: %v", err)
	}
}
