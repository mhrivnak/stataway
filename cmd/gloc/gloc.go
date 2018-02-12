package main

import (
	"fmt"
	"os"
	"github.com/mhrivnak/stataway/pkg/gloc"
)

func main() {
	username := os.Getenv("username")
	password := os.Getenv("password")

	fmt.Printf("authenticating as %s with %s\n", username, password)

	err := gloc.Get(username, password)
	if err != nil {
		fmt.Println("GOT AN ERROR")
		fmt.Println(err.Error())
	}
}
