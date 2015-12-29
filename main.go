package main

import "os"
import "fmt"

import "github.com/jwiklund/tmpl/tmpl"

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Usage: tmpl SOURCE")
		return
	}
	template, err := tmpl.New(os.Args[1])
	if err != nil {
		fmt.Printf("Failed to read template %s\n", err.Error())
		return
	}
	if err = template.Create(); err != nil {
		fmt.Printf("Failed to create template %s\n", err.Error())
		return
	}
	fmt.Println("Successfully created templated project")
}
