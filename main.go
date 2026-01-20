package main

import (
	"fmt"
	"log"

	"github.com/BrandonIrizarry/gogent"
	"github.com/BrandonIrizarry/gogent_repl/internal/cliargs"
	"github.com/BrandonIrizarry/gogent_repl/internal/promptbox"
	"github.com/charmbracelet/glamour"
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
		LogLevel:      cliArgs.LogLevel,
	}

	ask, err := g.Init()
	if err != nil {
		log.Fatal(err)
	}

	// The REPL loop.
	for {
		prompt, err := promptbox.GetPrompt()
		if err != nil {
			log.Fatal(err)
		}

		if prompt == "" {
			break
		}

		responseText, err := ask(prompt)
		if err != nil {
			log.Fatal(err)
		}

		if glamourText, err := glamour.Render(responseText, "light"); err != nil {
			log.Println("Glamour rendering failed, defaulting to plain text")
			fmt.Println(responseText)
		} else {
			fmt.Println(glamourText)
		}
	}

	fmt.Println("Bye!")
	fmt.Printf("Token counts: %+v\n", g.TokenCounts())
}
