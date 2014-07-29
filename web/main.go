package main

import (
	"fmt"
	"github.com/bjeanes/hk-deploy/policy"
	"os"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] != "curl" {
		Serve()
	} else {
		fmt.Println(policy.NewPolicy().ToCurl())
	}
}
