package main

import (
	"flag"
	"fmt"

	"github.com/hardsky/transforma/internal/transforma"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		args = []string{"."}
	}

	err := transforma.Generate(args[0])
	if err != nil {
		fmt.Println(err)
	}
}
