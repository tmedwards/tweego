/*
	Copyright © 2014–2021 Thomas Michael Edwards. All rights reserved.
	Use of this source code is governed by a Simplified BSD License which
	can be found in the LICENSE file.
*/

package main

import (
	"errors"
	"os"
	"os/user"
	"runtime"
)

func userHomeDir() (string, error) {
	// Prefer the user's `HOME` environment variable.
	if homeDir := os.Getenv("HOME"); homeDir != "" {
		return homeDir, nil
	}

	// Elsewise, use the user's `.HomeDir` info.
	if curUser, err := user.Current(); err == nil && curUser.HomeDir != "" {
		return curUser.HomeDir, nil
	}

	// Failovers for Windows, though they should be unnecessary in Go ≥v1.7.
	if runtime.GOOS == "windows" {
		// Prefer the user's `USERPROFILE` environment variable.
		if homeDir := os.Getenv("USERPROFILE"); homeDir != "" {
			return homeDir, nil
		}

		// Elsewise, use the user's `HOMEDRIVE` and `HOMEPATH` environment variables.
		homeDrive := os.Getenv("HOMEDRIVE")
		homePath := os.Getenv("HOMEPATH")
		if homeDrive != "" && homePath != "" {
			return homeDrive + homePath, nil
		}
	}

	return "", errors.New("Cannot find user home directory.")
}
