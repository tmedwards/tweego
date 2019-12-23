/*
	Copyright © 2014–2019 Thomas Michael Edwards. All rights reserved.
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
	"strings"
	"time"
	// external packages
	"github.com/radovskyb/watcher"
)

var programDir string
var workingDir string

func init() {
	// Attempt to get the program directory, failure is okay.
	if pp, err := os.Executable(); err == nil {
		programDir = filepath.Dir(pp)
	}

	// Attempt to get the working directory, failure is okay.
	if wd, err := os.Getwd(); err == nil {
		workingDir = wd
	}
}

var noOutToIn = fmt.Errorf("no output to input source")

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

		absolute, err := filepath.Abs(path)
		if err != nil {
			return err
		}
		if absolute == absOutFile {
			return noOutToIn
		}
		relative, _ := filepath.Rel(workingDir, absolute) // Failure is okay.
		if relative != "" {
			filenames = append(filenames, relative)
		} else {
			filenames = append(filenames, absolute)
		}
		return nil
	}

	// Get the absolute output filename.
	absOutFile, err := filepath.Abs(outFilename)
	if err != nil {
		log.Fatalf("error: path %s: %s", outFilename, err.Error())
	}

	for _, pathname := range pathnames {
		if pathname == "-" {
			log.Print("warning: path -: Reading from standard input is unsupported.")
			continue
		} else if err := filepath.Walk(pathname, fileWalker); err != nil {
			if err == noOutToIn {
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
						// NOTE: Format of Move/Rename event `Path` field: "oldName -> newName".
						// TODO: Should probably error out if we can't split the event.Path value.
						names := strings.Split(event.Path, " -> ")
						pathname = fmt.Sprintf("%s -> %s", relPath(names[0]), relPath(names[1]))
						if !build && !isDir {
							build = knownFileType(names[0]) || knownFileType(names[1])
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

func relPath(original string) string {
	absolute, err := filepath.Abs(original)
	if err != nil {
		// Failure is okay, just return the original path.
		return original
	}

	relative, err := filepath.Rel(workingDir, absolute)
	if err != nil {
		// Failure is okay, just return the absolute path.
		return absolute
	}

	return relative
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
