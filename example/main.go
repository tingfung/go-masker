package main

import (
	"fmt"

	masker "github.com/syoya/go-masker"
)

func main() {
	m, err := masker.New(map[string]string{
		"password": "**********",
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", m.Mask([]byte(`{"email":"foo@example.com","password":"p@ssw0rd"}`))) // -> {"email":"foo@example.com","password":"**********"}
}
