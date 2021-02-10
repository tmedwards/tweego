/*
	Copyright © 2014–2021 Thomas Michael Edwards. All rights reserved.
	Use of this source code is governed by a Simplified BSD License which
	can be found in the LICENSE file.
*/

package main

import (
	"unicode"
)

// StringsInsensitively provides for case insensitively sorting slices of strings.
type StringsInsensitively []string

func (p StringsInsensitively) Len() int {
	return len(p)
}

func (p StringsInsensitively) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p StringsInsensitively) Less(i, j int) bool {
	iRunes := []rune(p[i])
	jRunes := []rune(p[j])

	uBound := len(iRunes)
	if uBound > len(jRunes) {
		uBound = len(jRunes)
	}

	for pos := 0; pos < uBound; pos++ {
		iR := iRunes[pos]
		jR := jRunes[pos]

		iRLo := unicode.ToLower(iR)
		jRLo := unicode.ToLower(jR)

		if iRLo != jRLo {
			return iRLo < jRLo
		}
		if iR != jR {
			return iR < jR
		}
	}

	return false
}
