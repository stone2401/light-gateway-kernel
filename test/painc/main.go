package main

import (
	"fmt"
	"reflect"
)

type StructTest struct {
	Field  string
	Field2 string
	Field3 string
	Chi    ChiStructTest
}

type ChiStructTest struct {
	ChiField  string
	ChiField2 string
	ChiField3 string
}

func main() {
	my := "石振飞"
	value := reflect.ValueOf(my)
	fmt.Println("name", value.Kind())
}
