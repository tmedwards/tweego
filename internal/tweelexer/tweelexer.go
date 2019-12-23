/*
	Copyright © 2014–2019 Thomas Michael Edwards. All rights reserved.
	Use of this source code is governed by a Simplified BSD License which
	can be found in the LICENSE file.
*/

/*
	With kind regards to Rob Pike and his "Lexical Scanning in Go" talk.

	Any and all coding horrors within are my own.
*/

/*
	WARNING: Line Counts

	Ensuring proper line counts is fraught with peril as several methods
	modify the line count and it's entirely possible, if one is not careful,
	to count newlines multiple times.  For example, using `l.next()` to
	accept newlines, thus counting them, that are ultimately either emitted
	or ignored, which can cause them to be counted again.
*/
/*
	WARNING: Not Unicode Aware

	Twee syntax is strictly limited to US-ASCII, so there's no compelling
	reason to decode the UTF-8 input.
*/

package tweelexer

import (
	"bytes"
	"fmt"
)

// ItemType identifies the type of the items.
type ItemType int

// Item represents a lexed item, a lexeme.
type Item struct {
	Type ItemType // Type of the item.
	Line int      // Line within the input (1-base) of the item.
	Pos  int      // Starting position within the input, in bytes, of the item.
	Val  []byte   // Value of the item.
}

// String returns a formatted debugging string for the item.
func (i Item) String() string {
	var name string
	switch i.Type {
	case ItemEOF:
		return fmt.Sprintf("[EOF: %d/%d]", i.Line, i.Pos)
	case ItemError:
		name = "Error"
	case ItemHeader:
		name = "Header"
	case ItemName:
		name = "Name"
	case ItemTags:
		name = "Tags"
	case ItemMetadata:
		name = "Metadata"
	case ItemContent:
		name = "Content"
	}
	if i.Type != ItemError && len(i.Val) > 80 {
		return fmt.Sprintf("[%s: %d/%d] %.80q...", name, i.Line, i.Pos, i.Val)
	}
	return fmt.Sprintf("[%s: %d/%d] %q", name, i.Line, i.Pos, i.Val)
}

const eof = -1 // End of input value.

// TODO: golint claims ItemError, below, has no comment if the const
// block comment, below, is removed.  Report that lossage.

// Item type constants.
const (
	ItemError    ItemType = iota // Error.  Its value is the error message.
	ItemEOF                      // End of input.
	ItemHeader                   // '::', but only when starting a line.
	ItemName                     // Text w/ backslash escaped characters.
	ItemTags                     // '[tag1 tag2 tagN]'.
	ItemMetadata                 // JSON chunk, '{…}'.
	ItemContent                  // Plain text.
)

// stateFn state of the scanner as a function, which return the next state function.
type stateFn func(*Tweelexer) stateFn

// Tweelexer holds the state of the scanner.
type Tweelexer struct {
	input []byte    // Byte slice being scanned.
	line  int       // Number of newlines seen (1-base).
	start int       // Starting position of the current item.
	pos   int       // Current position within the input.
	items chan Item // Channel of scanned items.
}

// next returns the next byte, as a rune, in the input.
func (l *Tweelexer) next() rune {
	if l.pos >= len(l.input) {
		return eof
	}
	r := rune(l.input[l.pos])
	l.pos++
	if r == '\n' {
		l.line++
	}
	return r
}

// peek returns the next byte, as a rune, in the input, but does not consume it.
func (l *Tweelexer) peek() rune {
	if l.pos >= len(l.input) {
		return eof
	}
	return rune(l.input[l.pos])
}

// backup rewinds our position in the input by one byte.
func (l *Tweelexer) backup() {
	if l.pos > l.start {
		l.pos--
		if l.input[l.pos] == '\n' {
			l.line--
		}
	} else {
		panic(fmt.Errorf("backup would leave pos < start"))
	}
}

