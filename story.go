/*
	Copyright © 2014–2021 Thomas Michael Edwards. All rights reserved.
	Use of this source code is governed by a Simplified BSD License which
	can be found in the LICENSE file.
*/

package main

import (
	"fmt"
	"log"
	"strings"
)

// Twine 1 story metadata.
type twine1Metadata struct {
	// WARNING: Do not use individual fields here as story formats are allowed to
	// define their own `StorySettings` pairs—via a custom header file—so we have
	// no way of knowing what the keys might be prior to parsing the passage.
	settings map[string]string // Map of `StorySettings` key/value pairs.
}

// Twine 2 story metadata.
type twine2OptionsMap map[string]bool
type twine2TagColorsMap map[string]string
type twine2Metadata struct {
	format        string             // Name of the story format.
	formatVersion string             // SemVer of the story format.
	options       twine2OptionsMap   // Map of option-name/bool pairs.
	start         string             // Name of the starting passage.
	tagColors     twine2TagColorsMap // Unused by Tweego.  Map of tag-name/color pairs.
	zoom          float64            // Unused by Tweego.  Zoom level.  Why is this even a part of the story metadata?  It's editor configuration.
}

// Core story data.
type story struct {
	name     string
	ifid     string // A v4 random UUID, see: https://ifdb.tads.org/help-ifid.
	passages []*passage

	// Legacy fields from Tweego v1 StorySettings.
	legacyIFID string

	// Twine 1 & 2 compiler metadata.
	twine1 twine1Metadata
	twine2 twine2Metadata

	// Tweego compiler internals.
	format    *storyFormat
	processed map[string]bool
}

// newStory creates a new story instance.
func newStory() *story {
	return &story{
		passages: make([]*passage, 0, 64), // Initially create enough space for 64 passages.
		twine1: twine1Metadata{
			settings: make(map[string]string),
		},
		twine2: twine2Metadata{
			options:   make(twine2OptionsMap),
			tagColors: make(twine2TagColorsMap),
			zoom:      1,
		},
		processed: make(map[string]bool),
	}
}

func (s *story) count() int {
	return len(s.passages)
}

func (s *story) has(name string) bool {
	for _, p := range s.passages {
		if p.name == name {
			return true
		}
	}
	return false
}

func (s *story) index(name string) int {
	for i, p := range s.passages {
		if p.name == name {
			return i
		}
	}
	return -1
}

func (s *story) get(name string) (*passage, error) {
	for _, p := range s.passages {
		if p.name == name {
			return p, nil
		}
	}
	return nil, fmt.Errorf("get %s: No such passage.", name)
}

func (s *story) deleteAt(i int) error {
	upper := len(s.passages) - 1
	if 0 > i || i > upper {
		return fmt.Errorf("deleteAt %d: Index out of range.", i)
	}

	// TODO: Should the `copy()` only occur if `i < upper`?
	copy(s.passages[i:], s.passages[i+1:]) // shift elements down by one to overwrite the original element

	s.passages[upper] = nil         // zero the last element, which was itself duplicated by the last operation
	s.passages = s.passages[:upper] // reslice to remove the last element
	return nil
}

func (s *story) append(p *passage) {
	// Append the passage if new, elsewise replace the existing version.
	if i := s.index(p.name); i == -1 {
		s.passages = append(s.passages, p)
		stats.counts.passages++
		if p.isStoryPassage() {
			stats.counts.storyPassages++
			stats.counts.storyWords += p.countWords()
		}
	} else {
		log.Printf("warning: Replacing existing passage %q with duplicate.", p.name)
		s.passages[i] = p
	}
}

func (s *story) prepend(p *passage) {
	// Prepend the passage if new, elsewise replace the existing version.
	if i := s.index(p.name); i == -1 {
		s.passages = append([]*passage{p}, s.passages...)
		stats.counts.passages++
		if p.isStoryPassage() {
			stats.counts.storyPassages++
			stats.counts.storyWords += p.countWords()
		}
	} else {
		log.Printf("warning: Replacing existing passage %q with duplicate.", p.name)
		s.passages[i] = p
	}
}

func (s *story) add(p *passage) {
	// Preprocess compiler-oriented special passages.
	switch p.name {
	case "StoryIncludes":
		/*
			NOTE: StoryIncludes is a compiler special passage for Twine 1.4,
			and apparently Twee2.  Twee 1.4 does not support it—likely for
			the same reasons Tweego will not (outlined below).

			You may specify an arbitrary number of files and directories on
			the the command line for Tweego to process.  Furthermore, it will
			search all directories encountered during processing looking for
			additional files and directories.  Thus, supporting StoryIncludes
			would be beyond pointless.

			If we see StoryIncludes, log a warning.
		*/
		log.Print(`warning: Ignoring "StoryIncludes" compiler special passage; and it is ` +
			`recommended that you remove it.  Tweego allows you to specify project ` +
			`files and/or directories to recursively search for such files on the ` +
			`command line.  Thus, in practice, you only need to specify a project's ` +
			`root directory and Tweego will find all of its files automatically.`)
	case "StoryData":
		if err := s.unmarshalStoryData([]byte(p.text)); err == nil {
			// Validiate the IFID.
			if len(s.ifid) > 0 {
				if err := validateIFID(s.ifid); err != nil {
					log.Fatalf(`error: Cannot validate IFID; %s.`, err.Error())
				}
			}

			// Rebuild the passage contents to remove deleted and/or erroneous entries.
			p.text = string(s.marshalStoryData())
		} else {
			// log.Printf(`warning: Cannot unmarshal "StoryData" compiler special passage; %s.`, err.Error())
			log.Fatalf(`error: Cannot unmarshal "StoryData" compiler special passage; %s.`, err.Error())
		}
	case "StorySettings":
		if err := s.unmarshalStorySettings([]byte(p.text)); err != nil {
			log.Printf(`warning: Cannot unmarshal "StorySettings" special passage; %s.`, err.Error())
		}
	case "StoryTitle":
		// Rebuild the passage contents to trim erroneous whitespace surrounding the title.
		p.text = strings.TrimSpace(p.text)
		s.name = p.text
	}

	s.append(p)
}
