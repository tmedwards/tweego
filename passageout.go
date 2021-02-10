/*
	Copyright © 2014–2021 Thomas Michael Edwards. All rights reserved.
	Use of this source code is governed by a Simplified BSD License which
	can be found in the LICENSE file.
*/

package main

import (
	// standard packages
	"fmt"
	"strings"
)

func (p *passage) toJSON() *passageJSON {
	return &passageJSON{
		Name: p.name,
		Tags: p.tags,
		Text: p.text,
	}
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
