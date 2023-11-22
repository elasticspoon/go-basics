package main

import (
	"fmt"
	"log"

	"example.com/greetings"
	"golang.org/x/example/hello/reverse"
)

func main() {
	log.SetPrefix("greetings: ")
	log.SetFlags(0)

	names := []string{"Bandit", "Bella", "Buddy"}

	messages, err := greetings.Hellos(names)
	if err != nil {
		log.Fatal(err)
	}

	for _, message := range messages {
		fmt.Println(reverse.String(message))
	}
	fmt.Println(reverse.String("hello"), reverse.Int(24601))
}
