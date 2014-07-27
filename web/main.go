package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] != "curl" {
		Serve()
	} else {
		fmt.Println(NewPolicy().ToCurl())
	}
}
