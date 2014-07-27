package main

import (
	"fmt"
	"os"
)

func main() {
	policy := NewPolicy()

	if len(os.Args) < 2 || os.Args[1] != "curl" {
		fmt.Println(policy.ToJsonResponse())
	} else {
		fmt.Println(policy.ToCurl())
	}
}
