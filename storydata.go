/*
	Copyright © 2014–2021 Thomas Michael Edwards. All rights reserved.
	Use of this source code is governed by a Simplified BSD License which
	can be found in the LICENSE file.
*/

package main

import (
	"bytes"
	"encoding/json"
	"log"
	"strings"
)

type storyJSON struct {
	Name           string         `json:"name"`
	Ifid           string         `json:"ifid,omitempty"`
	Start          string         `json:"start,omitempty"`
	Options        []string       `json:"options,omitempty"`
	Format         string         `json:"format,omitempty"`
	FormatVersion  string         `json:"format-version,omitempty"`
	Creator        string         `json:"creator,omitempty"`
	CreatorVersion string         `json:"creator-version,omitempty"`
	Passages       []*passageJSON `json:"passages"`
}

type storyDataJSON struct {
	Ifid          string             `json:"ifid,omitempty"`
	Format        string             `json:"format,omitempty"`
	FormatVersion string             `json:"format-version,omitempty"`
	Options       []string           `json:"options,omitempty"`
	Start         string             `json:"start,omitempty"`
	TagColors     twine2TagColorsMap `json:"tag-colors,omitempty"`
	Zoom          float64            `json:"zoom,omitempty"`
}

func (s *story) marshalStoryData() []byte {
	marshaled, err := json.MarshalIndent(
		&storyDataJSON{
			s.ifid,
			s.twine2.format,
			s.twine2.formatVersion,
			twine2OptionsMapToSlice(s.twine2.options),
			s.twine2.start,
			s.twine2.tagColors,
			s.twine2.zoom,
		},
		"",
		"\t",
	)
	if err != nil {
		// NOTE: We should never be able to see an error here.  If we do,
		// then something truly exceptional—in a bad way—has happened, so
		// we get our panic on.
		panic(err)
	}
	return marshaled
}

func (s *story) unmarshalStoryData(marshaled []byte) error {
	storyData := storyDataJSON{}
	if err := json.Unmarshal(marshaled, &storyData); err != nil {
		return err
	}
	s.ifid = strings.ToUpper(storyData.Ifid) // NOTE: Force uppercase for consistency.
	s.twine2.format = storyData.Format
	s.twine2.formatVersion = storyData.FormatVersion
	s.twine2.options = twine2OptionsSliceToMap(storyData.Options)
	s.twine2.start = storyData.Start
	s.twine2.tagColors = storyData.TagColors
	if storyData.Zoom != 0 {
		s.twine2.zoom = storyData.Zoom
	}
	return nil
}

func twine2OptionsMapToSlice(optMap twine2OptionsMap) []string {
	optSlice := []string{}
	if len(optMap) > 0 {
		for opt, val := range optMap {
			if val {
				optSlice = append(optSlice, opt)
			}
		}
	}
	return optSlice
}

func twine2OptionsSliceToMap(optSlice []string) twine2OptionsMap {
	optMap := make(twine2OptionsMap)
	for _, opt := range optSlice {
		optMap[opt] = true
	}
	return optMap
}

// func (s *story) marshalStorySettings() []byte {
// 	var marshaled [][]byte
// 	for key, val := range s.twine1.settings {
// 		marshaled = append(marshaled, []byte(key+":"+val))
// 	}
// 	return bytes.Join(marshaled, []byte("\n"))
// }

func (s *story) unmarshalStorySettings(marshaled []byte) error {
	/*
		NOTE: (ca. Feb 28, 2019) Transition away from storing metadata within
		the StorySettings special passage and to the StoryData special passages
		for two reasons:

		1. I've discovered that it's not as Twine 1-safe as I'd originally believed.
		   When Twine 1 imports a StorySettings passage, it does not check if fields
		   exist before appending "missing" fields, so it's entirely possible to end
		   up with the first appended field essentially being concatenated to the end
		   of the last of the previously existing fields.  Not good.
		2. Twee 3 standardization
	*/
	/*
		LEGACY
	*/
	var obsolete []string
	/*
		END LEGACY
	*/
	for _, line := range bytes.Split(marshaled, []byte{'\n'}) {
		line = bytes.TrimSpace(line)
		if len(line) > 0 {
			if i := bytes.IndexRune(line, ':'); i != -1 {
				key := string(bytes.ToLower(bytes.TrimSpace(line[:i])))
				val := string(bytes.ToLower(bytes.TrimSpace(line[i+1:])))

				/*
					LEGACY
				*/
				switch key {
				case "ifid":
					if err := validateIFID(val); err == nil {
						s.legacyIFID = strings.ToUpper(val) // NOTE: Force uppercase for consistency.
					}
					obsolete = append(obsolete, `"ifid"`)
					continue
				case "zoom":
					// NOTE: Just drop it.
					obsolete = append(obsolete, `"zoom"`)
					continue
				}
				/*
					END LEGACY
				*/

				s.twine1.settings[key] = val
			} else {
				log.Printf(`warning: Malformed "StorySettings" entry; skipping %q.`, line)
			}
		}
	}
	/*
		LEGACY
	*/
	if len(obsolete) > 0 {
		var (
			entries string
			pronoun string
		)
		if len(obsolete) == 1 {
			entries = "entry"
			pronoun = "it"
		} else {
			entries = "entries"
			pronoun = "them"
		}
		log.Printf(
			`warning: Detected obsolete "StorySettings" %s: %s.  `+
				`Please remove %s from the "StorySettings" special passage.  If doing `+
				`so leaves the passage empty, please remove it as well.`,
			entries,
			strings.Join(obsolete, ", "),
			pronoun,
		)
	}
	/*
		END LEGACY
	*/

	return nil
}
