package main

import (
	"testing"
)

func TestProcessUser(t *testing.T) {
	user := UserInput{
		FirstName: "Ronen  ",
		LastName:  "  Konfortes",
		Email:     "  konfortes@gmail.com  ",
	}

	expected := UserInput{
		FirstName: "Ronen",
		LastName:  "Konfortes",
		Email:     "konfortes@gmail.com",
	}

	processUser(&user)

	if expected.FirstName != user.FirstName {
		t.Errorf("expected: %s. got: %s", expected.FirstName, user.FirstName)
	}

	if expected.LastName != user.LastName {
		t.Errorf("expected: %s. got: %s", expected.LastName, user.LastName)
	}

	if expected.Email != user.Email {
		t.Errorf("expected: %s. got: %s", expected.Email, user.Email)
	}
}
