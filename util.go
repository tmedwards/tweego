/*
	Copyright © 2014–2021 Thomas Michael Edwards. All rights reserved.
	Use of this source code is governed by a Simplified BSD License which
	can be found in the LICENSE file.
*/

package main

import (
	"path/filepath"
	"regexp"
	"strings"
)

func mediaTypeFromFilename(filename string) string {
	return mediaTypeFromExt(normalizedFileExt(filename))
}

func mediaTypeFromExt(ext string) string {
	var mediaType string
	switch ext {
	// AUDIO NOTES:
	//
	// The preferred media type for WAVE audio is `audio/wave`, however,
	// some browsers only recognize `audio/wav`, requiring its use instead.
	case "aac", "flac", "ogg", "wav":
		mediaType = "audio/" + ext
	case "mp3":
		mediaType = "audio/mpeg"
	case "m4a":
		mediaType = "audio/mp4"
	case "oga", "opus":
		mediaType = "audio/ogg"
	case "wave":
		mediaType = "audio/wav"
	case "weba":
		mediaType = "audio/webm"

	// FONT NOTES:
	//
	// (ca. 2017) The IANA deprecated the various font subtypes of the
	// "application" type in favor of the new "font" type.  While the
	// standards were new at that point, many browsers had long accepted
	// such media types due to existing use in the wild—erroneous at
	// that point or not.
	//     otf   : application/font-sfnt  → font/otf
	//     ttf   : application/font-sfnt  → font/ttf
	//     woff  : application/font-woff  → font/woff
	//     woff2 : application/font-woff2 → font/woff2
	case "otf", "ttf", "woff", "woff2":
		mediaType = "font/" + ext

	// IMAGE NOTES:
	case "gif", "jpeg", "png", "tiff", "webp":
		mediaType = "image/" + ext
	case "jpg":
		mediaType = "image/jpeg"
	case "svg":
		mediaType = "image/svg+xml"
	case "tif":
		mediaType = "image/tiff"

	// METADATA NOTES:
	//
	// Name aside, WebVTT files are generic media cue metadata files
	// that may be used with either `<audio>` or `<video>` elements.
	case "vtt": // WebVTT (Web Video Text Tracks)
		mediaType = "text/vtt"

	// VIDEO NOTES:
	case "mp4", "webm":
		mediaType = "video/" + ext
	case "ogv":
		mediaType = "video/ogg"
	}

	return mediaType
}

func normalizedFileExt(filename string) string {
	ext := filepath.Ext(filename)
	if ext == "" {
		return ext
	}
	return strings.ToLower(ext[1:])
}

// Returns a trimmed and encoded slug of the passed string that should be safe
// for use as a DOM ID or class name.
func slugify(original string) string {
	// NOTE: The range of illegal characters consists of: C0 controls, space, exclamation,
	// double quote, number, dollar, percent, ampersand, single quote, left paren, right
	// paren, asterisk, plus, comma, hyphen, period, forward slash, colon, semi-colon,
	// less-than, equals, greater-than, question, at, left bracket, backslash, right
	// bracket, caret, backquote/grave, left brace, pipe/vertical-bar, right brace, tilde,
	// delete, C1 controls.
	illegalRe := regexp.MustCompile(`[\x00-\x20!-/:-@[-^\x60{-\x9f]+`)

	return illegalRe.ReplaceAllLiteralString(original, "-")
}

func stringSliceContains(haystack []string, needle string) bool {
	if len(haystack) > 0 {
		for _, val := range haystack {
			if val == needle {
				return true
			}
		}
	}
	return false
}

func stringSliceDelete(haystack []string, needle string) []string {
	if len(haystack) > 0 {
		for i, val := range haystack {
			if val == needle {
				copy(haystack[i:], haystack[i+1:])
				haystack[len(haystack)-1] = ""
				haystack = haystack[:len(haystack)-1]
				break
			}
		}
	}
	return haystack
}
