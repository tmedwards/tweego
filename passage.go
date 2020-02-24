/*
	Copyright © 2014–2020 Thomas Michael Edwards. All rights reserved.
	Use of this source code is governed by a Simplified BSD License which
	can be found in the LICENSE file.
*/

package main

import (
	// standard packages
	"fmt"
	"regexp"
	"strings"
	// external packages
	"golang.org/x/text/unicode/norm"
)

// Info passages are passages which contain solely structural data, metadata,
// and code, rather than any actual story content.
var infoPassages = []string{
	// Story formats: Twine 1.4+ vanilla & SugarCube.
	"StoryAuthor", "StoryInit", "StoryMenu", "StorySubtitle", "StoryTitle",

	// Story formats: SugarCube.
	"PassageReady", "PassageDone", "PassageHeader", "PassageFooter", "StoryBanner", "StoryCaption",

	// Story formats: SugarCube (v1 only).
	"MenuOptions", "MenuShare", "MenuStory",

	// Story formats: SugarCube (v2 only).
	"StoryInterface", "StoryShare",

	// Story formats: Twine 1.4+ vanilla.
	// Compilers: Twine/Twee 1.4+, Twee2, & Tweego.
	"StorySettings",

	// Compilers: Tweego & (whatever Dan Cox's compiler is called).
	"StoryData",

	// Compilers: Twine/Twee 1.4+ & Twee2.
	"StoryIncludes",
}

type passageMetadata struct {
	position string // Unused by Tweego.  Twine 1 & 2 passage block X and Y coordinates CSV.
	size     string // Unused by Tweego.  Twine 2 passage block width and height CSV.
}

type passage struct {
	// Core.
	name string
	tags []string
	text string

	// Compiler metadata.
	metadata *passageMetadata
}

func newPassage(name string, tags []string, source string) *passage {
	return &passage{
		name: name,
		tags: tags,
		text: source,
	}
}

func (p *passage) equals(second passage) bool {
	return p.text == second.text
}

func (p *passage) tagsHas(needle string) bool {
	if len(p.tags) > 0 {
		for _, tag := range p.tags {
			if tag == needle {
				return true
			}
		}
	}
	return false
}

func (p *passage) tagsHasAny(needles ...string) bool {
	if len(p.tags) > 0 {
		for _, tag := range p.tags {
			for _, needle := range needles {
				if tag == needle {
					return true
				}
			}
		}
	}
	return false
}

func (p *passage) tagsContains(needle string) bool {
	if len(p.tags) > 0 {
		for _, tag := range p.tags {
			if strings.Contains(tag, needle) {
				return true
			}
		}
	}
	return false
}

func (p *passage) tagsStartsWith(needle string) bool {
	if len(p.tags) > 0 {
		for _, tag := range p.tags {
			if strings.HasPrefix(tag, needle) {
				return true
			}
		}
	}
	return false
}

func (p *passage) hasMetadataPosition() bool {
	return p.metadata != nil && p.metadata.position != ""
}

func (p *passage) hasMetadataSize() bool {
	return p.metadata != nil && p.metadata.size != ""
}

func (p *passage) hasAnyMetadata() bool {
	return p.metadata != nil && (p.metadata.position != "" || p.metadata.size != "")
}

func (p *passage) hasInfoTags() bool {
	return p.tagsHasAny("annotation", "script", "stylesheet", "widget") || p.tagsStartsWith("Twine.")
}

func (p *passage) hasInfoName() bool {
	return stringSliceContains(infoPassages, p.name)
}

func (p *passage) isInfoPassage() bool {
	return p.hasInfoName() || p.hasInfoTags()
}

func (p *passage) isStoryPassage() bool {
	return !p.hasInfoName() && !p.hasInfoTags()
}

func (p *passage) toTwee(outMode outputMode) string {
	var output string
	if outMode == outModeTwee3 {
		output = ":: " + tweeEscapeString(p.name)
		if len(p.tags) > 0 {
			output += " [" + tweeEscapeString(strings.Join(p.tags, " ")) + "]"
		}
		if p.hasAnyMetadata() {
			output += " " + string(p.marshalMetadata())
		}
	} else {
		output = ":: " + p.name
		if len(p.tags) > 0 {
			output += " [" + strings.Join(p.tags, " ") + "]"
		}
	}
	output += "\n"
	if len(p.text) > 0 {
		output += p.text + "\n"
	}
	output += "\n\n"
	return output
}

func (p *passage) toPassagedata(pid uint) string {
	var (
		position string
		size     string
	)
	if p.hasMetadataPosition() {
		position = p.metadata.position
	} else {
		// No position metadata, so generate something sensible on the fly.
		x := pid % 10
		y := pid / 10
		if x == 0 {
			x = 10
		} else {
			y++
		}
		position = fmt.Sprintf("%d,%d", x*125-25, y*125-25)
	}
	if p.hasMetadataSize() {
		size = p.metadata.size
	} else {
		// No size metadata, so default to the normal size.
		size = "100,100"
	}

	/*
		<tw-passagedata pid="…" name="…" tags="…" position="…" size="…">…</tw-passagedata>
	*/
	return fmt.Sprintf(`<tw-passagedata pid="%d" name=%q tags=%q position=%q size=%q>%s</tw-passagedata>`,
		pid,
		attrEscapeString(p.name),
		attrEscapeString(strings.Join(p.tags, " ")),
		attrEscapeString(position),
		attrEscapeString(size),
		htmlEscapeString(p.text),
	)
}

func (p *passage) toTiddler(pid uint) string {
	var position string
	if p.hasMetadataPosition() {
		position = p.metadata.position
	} else {
		// No position metadata, so generate something sensible on the fly.
		x := pid % 10
		y := pid / 10
		if x == 0 {
			x = 10
		} else {
			y++
		}
		position = fmt.Sprintf("%d,%d", x*140-130, y*140-130)
	}

	/*
		<div tiddler="…" tags="…" created="…" modifier="…" twine-position="…">…</div>
	*/
	return fmt.Sprintf(`<div tiddler=%q tags=%q twine-position=%q>%s</div>`,
		attrEscapeString(p.name),
		attrEscapeString(strings.Join(p.tags, " ")),
		attrEscapeString(position),
		tiddlerEscapeString(p.text),
	)
}

func (p *passage) countWords() uint64 {
	text := p.text

	// Strip newlines.
	text = strings.Replace(text, "\n", "", -1)

	// Strip comments.
	re := regexp.MustCompile(`(?s:/%.*?%/|/\*.*?\*/|<!--.*?-->)`)
	text = re.ReplaceAllString(text, "")

	// Count normalized "characters".
	var (
		count uint64
		ia    norm.Iter
	)
	ia.InitString(norm.NFKD, text)
	for !ia.Done() {
		count++
		ia.Next()
	}

	// Count "words", typing measurement style—i.e., 5 "characters" per "word".
	words := count / 5
	if count%5 > 0 {
		words++
	}

	return words
}
