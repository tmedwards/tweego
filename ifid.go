/*
	Copyright © 2014–2021 Thomas Michael Edwards. All rights reserved.
	Use of this source code is governed by a Simplified BSD License which
	can be found in the LICENSE file.
*/

package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"strings"
)

// An IFID (Interactive Fiction IDentifier) uniquely identifies compiled
// projects.  Most IFIDs are simply the string form of a v4 random UUID.
//
// IFIDs, in general, are defined within The Treaty of Babel.
//   SEE: http://babel.ifarchive.org/
//
// Twine ecosystem IFIDs are defined within both the Twee 3 Specification
// and Twine 2 HTML Output Specification.
//   SEE: https://github.com/iftechfoundation/twine-specs/

// newIFID generates a new IFID (UUID v4).
func newIFID() (string, error) {
	var uuid [16]byte
	if _, err := io.ReadFull(rand.Reader, uuid[:]); err != nil {
		return "", err
	}

	uuid[6] = (uuid[6] & 0x0f) | 0x40 // version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // variant 10

	return fmt.Sprintf("%X-%X-%X-%X-%X", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

// validateIFID validates ifid or returns an error.
func validateIFID(ifid string) error {
	switch len(ifid) {
	// xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	case 36:
		// no-op

	// UUID://xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx//
	case 36 + 9:
		if strings.ToUpper(ifid[:7]) != "UUID://" || ifid[43:] != "//" {
			return fmt.Errorf("invalid IFID UUID://…// format")
		}
		ifid = ifid[7:43]

	default:
		return fmt.Errorf("invalid IFID length: %d", len(ifid))
	}

	b := []byte(ifid)
	for i := 0; i < len(b); i++ {
		switch i {
		// hyphens
		case 8, 13, 18, 23:
			if b[i] != '-' {
				return fmt.Errorf("invalid IFID character %#U at position %d", b[i], i+1)
			}

		// version
		case 14:
			if '1' > b[i] || b[i] > '5' {
				return fmt.Errorf("invalid version %#U at position %d", b[i], i+1)
			}

		// variant
		case 19:
			switch b[i] {
			case '8', '9', 'a', 'A', 'b', 'B':
			default:
				return fmt.Errorf("invalid variant %#U at position %d", b[i], i+1)
			}

		// regular hex character
		default:
			switch {
			case '0' <= b[i] && b[i] <= '9':
			case 'a' <= b[i] && b[i] <= 'f':
			case 'A' <= b[i] && b[i] <= 'F':
			default:
				return fmt.Errorf("invalid IFID hex value %#U at position %d", b[i], i+1)
			}
		}
	}

	return nil
}
