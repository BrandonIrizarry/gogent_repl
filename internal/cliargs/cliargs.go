package cliargs

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type cliArguments struct {
	WorkingDir string
}

func New() (cliArguments, error) {
	var cliArgs cliArguments

	flag.StringVar(&cliArgs.WorkingDir, "dir", "", "Set project directory")
	flag.Parse()

	// Handle the 'dir' argument.
	hdir, err := os.UserHomeDir()
	if err != nil {
		return cliArguments{}, err
	}

	// If -dir is omitted, flag this to the caller so that it can
	// then obtain the working directory from a radio selection
	// widget, or something similar.
	if cliArgs.WorkingDir == "" {
		return cliArgs, nil
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
		// Remove some of the gobbledygook from the error
		// otherwise returned here.
		var pathErr *os.PathError
		if errors.As(err, &pathErr) {
			return cliArguments{}, fmt.Errorf("%s: %v", pathErr.Path, pathErr.Err)
		}

		// Else, for now, forward the error onward with a
		// meaningful preface.
		return cliArguments{}, fmt.Errorf("invalid working directory argument: %w", err)
	}

	if !info.IsDir() {
		return cliArguments{}, errors.New("Working directory argument not a directory")
	}

	return cliArgs, nil
}
