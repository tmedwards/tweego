/*
	Copyright © 2014–2021 Thomas Michael Edwards. All rights reserved.
	Use of this source code is governed by a Simplified BSD License which
	can be found in the LICENSE file.
*/

package main

import (
	// standard packages
	"bytes"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	// internal packages
	twee2 "github.com/tmedwards/tweego/internal/twee2compat"
	twlex "github.com/tmedwards/tweego/internal/tweelexer"

	// external packages
	"golang.org/x/net/html"
)

func (s *story) load(filenames []string, c *config) {
	for _, filename := range filenames {
		if s.processed[filename] {
			log.Printf("warning: load %s: Skipping duplicate.", filename)
			continue
		}

		switch normalizedFileExt(filename) {
		// NOTE: The case values here should match those in `filesystem.go:knownFileType()`.
		case "tw", "twee":
			if err := s.loadTwee(filename, c.encoding, c.trim, c.twee2Compat); err != nil {
				log.Fatalf("error: load %s: %s", filename, err.Error())
			}
		case "tw2", "twee2":
			if err := s.loadTwee(filename, c.encoding, c.trim, true); err != nil {
				log.Fatalf("error: load %s: %s", filename, err.Error())
			}
		case "htm", "html":
			if err := s.loadHTML(filename, c.encoding); err != nil {
				log.Fatalf("error: load %s: %s", filename, err.Error())
			}
		case "css":
			if err := s.loadTagged("stylesheet", filename, c.encoding); err != nil {
				log.Fatalf("error: load %s: %s", filename, err.Error())
			}
		case "js":
			if err := s.loadTagged("script", filename, c.encoding); err != nil {
				log.Fatalf("error: load %s: %s", filename, err.Error())
			}
		case "otf", "ttf", "woff", "woff2":
			if err := s.loadFont(filename); err != nil {
				log.Fatalf("error: load %s: %s", filename, err.Error())
			}
		case "gif", "jpeg", "jpg", "png", "svg", "tif", "tiff", "webp":
			if err := s.loadMedia("Twine.image", filename); err != nil {
				log.Fatalf("error: load %s: %s", filename, err.Error())
			}
		case "aac", "flac", "m4a", "mp3", "oga", "ogg", "opus", "wav", "wave", "weba":
			if err := s.loadMedia("Twine.audio", filename); err != nil {
				log.Fatalf("error: load %s: %s", filename, err.Error())
			}
		case "mp4", "ogv", "webm":
			if err := s.loadMedia("Twine.video", filename); err != nil {
				log.Fatalf("error: load %s: %s", filename, err.Error())
			}
		case "vtt":
			if err := s.loadMedia("Twine.vtt", filename); err != nil {
				log.Fatalf("error: load %s: %s", filename, err.Error())
			}
		default:
			// Simply ignore all other file types.
			continue
		}
		s.processed[filename] = true
		statsAddProjectFile(filename)
	}

	/*
		Postprocessing.
	*/

	// Prepend the `StoryTitle` special passage, if necessary.
	if s.name != "" && !s.has("StoryTitle") {
		s.prepend(newPassage("StoryTitle", []string{}, s.name))
	}
}

