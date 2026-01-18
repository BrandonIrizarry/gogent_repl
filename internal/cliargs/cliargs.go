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
	Debug      bool
}

func New() (cliArguments, error) {
	var cliArgs cliArguments

	flag.StringVar(&cliArgs.WorkingDir, "dir", ".", "Set project directory")
	flag.BoolVar(&cliArgs.Debug, "debug", false, "Whether agent should write debugging logs")
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
