package main

import (
	"github.com/gorilla/schema"
	"testing"
)

type Person struct {
	Name  string
	Phone string
}

func TestSchema(t *testing.T) {
	values := map[string][]string{
		"Name":  {"John"},
		"Phone": {"999-999-999"},
	}
	person := new(Person)
	decoder := schema.NewDecoder()
	decoder.Decode(person, values)
	t.Log(person)
}
