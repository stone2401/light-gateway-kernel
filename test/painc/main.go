package main

import (
	"fmt"
	"regexp"
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
	// my := "石振飞"
	// value := reflect.ValueOf(my)
	// fmt.Println("name", value.Kind())
	oldUrl := "/ai/v1/user"
	newUrl := `/$1`
	re := regexp.MustCompile(`^.+api/?`)
	fmt.Printf("re.MatchString(oldUrl): %v\n", re.MatchString(oldUrl))
	newPath := re.ReplaceAllString(oldUrl, newUrl)
	fmt.Println(newPath)
	// fmt.Printf("strings.Replace(oldUrl, \"^.+api/?\", newUrl, 1): %v\n", strings.Replace(oldUrl, "^.+api/?", newUrl, 1))
}
