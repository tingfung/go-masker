package main

import (
	"fmt"

	masker "github.com/tingfung/go-masker"
)

func main() {
	m, err := masker.New(masker.Options{
		Replacement: map[string]string{
			"password": "***",
		},
		Truncation: masker.Truncation{
			Length:   20,
			Omission: "...",
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", m.Mask([]byte(`{"Password":"p@Ssw0rd","long":"this should be truncated"}`)))
	//                             -> {"Password":"***","long":"this should be trunc..."}
}
