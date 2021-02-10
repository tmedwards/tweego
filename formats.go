/*
	Copyright © 2014–2021 Thomas Michael Edwards. All rights reserved.
	Use of this source code is governed by a Simplified BSD License which
	can be found in the LICENSE file.
*/

package main

import (
	// standard packages
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	// external packages
	"github.com/Masterminds/semver/v3"
)

type twine2FormatJSON struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	// Description string `json:"description"`
	// Author      string `json:"author"`
	// Image       string `json:"image"`
	// URL string `json:"url"`
	// License     string `json:"license"`
	Proofing bool   `json:"proofing"`
	Source   string `json:"source"`
	// Setup       string `json:"-"`
}

type storyFormat struct {
	id       string
	filename string
	twine2   bool
	name     string
	version  string
	proofing bool
}

func (f *storyFormat) isTwine1Style() bool {
	return !f.twine2
}

func (f *storyFormat) isTwine2Style() bool {
	return f.twine2
}

func (f *storyFormat) getStoryFormatData(source []byte) (*twine2FormatJSON, error) {
	if !f.twine2 {
		return nil, errors.New("Not a Twine 2 style story format.")
	}

	// get the JSON chunk from the source
	first := bytes.Index(source, []byte("{"))
	last := bytes.LastIndex(source, []byte("}"))
	if first == -1 || last == -1 {
		return nil, errors.New("Could not find Twine 2 style story format JSON chunk.")
	}
	source = source[first : last+1]

	// parse the JSON
	data := &twine2FormatJSON{}
	if err := json.Unmarshal(source, data); err != nil {
		/*
			START Harlowe malformed JSON chunk workaround

			TODO: Remove this workaround that attempts to handle Harlowe's
			broken JSON chunk.

			NOTE: This worksaround is only possible because, currently,
			Harlowe's "setup" property is the last entry in the chunk.
		*/
		if strings.HasPrefix(strings.ToLower(f.id), "harlowe") {
			if i := bytes.LastIndex(source, []byte(`,"setup": function`)); i != -1 {
				// cut the "setup" property and its invalid value
				j := len(source) - 1
				source = append(source[:i], source[j:]...)
				return f.getStoryFormatData(source)
			}
		}
		/*
			If we've reached this point, either the format is not Harlowe
			or we cannot find the start of its "setup" property, so just
			return the JSON decoding error as normal.

			END Harlowe malformed JSON chunk workaround
		*/

		return nil, errors.New("Could not decode story format JSON chunk.")
	}

	return data, nil
}

func (f *storyFormat) unmarshalMetadata() error {
	if !f.twine2 {
		return nil
	}

	var (
		data   *twine2FormatJSON
		source []byte
		err    error
	)

	// read in the story format
	if source, err = fileReadAllAsUTF8(f.filename); err != nil {
		return err
	}

	// load various bits of metadata from the JSON
	if data, err = f.getStoryFormatData(source); err != nil {
		return err
	}
	f.name = data.Name
	f.version = data.Version
	f.proofing = data.Proofing

	return nil
}

func (f *storyFormat) source() []byte {
	var (
		source []byte
		err    error
	)

	// read in the story format
	if source, err = fileReadAllAsUTF8(f.filename); err != nil {
		log.Fatalf("error: format %s", err.Error())
	}

	// if in Twine 2 style, extract the actual source from the JSON
	if f.twine2 {
		var data *twine2FormatJSON
		if data, err = f.getStoryFormatData(source); err != nil {
			log.Fatalf("error: format %s: %s", f.id, err.Error())
		}
		source = []byte(data.Source)
	}

	return source
}

type storyFormatsMap map[string]*storyFormat

func newStoryFormatsMap(searchPaths []string) storyFormatsMap {
	var (
		baseFilenames = []string{"format.js", "header.html"}
		formats       = make(storyFormatsMap)
	)

	for _, searchDirname := range searchPaths {
		if info, err := os.Stat(searchDirname); err != nil || !info.IsDir() {
			continue
		}

		d, err := os.Open(searchDirname)
		if err != nil {
			continue
		}

		baseDirnames, err := d.Readdirnames(0)
		if err != nil {
			continue
		}

		for _, baseDirname := range baseDirnames {
			formatDirname := filepath.Join(searchDirname, baseDirname)
			if info, err := os.Stat(formatDirname); err != nil || !info.IsDir() {
				continue
			}

			for _, baseFilename := range baseFilenames {
				formatFilename := filepath.Join(formatDirname, baseFilename)
				if info, err := os.Stat(formatFilename); err == nil && info.Mode().IsRegular() {
					f := &storyFormat{
						id:       baseDirname,
						filename: formatFilename,
						twine2:   baseFilename == "format.js",
					}
					if err := f.unmarshalMetadata(); err != nil {
						log.Printf("warning: format %s: Skipping format; %s", f.id, err.Error())
						continue
					}
					formats[baseDirname] = f
					break
				}
			}
		}
	}

	return formats
}

func (m storyFormatsMap) isEmpty() bool {
	return len(m) == 0
}

func (m storyFormatsMap) getIDFromTwine2Name(name string) string {
	var (
		found *semver.Version
		id    string
	)

	for _, f := range m {
		if !f.twine2 {
			continue
		}

		if f.name == name {
			if have, err := semver.NewVersion(f.version); err == nil {
				if found == nil || have.GreaterThan(found) {
					found = have
					id = f.id
				}
			}
		}
	}

	return id
}

func (m storyFormatsMap) getIDFromTwine2NameAndVersion(name, version string) string {
	var (
		wanted *semver.Version
		found  *semver.Version
		id     string
	)
	if v, err := semver.NewVersion(version); err == nil {
		wanted = v
	} else {
		log.Printf("warning: format %q: Auto-selecting greatest version; Could not parse version %q.", name, version)
	}

	for _, f := range m {
		if !f.twine2 {
			continue
		}

		if f.name == name {
			if have, err := semver.NewVersion(f.version); err == nil {
				if wanted == nil || have.Major() == wanted.Major() && have.Compare(wanted) > -1 {
					if found == nil || have.GreaterThan(found) {
						found = have
						id = f.id
					}
				}
			}
		}
	}

	return id
}

func (m storyFormatsMap) hasByID(id string) bool {
	_, ok := m[id]
	return ok
}

func (m storyFormatsMap) hasByTwine2Name(name string) bool {
	_, ok := m[m.getIDFromTwine2Name(name)]
	return ok
}

func (m storyFormatsMap) hasByTwine2NameAndVersion(name, version string) bool {
	_, ok := m[m.getIDFromTwine2NameAndVersion(name, version)]
	return ok
}

func (m storyFormatsMap) getByID(id string) *storyFormat {
	return m[id]
}

func (m storyFormatsMap) getByTwine2Name(name string) *storyFormat {
	return m[m.getIDFromTwine2Name(name)]
}

func (m storyFormatsMap) getByTwine2NameAndVersion(name, version string) *storyFormat {
	return m[m.getIDFromTwine2NameAndVersion(name, version)]
}

func (m storyFormatsMap) ids() []string {
	var ids []string
	for id := range m {
		ids = append(ids, id)
	}
	return ids
}
