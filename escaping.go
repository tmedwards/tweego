/*
	Copyright © 2014–2021 Thomas Michael Edwards. All rights reserved.
	Use of this source code is governed by a Simplified BSD License which
	can be found in the LICENSE file.
*/

package main

import (
	"strings"
)

/*
	HTML escaping/unescaping utilities.
*/

// Escape the minimum characters required for attribute values.
var attrEscaper = strings.NewReplacer(
	`&`, `&amp;`,
	`"`, `&quot;`,
	// QUESTION: Keep the following?  All markup we generate double quotes attribute
	// values, so escaping single quotes/apostrophes isn't actually necessary.
	`'`, `&#39;`,
)

func attrEscapeString(s string) string {
	if len(s) == 0 {
		return s
	}
	return attrEscaper.Replace(s)
}

// Escape the minimum characters required for general HTML escaping—i.e., only
// the special characters (`&`, `<`, `>`, `"`, `'`).
//
// NOTE: The following exists because `html.EscapeString()` converts double
// quotes (`"`) to their decimal numeric character reference (`&#34;`) rather
// than to their entity (`&quot;`).  While the behavior is entirely legal, and
// browsers will happily accept the NCRs, a not insignificant amount of code in
// the wild only checks for `&quot;` and will fail to properly unescape the NCR.
//
// The primary special characters (`&`, `<`, `>`, `"`) should always be
// converted to their entity forms and never to an NCR form.  Saving one byte
// (5 vs. 6) is not worth the issues it causes.
var htmlEscaper = strings.NewReplacer(
	`&`, `&amp;`,
	`<`, `&lt;`,
	`>`, `&gt;`,
	`"`, `&quot;`,
	`'`, `&#39;`,
)

func htmlEscapeString(s string) string {
	if len(s) == 0 {
		return s
	}
	return htmlEscaper.Replace(s)
}

var tiddlerEscaper = strings.NewReplacer(
	`&`, `&amp;`,
	`<`, `&lt;`,
	`>`, `&gt;`,
	`"`, `&quot;`,
	`\`, `\s`,
	"\t", `\t`,
	"\n", `\n`,
)

func tiddlerEscapeString(s string) string {
	if len(s) == 0 {
		return s
	}
	return tiddlerEscaper.Replace(s)
}

// NOTE: We only need the newline, tab, and backslash escapes here since
// `tiddlerUnescapeString()` is only used when loading Twine 1 HTML and the
// `x/net/html` package already handles entity/reference unescaping for us.
var tiddlerUnescaper = strings.NewReplacer(
	`\n`, "\n",
	`\t`, "\t",
	`\s`, `\`,
)

func tiddlerUnescapeString(s string) string {
	if len(s) == 0 {
		return s
	}
	return tiddlerUnescaper.Replace(s)
}

/*
	Twee escaping/unescaping utilities.
*/

// Encode set: '\\', '[', ']', '{', '}'.

func tweeEscapeBytes(s []byte) []byte {
	if len(s) == 0 {
		return []byte(nil)
	}

	// NOTE: The slices this will be used with will be short enough that
	// iterating a slice twice shouldn't be problematic.  That said,
	// assuming an escape count of 8 or so wouldn't be a terrible way to
	// handle this either.
	cnt := 0
	for _, b := range s {
		switch b {
		case '\\', '[', ']', '{', '}':
			cnt++
		}
	}
	e := make([]byte, 0, len(s)+cnt)
	for _, b := range s {
		switch b {
		case '\\', '[', ']', '{', '}':
			e = append(e, '\\')
		}
		e = append(e, b)
	}
	return e
}

var tweeEscaper = strings.NewReplacer(
	`\`, `\\`,
	`[`, `\[`,
	`]`, `\]`,
	`{`, `\{`,
	`}`, `\}`,
)

func tweeEscapeString(s string) string {
	if len(s) == 0 {
		return s
	}
	return tweeEscaper.Replace(s)
}

func tweeUnescapeBytes(s []byte) []byte {
	if len(s) == 0 {
		return []byte(nil)
	}
	u := make([]byte, 0, len(s))
	for i, l := 0, len(s); i < l; i++ {
		if s[i] == '\\' {
			i++
			if i >= l {
				break
			}
		}
		u = append(u, s[i])
	}
	return u
}

var tweeUnescaper = strings.NewReplacer(
	`\\`, `\`,
	`\[`, `[`,
	`\]`, `]`,
	`\{`, `{`,
	`\}`, `}`,
)

func tweeUnescapeString(s string) string {
	if len(s) == 0 {
		return s
	}
	return tweeUnescaper.Replace(s)
}