func (s *story) loadTwee(filename, encoding string, trim, twee2Compat bool) error {
	source, err := fileReadAllWithEncoding(filename, encoding)
	if err != nil {
		return err
	}

	if twee2Compat {
		source = twee2.ToV3(source)
	}

	var (
		pCount   = 0
		lastType twlex.ItemType
		lex      = twlex.NewTweelexer(source)
	)

ParseLoop:
	for {
		p := &passage{}
		for item, ok := lex.NextItem(); ok; item, ok = lex.NextItem() {
			switch item.Type {
			case twlex.ItemError:
				return fmt.Errorf("line %d: Malformed twee source; %s.", item.Line, item.Val)

			case twlex.ItemEOF:
				// Add the final passage, if any.
				if pCount > 0 {
					s.add(p)
				}
				break ParseLoop

			case twlex.ItemHeader:
				pCount++
				if pCount > 1 {
					s.add(p)
					p = &passage{}
				}

			case twlex.ItemName:
				p.name = string(bytes.TrimSpace(tweeUnescapeBytes(item.Val)))
				if len(p.name) == 0 {
					lex.Drain()
					return fmt.Errorf("line %d: Malformed twee source; passage with no name.", item.Line)
				}

			case twlex.ItemTags:
				if lastType != twlex.ItemName {
					lex.Drain()
					return fmt.Errorf("line %d: Malformed twee source; optional tags block must immediately follow the passage name.", item.Line)
				}
				p.tags = strings.Fields(string(tweeUnescapeBytes(item.Val[1 : len(item.Val)-1])))

			case twlex.ItemMetadata:
				if lastType != twlex.ItemName && lastType != twlex.ItemTags {
					lex.Drain()
					return fmt.Errorf("line %d: Malformed twee source; optional metadata block must immediately follow the passage name or tags block.", item.Line)
				}
				if err := p.unmarshalMetadata(item.Val); err != nil {
					log.Printf("warning: load %s: line %d: Malformed twee source; could not decode metadata (reason: %s).", filename, item.Line, err.Error())
				}

			case twlex.ItemContent:
				if trim {
					// Trim whitespace surrounding (leading and trailing) passages.
					p.text = string(bytes.TrimSpace(item.Val))
				} else {
					// Do not trim whitespace surrounding passages.
					p.text = string(item.Val)
				}
			}

			lastType = item.Type
		}
	}

	return nil
}

func (s *story) loadHTML(filename, encoding string) error {
	source, err := fileReadAllWithEncoding(filename, encoding)
	if err != nil {
		return err
	}

	doc, err := getDocumentTree(bytes.TrimSpace(source))
	if err != nil {
		return fmt.Errorf("Malformed HTML source; %s.", err.Error())
	}

	if storyData := getElementByTag(doc, "tw-storydata"); storyData != nil {
		// Twine 2 style story data chunk.
		/*
			<tw-storydata name="…" startnode="…" creator="…" creator-version="…" ifid="…"
				zoom="…" format="…" format-version="…" options="…" hidden>…</tw-storydata>
		*/

		var startnode int

		// Content attribute processing.
		for _, a := range storyData.Attr {
			switch a.Key {
			case "name":
				s.name = a.Val
			case "startnode":
				if iVal, err := strconv.Atoi(a.Val); err == nil {
					startnode = iVal
				} else {
					log.Printf(`warning: Cannot parse "tw-storydata" content attribute "startnode" as an integer; value %q.`, a.Val)
				}
			// case "creator": Discard.
			// case "creator-version": Discard.
			case "ifid":
				s.ifid = strings.ToUpper(a.Val) // Force uppercase for consistency.
			case "zoom":
				if fVal, err := strconv.ParseFloat(a.Val, 64); err == nil {
					s.twine2.zoom = fVal
				} else {
					log.Printf(`warning: Cannot parse "tw-storydata" content attribute "zoom" as a float; value %q.`, a.Val)
				}
			case "format":
				s.twine2.format = a.Val
			case "format-version":
				s.twine2.formatVersion = a.Val
			case "options":
				// FIXME: I'm unsure whether the `options` content attribute is
				// intended to be a space delimited list.  That does seem likely,
				// so we treat it as such for now.
				for _, opt := range strings.Fields(a.Val) {
					s.twine2.options[opt] = true
				}
			}
		}

		// Node processing.
		for node := storyData.FirstChild; node != nil; node = node.NextSibling {
			if node.Type != html.ElementNode {
				continue
			}

			var (
				pid      int
				name     string
				tags     []string
				content  string
				metadata *passageMetadata
			)

			switch node.Data {
			case "style", "script":
				/*
					<style role="stylesheet" id="twine-user-stylesheet" type="text/twine-css">…</style>

					<script role="script" id="twine-user-script" type="text/twine-javascript">…</script>
				*/
				if node.FirstChild == nil {
					// skip empty elements
					continue
				} else {
					nodeData := strings.TrimSpace(node.FirstChild.Data)
					if len(nodeData) == 0 {
						// NOTE: Skip elements that are empty after trimming; this additional
						// "empty" check is necessary because most (all?) versions of Twine 2
						// habitually append newlines to the nodes, so they're almost never
						// actually empty.
						continue
					}
					if node.Data == "style" {
						name = "Story Stylesheet"
						tags = []string{"stylesheet"}
					} else {
						name = "Story JavaScript"
						tags = []string{"script"}
					}
					content = nodeData
				}
			case "tw-tag":
				/*
					<tw-tag name="…" color="…"></tw-tag>
				*/
				{
					var (
						tagName  string
						tagColor string
					)
					for _, a := range node.Attr {
						switch a.Key {
						case "name":
							tagName = a.Val
						case "color":
							tagColor = a.Val
						}
					}
					s.twine2.tagColors[tagName] = tagColor
				}
				continue
			case "tw-passagedata":
				/*
					<tw-passagedata pid="…" name="…" tags="…" position="…" size="…">…</tw-passagedata>
				*/
				metadata = &passageMetadata{}
				for _, a := range node.Attr {
					switch a.Key {
					case "pid":
						if iVal, err := strconv.Atoi(a.Val); err == nil {
							pid = iVal
						} else {
							log.Printf(`warning: Cannot parse "tw-passagedata" content attribute "pid" as an integer; value %q.`, a.Val)
						}
					case "name":
						name = a.Val
					case "tags":
						tags = strings.Fields(a.Val)
					case "position":
						metadata.position = a.Val
					case "size":
						metadata.size = a.Val
					}
				}
				if pid == startnode {
					s.twine2.start = name
				}
				if node.FirstChild != nil {
					content = node.FirstChild.Data
				}
			default:
				continue
			}

			p := newPassage(name, tags, strings.TrimSpace(content))
			if metadata != nil {
				p.metadata = metadata
			}
			s.add(p)
		}

		// Prepend the `StoryData` special passage.  Includes the story IFID and Twine 2 metadata.
		s.prepend(newPassage("StoryData", []string{}, string(s.marshalStoryData())))
	} else if storyData := getElementByID(doc, "store(?:-a|A)rea"); storyData != nil {
		// Twine 1 style story data chunk.
		/*
			<div id="store-area" data-size="…" hidden>…</div>
		*/
		for node := storyData.FirstChild; node != nil; node = node.NextSibling {
			if node.Type != html.ElementNode || node.Data != "div" || !hasAttr(node, "tiddler") {
				continue
			}

			var (
				name     string
				tags     []string
				content  string
				metadata = &passageMetadata{}
			)

			/*
				<div tiddler="…" tags="…" created="…" modified="…" modifier="…" twine-position="…">…</div>
			*/
			for _, a := range node.Attr {
				// NOTE: Ignore the following content attributes: `created`, `modified`, `modifier`.
				switch a.Key {
				case "tiddler":
					name = a.Val
				case "tags":
					tags = strings.Fields(a.Val)
				case "twine-position":
					metadata.position = a.Val
				}
			}
			if node.FirstChild != nil {
				content = tiddlerUnescapeString(node.FirstChild.Data)
			}

			p := newPassage(name, tags, strings.TrimSpace(content))
			if metadata != nil {
				p.metadata = metadata
			}
			s.add(p)
		}
	} else {
		return fmt.Errorf("Malformed HTML source; story data not found.")
	}

	return nil
}

