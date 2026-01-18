package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

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

	// The source of the configuration of the various fields is up
	// to the particular frontend to decide (e.g. CLI arguments,
	// YAML file, TUI/GUI widget, etc.)
	g := gogent.Gogent{
		WorkingDir:    ".",
		MaxFilesize:   100_000,
		MaxIterations: 20,
		LLMModel:      "gemini-2.5-flash-lite-preview-09-2025",
	}

	// The REPL loop.
	for {
		prompt, quit := getPrompt()
		if quit {
			break
		}

		response, err := g.Ask(prompt)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(response)
	}
}

func getPrompt() (string, bool) {
	fmt.Println()
	fmt.Println("Ask the agent something (press Enter twice to submit your prompt)")
	fmt.Println("Submit a blank prompt to exit")
	fmt.Print("> ")

	scanner := bufio.NewScanner(os.Stdin)
	var bld strings.Builder

	for scanner.Scan() {
		text := scanner.Text()

		if strings.TrimSpace(text) == "" {
			break
		}

		// Write an extra space, to make sure that words
		// across newline boundaries don't run on to each
		// other.
		bld.WriteString(" ")
		bld.WriteString(text)
	}

	// Nothing was written, meaning we must signal to our caller
	// to not invoke the agent REPL.
	if bld.Len() == 0 {
		fmt.Println("Bye!")
		return "", true
	}

	fmt.Println("Thinking...")
	return bld.String(), false
}
