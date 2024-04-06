package main

import (
	"encoding/json"
	"fmt"
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
	s := StructTest{
		Field:  "field",
		Field2: "field2",
		Field3: "field3",
		Chi: ChiStructTest{
			ChiField:  "chiField",
			ChiField2: "chiField2",
			ChiField3: "chiField3",
		},
	}
	data, _ := json.MarshalIndent(s, "", "  ")
	// zlog.Zlog().Info("test", zap.String("data", string(data)))
	fmt.Println(string(data))
}
