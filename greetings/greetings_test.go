package greetings

import (
	"regexp"
	"testing"
)

func TestHelloName(t *testing.T) {
	name := "Bandit"
	want := regexp.MustCompile(`\b` + name + `\b`)
	res, err := Hello(name)
	if !want.MatchString(res) || err != nil {
		t.Fatalf(`Hello("Bandit") = %q, %v, want match for %#q, nil`, res, err, want)
	}
}

func TestHelloEmpty(t *testing.T) {
	res, err := Hello("")
	if res != "" || err == nil {
		t.Fatalf(`Hello("") = %q, %v, want "", error`, res, err)
	}
}
