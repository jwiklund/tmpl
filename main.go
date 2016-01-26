package main

import "os"
import "fmt"

import "github.com/jwiklund/tmpl/tmpl"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: tmpl SOURCE [PROPERTY=VALUE]")
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
	env, err := tmpl.GetEnvironment(template, os.Args[2:]...)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if err = template.Create(env, target); err != nil {
		fmt.Printf("Failed to create template %s\n", err.Error())
		return
	}
	fmt.Println("Successfully created templated project")
}
