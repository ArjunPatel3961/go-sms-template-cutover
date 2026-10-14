package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("usage: go run . <release> <signature> <template>")
		os.Exit(2)
	}
	c, err := NewInfraiClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	r := Release{Name: os.Args[1], Signature: os.Args[2], Template: os.Args[3]}
	if err := PublishRelease(c, r); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("release %q published; rollback by restoring the incumbent route\n", r.Name)
}
