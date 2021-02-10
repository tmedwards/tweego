/*
	Copyright © 2014–2021 Thomas Michael Edwards. All rights reserved.
	Use of this source code is governed by a Simplified BSD License which
	can be found in the LICENSE file.
*/

package main

import (
	// standard packages
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	// external packages
	"github.com/radovskyb/watcher"
)

var programDir string
var workingDir string

func init() {
	// Attempt to get the canonical program directory, failure is okay.
	if pp, err := os.Executable(); err == nil {
		if pd, err := filepath.EvalSymlinks(filepath.Dir(pp)); err == nil {
			programDir = pd
		}
	}

	// Attempt to get the canonical working directory, failure is okay.
	if wd, err := os.Getwd(); err == nil {
		if wd, err := filepath.EvalSymlinks(wd); err == nil {
			workingDir = wd
		}
	}
}

var errNoOutToIn = fmt.Errorf("no output to input source")

// Walk the specified pathnames, collecting regular files.
func getFilenames(pathnames []string, outFilename string) []string {
	var (
		filenames  []string
		absOutFile string
	)
	var fileWalker filepath.WalkFunc = func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		// Get the absolute path.
		absolute, err := absPath(path)
		if err != nil {
			return err
		}

		// The output file must not be one of the input files.
		if absolute == absOutFile {
			return errNoOutToIn
		}

		// Return a relative path if one can be resolved and it is shorter than the
		// absolute path, failure is okay.
		if relative, err := filepath.Rel(workingDir, absolute); err == nil {
			if relative, err := filepath.EvalSymlinks(relative); err == nil {
				if len(relative) < len(absolute) {
					filenames = append(filenames, relative)
					return nil
				}
			}
		}

		filenames = append(filenames, absolute)
		return nil
	}

	// Get the absolute output filename.
	if abs, err := absPath(outFilename); err == nil {
		absOutFile = abs
	} else {
		log.Fatalf("error: path %s: %s", outFilename, err.Error())
	}

	// Walk the pathnames.
	for _, pathname := range pathnames {
		if pathname == "-" {
			log.Print("warning: path -: Reading from standard input is unsupported.")
			continue
		} else if err := filepath.Walk(pathname, fileWalker); err != nil {
			if err == errNoOutToIn {
				log.Fatalf("error: path %s: Output file cannot be an input source.", pathname)
			} else {
				log.Printf("warning: path %s: %s", pathname, err.Error())
				continue
			}
		}
	}

	return filenames
}

// Watch the specified pathnames, calling the build callback as necessary.
func watchFilesystem(pathnames []string, outFilename string, buildCallback func()) {
	var (
		buildRate = time.Millisecond * 500
		pollRate  = buildRate * 2
	)

	// Create a new watcher instance.
	w := watcher.New()

	// Only notify on certain events.
	w.FilterOps(
		watcher.Create,
		watcher.Write,
		watcher.Remove,
		watcher.Rename,
		watcher.Move,
	)

	// Ignore the output file.
	w.Ignore(outFilename)

	// Start a goroutine to handle the event loop.
	go func() {
		build := false

		for {
			select {
			case <-time.After(buildRate):
				if build {
					buildCallback()
					build = false
				}
			case event := <-w.Event:
				if event.FileInfo != nil {
					isDir := event.IsDir()

					if event.Op == watcher.Write && isDir {
						continue
					}

					var pathname string
					switch event.Op {
					case watcher.Move, watcher.Rename:
						pathname = fmt.Sprintf("%s -> %s", relPath(event.OldPath), relPath(event.Path))
						if !build && !isDir {
							build = knownFileType(event.OldPath) || knownFileType(event.Path)
						}
					default:
						pathname = relPath(event.Path)
						if !build && !isDir {
							build = knownFileType(event.Path)
						}
					}
					log.Printf("%s: %s", event.Op, pathname)
				}
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	// Recursively watch the specified paths for changes.
	for _, pathname := range pathnames {
		if err := w.AddRecursive(pathname); err != nil {
			log.Fatalln(err)
		}
	}

	// Print a message telling the user how to cancel watching
	// and list all paths being watched.
	log.Print()
	log.Print("Watch mode started.  Press CTRL+C to stop.")
	log.Print()
	log.Printf("Recursively watched paths: %d", len(pathnames))
	for _, pathname := range pathnames {
		log.Printf("  %s", relPath(pathname))
	}
	log.Print()

	// Build the ouput once before the watcher starts.
	buildCallback()

	// Start watching.
	if err := w.Start(pollRate); err != nil {
		log.Fatalln(err)
	}
}

func absPath(original string) (string, error) {
	var absolute string

	if abs, err := filepath.Abs(original); err == nil {
		absolute = abs
		if abs, err := filepath.EvalSymlinks(absolute); err == nil {
			absolute = abs
		} else {
			dir, base := filepath.Split(absolute)
			if abs, err := filepath.EvalSymlinks(dir); err == nil {
				absolute = filepath.Join(abs, base)
			} else {
				return "", err
			}
		}
	} else {
		return "", err
	}

	return absolute, nil
}

func relPath(original string) string {
	absolute, err := absPath(original)
	if err != nil {
		// Failure is okay, just return the original path.
		return original
	}

	if relative, err := filepath.Rel(workingDir, absolute); err == nil {
		if relative, err := filepath.EvalSymlinks(relative); err == nil {
			if len(relative) < len(absolute) {
				return relative
			}
		}
	}

	// Failure is okay, just return the absolute path.
	return absolute
}

func knownFileType(filename string) bool {
	switch normalizedFileExt(filename) {
	// NOTE: The case values here should match those in `storyload.go:(*story).load()`.
	case "tw", "twee",
		"tw2", "twee2",
		"htm", "html",
		"css",
		"js",
		"otf", "ttf", "woff", "woff2",
		"gif", "jpeg", "jpg", "png", "svg", "tif", "tiff", "webp",
		"aac", "flac", "m4a", "mp3", "oga", "ogg", "opus", "wav", "wave", "weba",
		"mp4", "ogv", "webm",
		"vtt":
		return true
	}

	return false
}
