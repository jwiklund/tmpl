package main

import "os"
import "fmt"

import "github.com/jwiklund/tmpl/tmpl"

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Usage: tmpl SOURCE")
		return
	}
	template, err := tmpl.NewTemplate(os.Args[1])
	if err != nil {
		fmt.Printf("Failed to read template %s\n", err.Error())
		return
	}
	target, err := tmpl.GetTarget(".")
	if err != nil {
		fmt.Println("Failed to create target %s\n", err.Error())
		return
	}
	if err = template.Create(target); err != nil {
		fmt.Printf("Failed to create template %s\n", err.Error())
		return
	}
	fmt.Println("Successfully created templated project")
}
