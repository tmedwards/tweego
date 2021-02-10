/*
	Copyright © 2014–2021 Thomas Michael Edwards. All rights reserved.
	Use of this source code is governed by a Simplified BSD License which
	can be found in the LICENSE file.
*/

package main

import (
	"log"
)

// statistics are collected statistics about the compiled story.
type statistics struct {
	files struct {
		project  []string // Project source files.
		external []string // Both modules and the head file.
	}
	counts struct {
		passages      uint64 // Count of all passages.
		storyPassages uint64 // Count of story passages.
		storyWords    uint64 // Count of story passage "words" (typing measurement style).
	}
}

var stats = statistics{}

func statsAddProjectFile(filepath string) {
	stats.files.project = append(stats.files.project, filepath)
}

func statsAddExternalFile(filepath string) {
	stats.files.external = append(stats.files.external, filepath)
}

func statsLog() {
	log.Print("Statistics")
	log.Printf("  Total> Passages: %d", stats.counts.passages)
	log.Printf("  Story> Passages: %d, Words: %d", stats.counts.storyPassages, stats.counts.storyWords)
}

func statsLogFiles() {
	log.Println("Processed files (in order)")
	log.Printf("  Project files: %d", len(stats.files.project))
	for _, file := range stats.files.project {
		log.Printf("    %s", file)
	}
	log.Printf("  External files: %d", len(stats.files.external))
	for _, file := range stats.files.external {
		log.Printf("    %s", file)
	}
}
