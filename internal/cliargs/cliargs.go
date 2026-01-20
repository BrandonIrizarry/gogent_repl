package cliargs

import (
	"errors"
	"flag"
	"os"
	"path/filepath"
	"strings"
)

type cliArguments struct {
	WorkingDir string
	LogLevel   string
}

func New() (cliArguments, error) {
	var cliArgs cliArguments

	flag.StringVar(&cliArgs.WorkingDir, "dir", ".", "Set project directory")
	flag.StringVar(&cliArgs.LogLevel, "log", "error", "Set log level to one of debug, info, warn, or error")
	flag.Parse()

	// Handle the 'dir' argument.
	hdir, err := os.UserHomeDir()
	if err != nil {
		return cliArguments{}, err
	}

	cliArgs.WorkingDir, err = filepath.Abs(cliArgs.WorkingDir)
	if err != nil {
		return cliArguments{}, err
	}

	if !strings.HasPrefix(cliArgs.WorkingDir, hdir) {
		return cliArguments{}, errors.New("Working directory not inside user's $HOME")
	}

	info, err := os.Stat(cliArgs.WorkingDir)
	if err != nil {
		return cliArguments{}, err
	}

	if !info.IsDir() {
		return cliArguments{}, errors.New("Working directory argument not a directory")
	}

	return cliArgs, nil
}
