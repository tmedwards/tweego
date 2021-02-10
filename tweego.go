/*
	tweego (a twee compiler in Go)

	Copyright © 2014–2021 Thomas Michael Edwards. All rights reserved.
	Use of this source code is governed by a Simplified BSD License which
	can be found in the LICENSE file.
*/

package main

import (
	"log"
)

const tweegoName = "tweego"

func init() {
	// Clear standard logger flags.
	log.SetFlags(0)
}

func main() {
	// Create a new config instance.
	c := newConfig()

	// Build the output and, possibly, log various stats.
	if c.watchFiles {
		buildName := relPath(c.outFile)
		paths := append(c.sourcePaths, c.modulePaths...)
		watchFilesystem(paths, c.outFile, func() {
			log.Printf("BUILDING: %s", buildName)
			buildOutput(c)
		})
	} else {
		buildOutput(c)

		// Logging.
		if c.logFiles {
			log.Println()
			statsLogFiles()
			log.Println()
		}
		if c.logStats {
			if !c.logFiles {
				log.Println()
			}
			statsLog()
			log.Println()
		}
	}
}

func buildOutput(c *config) *story {
	// Get the source and module paths.
	sourcePaths := getFilenames(c.sourcePaths, c.outFile)
	modulePaths := getFilenames(c.modulePaths, c.outFile)

	// Create a new story instance and load the source files.
	s := newStory()
	s.load(sourcePaths, c)

	// Finalize the config with values from the `StoryData` passage, if any.
	c.mergeStoryConfig(s)

	// Write the output.
	switch c.outMode {
	case outModeJSON:
		// Write out the project as JSON.
		if _, err := fileWriteAll(c.outFile, alignRecordSeparators(s.toJSON(c.startName))); err != nil {
			log.Fatalf(`error: %s`, err.Error())
		}
	case outModeTwee3, outModeTwee1:
		// Write out the project as Twee source.
		if _, err := fileWriteAll(c.outFile, alignRecordSeparators(s.toTwee(c.outMode))); err != nil {
			log.Fatalf(`error: %s`, err.Error())
		}
	case outModeTwine2Archive:
		// Write out the project as Twine 2 archived HTML.
		if _, err := fileWriteAll(c.outFile, s.toTwine2Archive(c.startName)); err != nil {
			log.Fatalf(`error: %s`, err.Error())
		}
	case outModeTwine1Archive:
		// Write out the project as Twine 1 archived HTML.
		if _, err := fileWriteAll(c.outFile, s.toTwine1Archive(c.startName)); err != nil {
			log.Fatalf(`error: %s`, err.Error())
		}
	default:
		// Basic sanity checks.
		if !s.has(c.startName) {
			log.Fatalf("error: Starting passage %q not found.", c.startName)
		}
		if (s.format.isTwine1Style() || s.name == "") && !s.has("StoryTitle") {
			log.Fatal(`error: Special passage "StoryTitle" not found.`)
		}

		if s.format.isTwine2Style() {
			// Write out the project as Twine 2 compiled HTML.
			if _, err := fileWriteAll(
				c.outFile,
				modifyHead(
					s.toTwine2HTML(c.startName),
					modulePaths,
					c.headFile,
					c.encoding,
				),
			); err != nil {
				log.Fatalf(`error: %s`, err.Error())
			}
		} else {
			// Write out the project as Twine 1 compiled HTML.
			if _, err := fileWriteAll(
				c.outFile,
				modifyHead(
					s.toTwine1HTML(c.startName),
					modulePaths,
					c.headFile,
					c.encoding,
				),
			); err != nil {
				log.Fatalf(`error: %s`, err.Error())
			}
		}
	}

	return s
}
