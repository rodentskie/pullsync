package mapstruct

import (
	"reflect"
	"testing"

	"github.com/go-faker/faker/v4"
)

type SampleStruct struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

func TestMapToStruct(t *testing.T) {

	randomInt, _ := faker.RandomInt(100, 200)

	id := randomInt[0]
	name := faker.Name()
	age := randomInt[1]
	email := faker.Email()
	data := map[string]interface{}{
		"id":    id,
		"name":  name,
		"age":   age,
		"email": email,
	}

	var target SampleStruct

	err := MapToStruct(data, &target)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	if target.ID != id {
		t.Errorf("Expected ID to be %d, but got %d", id, target.ID)
	}
	if target.Name != name {
		t.Errorf("Expected Name to be '%s', but got '%s'", name, target.Name)
	}
	if target.Age != age {
		t.Errorf("Expected Age to be %d, but got %d", age, target.Age)
	}
	if target.Email != email {
		t.Errorf("Expected Email to be '%s', but got '%s'", email, target.Email)
	}
}

func TestStructToMap(t *testing.T) {
	sample := SampleStruct{
		ID:    1,
		Name:  "John",
		Age:   30,
		Email: "john@example.com",
	}

	result := StructToMap(sample)

	expected := map[string]interface{}{
		"id":    1,
		"name":  "John",
		"age":   30,
		"email": "john@example.com",
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}
