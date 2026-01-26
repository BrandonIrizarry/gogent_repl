package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"time"

	"github.com/BrandonIrizarry/gogent"
	"github.com/BrandonIrizarry/gogent_repl/internal/cliargs"
	"github.com/BrandonIrizarry/gogent_repl/internal/promptbox"
	"github.com/BrandonIrizarry/gogent_repl/internal/radioselect"
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
	// Load a log file to write logs to as well.
	logFilename := fmt.Sprintf(".logs-%d", time.Now().Unix())
	logFile, err := os.OpenFile(logFilename, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	log.SetOutput(io.MultiWriter(os.Stdout, logFile))

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

	historyFile, err := os.OpenFile(".history", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer historyFile.Close()

	// Load these up front unconditionally, since we either use
	// them as our selections, or else use them to avoid writing a
	// duplicate history entry.
	var choices []string
	scanner := bufio.NewScanner(historyFile)
	for scanner.Scan() {
		choices = append(choices, scanner.Text())
	}

	// If -dir wasn't provided, present the radio-button selection
	// widget to the user.
	if cliArgs.WorkingDir == "" {
		// The choices slice should have something in it, else
		// SelectWorkingDir will panic with an out-of-bounds
		// access to its counterpart inside the TUI model.
		if len(choices) == 0 {
			log.Fatal("-dir missing, and no saved choices inside history file")
		}

		wdir, err := radioselect.LoadList("Recent projects", choices, choices[0])
		if err != nil {
			log.Fatal(err)
		}

		// An empty wdir result means that the user ctrl+c'ed
		// out of the radio selection widget; for now, simply
		// quit the application.
		if wdir == "" {
			fmt.Println("Bye!")
			os.Exit(0)
		}

		cliArgs.WorkingDir = wdir
	} else {
		if !slices.Contains(choices, cliArgs.WorkingDir) {
			historyFile.WriteString(cliArgs.WorkingDir + "\n")
		}
	}

	// Ask the user for the model they wish to use
	modelNames := []string{
		GeminiTwoFlash,
		GeminiTwoPointFiveFlash,
		GeminiTwoPointFiveFlashLite,
		GeminiTwoPointFiveFlashLitePreview,
	}

	// The user's selection from the previous session is found in
	// './default-model'; use it as the default list selection
	// here.
	defaultModelFile, err := os.OpenFile(".default-model", os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer defaultModelFile.Close()

	var defaultModel string
	buf, err := io.ReadAll(defaultModelFile)
	if err != nil {
		log.Fatal(err)
	}

	if len(buf) == 0 {
		defaultModel = modelNames[0]
	} else {
		defaultModel = string(buf)
	}

	modelName, err := radioselect.LoadList("Model name", modelNames, defaultModel)
	if err != nil {
		log.Fatal(err)
	}

	defaultModelFile.WriteString(modelName)

	g := gogent.Gogent{
		WorkingDir:    cliArgs.WorkingDir,
		MaxFilesize:   100_000,
		MaxIterations: 20,
		LLMModel:      modelName,
		LogLevel:      cliArgs.LogLevel,
	}

	// Initialize Gogent. This also loads Gogent's "ask" function,
	// which is what we use to drive our conversation with the
	// agent.
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

		var responseText string

		responseText, err = ask(prompt)
		if err != nil {
			log.Fatal(err)
		}

		if responseText, err = glamour.Render(responseText, "light"); err != nil {
			log.Println("Glamour rendering failed, defaulting to plain text")
		}

		fmt.Println(responseText)
	}

	fmt.Println("Bye!")
	fmt.Printf("Token counts: %+v\n", g.TokenCounts())
}
