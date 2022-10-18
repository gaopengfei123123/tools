package convert

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestJsonCamelCase_MarshalJSON(t *testing.T) {
	type Person struct {
		HelloWold       string
		LightWeightBaby string
	}
	var a = Person{HelloWold: "GPF", LightWeightBaby: "muscle"}
	res, _ := json.Marshal(JsonCamelCase{Value: a})
	fmt.Printf("%s", res)
}

func TestJsonSnakeCase_MarshalJSON(t *testing.T) {
	type Person struct {
		HelloWold       string
		LightWeightBaby string
	}
	var a = Person{HelloWold: "GPF", LightWeightBaby: "muscle"}
	res, _ := json.Marshal(JsonSnakeCase{Value: a})
	fmt.Printf("%s", res)
}