// emit sends an item to the item channel.
func (l *Tweelexer) emit(t ItemType) {
	l.items <- Item{t, l.line, l.start, l.input[l.start:l.pos]}
	// Some items may contain newlines that must be counted.
	if t == ItemContent {
		l.line += bytes.Count(l.input[l.start:l.pos], []byte("\n"))
	}
	l.start = l.pos
}

// ignore skips over the pending input.
func (l *Tweelexer) ignore() {
	l.line += bytes.Count(l.input[l.start:l.pos], []byte("\n"))
	l.start = l.pos
}

// accept consumes the next byte if it's from the valid set.
func (l *Tweelexer) accept(valid []byte) bool {
	if bytes.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of bytes from the valid set.
func (l *Tweelexer) acceptRun(valid []byte) {
	var r rune
	for r = l.next(); bytes.ContainsRune(valid, r); r = l.next() {
	}
	if r != eof {
		l.backup()
	}
}

// errorf emits an error item and returns nil, allowing the scan to be terminated
// simply by returning the call to errorf.
func (l *Tweelexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- Item{ItemError, l.line, l.start, []byte(fmt.Sprintf(format, args...))}
	return nil
}

// run runs the state machine for tweelexer.
func (l *Tweelexer) run() {
	for state := lexProlog; state != nil; {
		state = state(l)
	}
	close(l.items)
}

// NewTweelexer creates a new scanner for the input text.
func NewTweelexer(input []byte) *Tweelexer {
	l := &Tweelexer{
		input: input,
		line:  1,
		items: make(chan Item),
	}
	go l.run()
	return l
}

// GetItems returns the item channel.
// Called by the parser, not tweelexer.
func (l *Tweelexer) GetItems() chan Item {
	return l.items
}

// NextItem returns the next item and its ok status from the item channel.
// Called by the parser, not tweelexer.
func (l *Tweelexer) NextItem() (Item, bool) {
	// return <-l.items
	item, ok := <-l.items
	return item, ok
}

// Drain drains the item channel so the lexing goroutine will close the item channel and exit.
// Called by the parser, not tweelexer.
func (l *Tweelexer) Drain() {
	for range l.items {
	}
}

// acceptQuoted accepts a quoted string.
// The opening quote has already been seen.
func acceptQuoted(l *Tweelexer, quote rune) error {
Loop:
	for {
		switch l.next() {
		case '\\':
			if r := l.next(); r != '\n' && r != eof {
				break
			}
			fallthrough
		case '\n', eof:
			return fmt.Errorf("unterminated quoted string")
		case quote:
			break Loop
		}
	}
	return nil
}

// State functions.

var (
	headerDelim        = []byte("::")
	newlineHeaderDelim = []byte("\n::")
)

// lexProlog skips until the first passage header delimiter.
func lexProlog(l *Tweelexer) stateFn {
	if bytes.HasPrefix(l.input[l.pos:], headerDelim) {
		return lexHeaderDelim
	} else if i := bytes.Index(l.input[l.pos:], newlineHeaderDelim); i > -1 {
		l.pos += i + 1
		l.ignore()
		return lexHeaderDelim
	}
	l.emit(ItemEOF)
	return nil
}

// lexContent scans until a passage header delimiter.
func lexContent(l *Tweelexer) stateFn {
	if bytes.HasPrefix(l.input[l.pos:], headerDelim) {
		return lexHeaderDelim
	} else if i := bytes.Index(l.input[l.pos:], newlineHeaderDelim); i > -1 {
		l.pos += i + 1
		l.emit(ItemContent)
		return lexHeaderDelim
	}
	l.pos = len(l.input)
	if l.pos > l.start {
		l.emit(ItemContent)
	}
	l.emit(ItemEOF)
	return nil
}

// lexHeaderDelim scans a passage header delimiter.
func lexHeaderDelim(l *Tweelexer) stateFn {
	l.pos += len(headerDelim)
	l.emit(ItemHeader)
	return lexName
}

// lexName scans a passage name until: one of the optional block delimiters, newline, or EOF.
func lexName(l *Tweelexer) stateFn {
	var r rune
Loop:
	for {
		r = l.next()
		switch r {
		case '\\':
			r = l.next()
			if r != '\n' && r != eof {
				break
			}
			fallthrough
		case '[', ']', '{', '}', '\n', eof:
			if r != eof {
				l.backup()
			}
			break Loop
		}
	}
	// Always emit a name item, even if it's empty.
	l.emit(ItemName)

	switch r {
	case '[':
		return lexTags
	case ']':
		return l.errorf("unexpected right square bracket %#U", r)
	case '{':
		return lexMetadata
	case '}':
		return l.errorf("unexpected right curly brace %#U", r)
	case '\n':
		l.pos++
		l.ignore()
		return lexContent
	}
	l.emit(ItemEOF)
	return nil
}

// lexNextOptionalBlock scans within a header for the next optional block.
func lexNextOptionalBlock(l *Tweelexer) stateFn {
	// Consume space.
	l.acceptRun([]byte(" \t"))
	l.ignore()

	r := l.peek()
	// panic(fmt.Sprintf("[lexNextOptionalBlock: %d, %d:%d]", l.line, l.start, l.pos))
	switch r {
	case '[':
		return lexTags
	case ']':
		return l.errorf("unexpected right square bracket %#U", r)
	case '{':
		return lexMetadata
	case '}':
		return l.errorf("unexpected right curly brace %#U", r)
	case '\n':
		l.pos++
		l.ignore()
		return lexContent
	case eof:
		l.emit(ItemEOF)
		return nil
	}
	return l.errorf("illegal character %#U amid the optional blocks", r)
}

// lexTags scans an optional tags block.
func lexTags(l *Tweelexer) stateFn {
	// Consume the left delimiter '['.
	l.pos++

Loop:
	for {
		r := l.next()
		switch r {
		case '\\':
			r = l.next()
			if r != '\n' && r != eof {
				break
			}
			fallthrough
		case '\n', eof:
			if r == '\n' {
				l.backup()
			}
			return l.errorf("unterminated tag block")
		case ']':
			break Loop
		case '[':
			return l.errorf("unexpected left square bracket %#U", r)
		case '{':
			return l.errorf("unexpected left curly brace %#U", r)
		case '}':
			return l.errorf("unexpected right curly brace %#U", r)
		}
	}
	if l.pos > l.start {
		l.emit(ItemTags)
	}

	return lexNextOptionalBlock
}

// lexMetadata scans an optional (JSON) metadata block.
func lexMetadata(l *Tweelexer) stateFn {
	// Consume the left delimiter '{'.
	l.pos++

	depth := 1
Loop:
	for {
		r := l.next()
		// switch r {
		// case '"': // Only double quoted strings are legal within JSON chunks.
		// 	if err := acceptQuoted(l, '"'); err != nil {
		// 		return l.errorf(err.Error())
		// 	}
		// case '\\':
		// 	r = l.next()
		// 	if r != '\n' && r != eof {
		// 		break
		// 	}
		// 	fallthrough
		// case '\n', eof:
		// 	if r == '\n' {
		// 		l.backup()
		// 	}
		// 	return l.errorf("unterminated metadata block")
		// case '{':
		// 	depth++
		// case '}':
		// 	depth--
		// 	switch {
		// 	case depth == 0:
		// 		break Loop
		// 	case depth < 0:
		// 		return l.errorf("unbalanced curly braces in metadata block")
		// 	}
		// }
		switch r {
		case '"': // Only double quoted strings are legal within JSON chunks.
			if err := acceptQuoted(l, '"'); err != nil {
				return l.errorf(err.Error())
			}
		case '\n':
			l.backup()
			fallthrough
		case eof:
			return l.errorf("unterminated metadata block")
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				break Loop
			}
		}
	}
	if l.pos > l.start {
		l.emit(ItemMetadata)
	}

	return lexNextOptionalBlock
}