func (s *story) loadTagged(tag, filename, encoding string) error {
	source, err := fileReadAllWithEncoding(filename, encoding)
	if err != nil {
		return err
	}

	s.add(newPassage(
		filepath.Base(filename),
		[]string{tag},
		string(source),
	))

	return nil
}

func (s *story) loadMedia(tag, filename string) error {
	source, err := fileReadAllAsBase64(filename)
	if err != nil {
		return err
	}

	s.add(newPassage(
		strings.Split(filepath.Base(filename), ".")[0],
		[]string{tag},
		"data:"+mediaTypeFromFilename(filename)+";base64,"+string(source),
	))

	return nil
}

func (s *story) loadFont(filename string) error {
	source, err := fileReadAllAsBase64(filename)
	if err != nil {
		return err
	}

	var (
		name      = filepath.Base(filename)
		family    = strings.Split(name, ".")[0]
		ext       = normalizedFileExt(filename)
		mediaType = mediaTypeFromExt(ext)
		hint      string
	)
	switch ext {
	case "ttf":
		hint = "truetype"
	case "otf":
		hint = "opentype"
	default:
		hint = ext
	}

	s.add(newPassage(
		name,
		[]string{"stylesheet"},
		fmt.Sprintf(
			"@font-face {\n\tfont-family: %q;\n\tsrc: url(\"data:%s;base64,%s\") format(%q);\n}",
			family,
			mediaType,
			source,
			hint,
		),
	))

	return nil
}
