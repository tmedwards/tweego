/*
	Copyright © 2014–2019 Thomas Michael Edwards. All rights reserved.
	Use of this source code is governed by a Simplified BSD License which
	can be found in the LICENSE file.
*/

package twee2compat

import (
	"regexp"
)

// Twee2 line regexp: `^:: *([^\[]*?) *(\[(.*?)\])? *(<(.*?)>)? *$`
// See: https://github.com/Dan-Q/twee2/blob/d7659d84b5415d594dcc868628d74c3c9b48f496/lib/twee2/story_file.rb#L61

var (
	twee2DetectRe *regexp.Regexp
	twee2HeaderRe *regexp.Regexp
	twee2BadPosRe *regexp.Regexp
)

func hasTwee2Syntax(s []byte) bool {
	// Initialize and cache the regular expressions if necessary.
	if twee2DetectRe == nil {
		twee2DetectRe = regexp.MustCompile(`(?m)^:: *[^\[]*?(?: *\[.*?\])? *<(.*?)> *$`)
		twee2HeaderRe = regexp.MustCompile(`(?m)^(:: *[^\[]*?)( *\[.*?\])?(?: *<(.*?)>)? *$`)
		twee2BadPosRe = regexp.MustCompile(`(?m)^(::.*?) *{"position":" *"}$`)
	}
	return twee2DetectRe.Match(s)
}

// ToV3 returns a copy of the slice s with all instances of Twee2 position blocks
// replaced with Twee v3 metadata blocks.
func ToV3(s []byte) []byte {
	if hasTwee2Syntax(s) {
		s = twee2HeaderRe.ReplaceAll(s, []byte(`${1}${2} {"position":"${3}"}`))
		s = twee2BadPosRe.ReplaceAll(s, []byte(`$1`))
	}
	return s
}
