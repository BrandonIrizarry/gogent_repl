package main

import (
	"log"

	"github.com/BrandonIrizarry/gogent"
	"github.com/joho/godotenv"
)

func main() {
	// Load our environment variables (including the Gemini API
	// key.)
	//
	// Note that, since we don't have our custom logger yet, we're
	// using the default logger for now.
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	var g gogent.Gogent

	_, err := g.Write([]byte("What is the name of this program?"))
	if err != nil {
		log.Println(err)
	}

	resp := make([]byte, 10000)
	_, err = g.Read(resp)
	if err != nil {
		log.Println(err)
	}

	log.Println(string(resp))
}
