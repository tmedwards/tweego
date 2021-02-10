/*
	Copyright © 2014–2021 Thomas Michael Edwards. All rights reserved.
	Use of this source code is governed by a Simplified BSD License which
	can be found in the LICENSE file.
*/

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func (s *story) toJSON(startName string) []byte {
	var passages = make([]*passageJSON, 0, 32)
	for _, p := range s.passages {
		if p.name == "StoryTitle" || p.name == "StoryData" {
			continue
		}

		passages = append(passages, p.toJSON())
	}

	marshaled, err := json.MarshalIndent(
		&storyJSON{
			s.name,
			s.ifid,
			startName,
			twine2OptionsMapToSlice(s.twine2.options),
			s.twine2.format,
			s.twine2.formatVersion,
			strings.Title(tweegoName),
			tweegoVersion.Version(),
			passages,
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

func (s *story) toTwee(outMode outputMode) []byte {
	var data []byte
	for _, p := range s.passages {
		data = append(data, p.toTwee(outMode)...)
	}
	return data
}

func (s *story) toTwine2Archive(startName string) []byte {
	return append(s.getTwine2DataChunk(startName), '\n')
}

func (s *story) toTwine1Archive(startName string) []byte {
	var (
		count    uint
		data     []byte
		template []byte
	)

	data, count = s.getTwine1PassageChunk()
	// NOTE: In Twine 1.4, the passage data wrapper is part of the story formats
	// themselves, so we have to create/forge one here.  We use the Twine 1.4 vanilla
	// `storeArea` ID, rather than SugarCube's preferred `store-area` ID, for maximum
	// compatibility and interoperability.
	template = append(template, fmt.Sprintf(`<div id="storeArea" data-size="%d" hidden>`, count)...)
	template = append(template, data...)
	template = append(template, "</div>\n"...)
	return template
}

func (s *story) toTwine2HTML(startName string) []byte {
	var template = s.format.source()

	// Story instance replacements.
	if bytes.Contains(template, []byte("{{STORY_NAME}}")) {
		template = bytes.Replace(template, []byte("{{STORY_NAME}}"), []byte(htmlEscapeString(s.name)), -1)
	}
	if bytes.Contains(template, []byte("{{STORY_DATA}}")) {
		template = bytes.Replace(template, []byte("{{STORY_DATA}}"), s.getTwine2DataChunk(startName), 1)
	}

	return template
}

func (s *story) toTwine1HTML(startName string) []byte {
	var (
		formatDir = filepath.Dir(s.format.filename)
		parentDir = filepath.Dir(formatDir)
		template  = s.format.source()
		count     uint
		data      []byte
		component []byte
		err       error
	)

	// Get the story data.
	data, count = s.getTwine1PassageChunk()

	// Story format compiler byline replacement.
	if search := []byte(`<a href="http://twinery.org/"`); bytes.Contains(template, search) {
		template = bytes.Replace(template, search, []byte(`<a href="http://www.motoslave.net/tweego/"`), 1)
		template = bytes.Replace(template, []byte(`>Twine</a>`), []byte(`>Tweego</a>`), 1)
	}

	// Story format component replacements (SugarCube).
	if search := []byte(`"USER_LIB"`); bytes.Contains(template, search) {
		component, err = fileReadAllAsUTF8(filepath.Join(formatDir, "userlib.js"))
		if err == nil {
			template = bytes.Replace(template, search, component, 1)
		} else if !os.IsNotExist(err) {
			log.Fatalf("error: %s", err.Error())
		}
	}

	// Story format component replacements (Twine 1.4+ vanilla story formats).
	if search := []byte(`"ENGINE"`); bytes.Contains(template, search) {
		component, err = fileReadAllAsUTF8(filepath.Join(parentDir, "engine.js"))
		if err != nil {
			log.Fatalf("error: %s", err.Error())
		}
		template = bytes.Replace(template, search, component, 1)
	}
	for _, pattern := range []string{`"SUGARCANE"`, `"JONAH"`} {
		if search := []byte(pattern); bytes.Contains(template, search) {
			component, err = fileReadAllAsUTF8(filepath.Join(formatDir, "code.js"))
			if err != nil {
				log.Fatalf("error: %s", err.Error())
			}
			template = bytes.Replace(template, search, component, 1)
		}
	}
	if s.twine1.settings["jquery"] == "on" {
		if search := []byte(`"JQUERY"`); bytes.Contains(template, search) {
			component, err = fileReadAllAsUTF8(filepath.Join(parentDir, "jquery.js"))
			if err != nil {
				log.Fatalf("error: %s", err.Error())
			}
			template = bytes.Replace(template, search, component, 1)
		}
	}
	if s.twine1.settings["modernizr"] == "on" {
		if search := []byte(`"MODERNIZR"`); bytes.Contains(template, search) {
			component, err = fileReadAllAsUTF8(filepath.Join(parentDir, "modernizr.js"))
			if err != nil {
				log.Fatalf("error: %s", err.Error())
			}
			template = bytes.Replace(template, search, component, 1)
		}
	}

	// Story instance replacements.
	if startName == defaultStartName {
		startName = ""
	}
	template = bytes.Replace(template, []byte(`"VERSION"`),
		[]byte(fmt.Sprintf("Compiled with %s, %s", tweegoName, tweegoVersion.Version())), 1)
	template = bytes.Replace(template, []byte(`"TIME"`),
		[]byte(fmt.Sprintf("Built on %s", time.Now().Format(time.RFC1123Z))), 1)
	template = bytes.Replace(template, []byte(`"START_AT"`),
		[]byte(fmt.Sprintf(`%q`, startName)), 1)
	template = bytes.Replace(template, []byte(`"STORY_SIZE"`),
		[]byte(fmt.Sprintf(`"%d"`, count)), 1)
	if bytes.Contains(template, []byte(`"STORY"`)) {
		// Twine/Twee ≥1.4 style story format.
		template = bytes.Replace(template, []byte(`"STORY"`), data, 1)
	} else {
		// Twine/Twee <1.4 style story format.
		var footer []byte
		footer, err = fileReadAllAsUTF8(filepath.Join(formatDir, "footer.html"))
		if err != nil {
			if os.IsNotExist(err) {
				footer = []byte("</div>\n</body>\n</html>\n")
			} else {
				log.Fatalf("error: %s", err.Error())
			}
		}
		template = append(template, data...)
		template = append(template, footer...)
	}

	// IFID replacement.
	if s.ifid != "" {
		if bytes.Contains(template, []byte(`<div id="store-area"`)) {
			// SugarCube
			template = bytes.Replace(template, []byte(`<div id="store-area"`),
				[]byte(fmt.Sprintf(`<!-- UUID://%s// --><div id="store-area"`, s.ifid)), 1)
		} else {
			// Twine/Twee vanilla story formats.
			template = bytes.Replace(template, []byte(`<div id="storeArea"`),
				[]byte(fmt.Sprintf(`<!-- UUID://%s// --><div id="storeArea"`, s.ifid)), 1)
		}
	}

	return template
}

func (s *story) getTwine2DataChunk(startName string) []byte {
	var (
		data    []byte
		startID string
		options string
		pid     uint
	)

	// Check the IFID status.
	if s.ifid == "" {
		var (
			ifid string
			err  error
		)
		if s.legacyIFID != "" {
			/*
				LEGACY
			*/
			log.Print(`error: Story IFID not found; reusing "ifid" entry from the "StorySettings" special passage.`)
			log.Println()
			ifid = s.legacyIFID
			/*
				END LEGACY
			*/
		} else {
			log.Print("error: Story IFID not found; generating one for your project.")
			log.Println()
			ifid, err = newIFID()
			if err != nil {
				log.Fatalf("error: IFID generation failed; %s", err.Error())
			}
		}
		ifid = fmt.Sprintf(`"ifid": %q`, ifid)
		base := "Copy the following "
		if s.has("StoryData") {
			ifid += ","
			log.Printf("%sline into the \"StoryData\" special passage's JSON block (at the top):\n\n\t%s\n\n", base, ifid)
			log.Printf("E.g., it should look something like the following:\n\n:: StoryData\n%s\n\n",
				bytes.Replace(s.marshalStoryData(), []byte("{"), []byte("{\n\t"+ifid), 1))
		} else {
			log.Printf("%s\"StoryData\" special passage into one of your project's twee source files:\n\n:: StoryData\n{\n\t%s\n}", base, ifid)
		}
		log.Fatalln()
	}

	// Gather all script and stylesheet passages.
	var (
		scripts     = make([]*passage, 0, 4)
		stylesheets = make([]*passage, 0, 4)
	)
	for _, p := range s.passages {
		if p.tagsHas("Twine.private") {
			continue
		}
		if p.tagsHas("script") {
			scripts = append(scripts, p)
		} else if p.tagsHas("stylesheet") {
			stylesheets = append(stylesheets, p)
		}
	}

	// Prepare the style element.
	/*
		<style role="stylesheet" id="twine-user-stylesheet" type="text/twine-css">…</style>
	*/
	data = append(data, `<style role="stylesheet" id="twine-user-stylesheet" type="text/twine-css">`...)
	if len(stylesheets) == 1 {
		data = append(data, stylesheets[0].text...)
	} else if len(stylesheets) > 1 {
		pid = 1
		for _, p := range stylesheets {
			if pid > 1 && data[len(data)-1] != '\n' {
				data = append(data, '\n')
			}
			data = append(data, fmt.Sprintf("/* twine-user-stylesheet #%d: %q */\n", pid, p.name)...)
			data = append(data, p.text...)
			pid++
		}
	}
	data = append(data, `</style>`...)

	// Prepare the script element.
	/*
		<script role="script" id="twine-user-script" type="text/twine-javascript">…</script>
	*/
	data = append(data, `<script role="script" id="twine-user-script" type="text/twine-javascript">`...)
	if len(scripts) == 1 {
		data = append(data, scripts[0].text...)
	} else if len(scripts) > 1 {
		pid = 1
		for _, p := range scripts {
			if pid > 1 && data[len(data)-1] != '\n' {
				data = append(data, '\n')
			}
			data = append(data, fmt.Sprintf("/* twine-user-script #%d: %q */\n", pid, p.name)...)
			data = append(data, p.text...)
			pid++
		}
	}
	data = append(data, `</script>`...)

	// Prepare tw-tag elements.
	/*
		<tw-tag name="…" color="…"></tw-tag>
	*/
	if s.twine2.tagColors != nil {
		for tag, color := range s.twine2.tagColors {
			data = append(data, fmt.Sprintf(`<tw-tag name=%q color=%q></tw-tag>`, tag, color)...)
		}
	}

	// Prepare normal passage elements.
	pid = 1
	for _, p := range s.passages {
		if p.name == "StoryTitle" || p.name == "StoryData" || p.tagsHasAny("script", "stylesheet", "Twine.private") {
			continue
		}

		/*
			LEGACY
		*/
		// TODO: Should we actually drop an empty StorySettings passage?
		if p.name == "StorySettings" && len(s.twine1.settings) == 0 {
			continue
		}
		/*
			END LEGACY
		*/

		data = append(data, p.toPassagedata(pid)...)
		if startName == p.name {
			startID = fmt.Sprint(pid)
		}
		pid++
	}

	// Add the <tw-storydata> wrapper.
	/*
		<tw-storydata name="…" startnode="…" creator="…" creator-version="…" ifid="…"
			zoom="…" format="…" format-version="…" options="…" hidden>…</tw-storydata>
	*/
	if optCount := len(s.twine2.options); optCount > 0 {
		opts := make([]string, 0, optCount)
		for opt, val := range s.twine2.options {
			if val {
				opts = append(opts, opt)
			}
		}
		options = strings.Join(opts, " ")
	}
	data = append([]byte(fmt.Sprintf(
		`<!-- UUID://%s// -->`+
			`<tw-storydata name=%q startnode=%q creator=%q creator-version=%q ifid=%q zoom=%q format=%q format-version=%q options=%q hidden>`,
		s.ifid,
		attrEscapeString(s.name),
		startID,
		attrEscapeString(strings.Title(tweegoName)),
		attrEscapeString(tweegoVersion.Version()),
		attrEscapeString(s.ifid),
		attrEscapeString(strconv.FormatFloat(s.twine2.zoom, 'f', -1, 32)),
		attrEscapeString(s.format.name),
		attrEscapeString(s.format.version),
		attrEscapeString(options),
	)), data...)
	data = append(data, `</tw-storydata>`...)

	return data
}

func (s *story) getTwine1PassageChunk() ([]byte, uint) {
	var (
		data  []byte
		count uint
	)

	for _, p := range s.passages {
		if p.tagsHas("Twine.private") {
			continue
		}

		count++
		data = append(data, p.toTiddler(count)...)
	}
	return data, count
}
