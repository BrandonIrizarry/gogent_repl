package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/BrandonIrizarry/gogent"
	"github.com/BrandonIrizarry/gogent_repl/internal/cliargs"
	"github.com/joho/godotenv"
)

// A list of models we can conveniently select from.
const (
	GeminiTwoPointFiveFlash            = "gemini-2.5-flash"
	GeminiTwoFlash                     = "gemini-2.0-flash"
	GeminiTwoPointFiveFlashLite        = "gemini-2.5-flash-lite"
	GeminiTwoPointFiveFlashLitePreview = "gemini-2.5-flash-lite-preview-09-2025"
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

	cliArgs, err := cliargs.New()
	if err != nil {
		log.Fatal(err)
	}

	g := gogent.Gogent{
		WorkingDir:    cliArgs.WorkingDir,
		MaxFilesize:   100_000,
		MaxIterations: 20,
		LLMModel:      GeminiTwoPointFiveFlashLite,
		Debug:         cliArgs.Debug,
	}

	ask, err := g.Init()
	if err != nil {
		log.Fatal(err)
	}

	// The REPL loop.
	for {
		prompt, quit := getPrompt()
		if quit {
			break
		}

		response, err := ask(prompt)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(response)
	}

	fmt.Printf("Token counts: %+v\n", g.TokenCounts())
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
