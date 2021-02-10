/*
	Copyright © 2014–2021 Thomas Michael Edwards. All rights reserved.
	Use of this source code is governed by a Simplified BSD License which
	can be found in the LICENSE file.
*/

package main

import (
	"bytes"
	"fmt"
	"log"
	"path/filepath"
	"strings"
)

func loadModules(filenames []string, encoding string) []byte {
	var (
		processedModules = make(map[string]bool)
		headTags         [][]byte
	)

	for _, filename := range filenames {
		if processedModules[filename] {
			log.Printf("warning: load %s: Skipping duplicate.", filename)
			continue
		}

		var (
			source []byte
			err    error
		)
		switch normalizedFileExt(filename) {
		// NOTE: The case values here should match those in `filesystem.go:knownFileType()`.
		case "css":
			source, err = loadModuleTagged("style", filename, encoding)
		case "js":
			source, err = loadModuleTagged("script", filename, encoding)
		case "otf", "ttf", "woff", "woff2":
			source, err = loadModuleFont(filename)
		default:
			// Simply ignore all other file types.
			continue
		}
		if err != nil {
			log.Fatalf("error: load %s: %s", filename, err.Error())
		}
		if len(source) > 0 {
			headTags = append(headTags, source)
		}
		processedModules[filename] = true
		statsAddExternalFile(filename)
	}

	return bytes.Join(headTags, []byte("\n"))
}

func loadModuleTagged(tag, filename, encoding string) ([]byte, error) {
	source, err := fileReadAllWithEncoding(filename, encoding)
	if err != nil {
		return nil, err
	}
	source = bytes.TrimSpace(source)
	if len(source) == 0 {
		return source, nil
	}

	var (
		idSlug   = tag + "-module-" + slugify(strings.Split(filepath.Base(filename), ".")[0])
		mimeType string
		b        bytes.Buffer
	)
	switch tag {
	case "script":
		mimeType = "text/javascript"
	case "style":
		mimeType = "text/css"
	}

	if _, err := fmt.Fprintf(
		&b,
		`<%s id=%q type=%q>%s</%[1]s>`,
		tag,
		idSlug,
		mimeType,
		source,
	); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func loadModuleFont(filename string) ([]byte, error) {
	source, err := fileReadAllAsBase64(filename)
	if err != nil {
		return nil, err
	}

	var (
		family    = strings.Split(filepath.Base(filename), ".")[0]
		idSlug    = "style-module-" + slugify(family)
		ext       = normalizedFileExt(filename)
		mediaType = mediaTypeFromExt(ext)
		hint      string
		b         bytes.Buffer
	)
	switch ext {
	case "ttf":
		hint = "truetype"
	case "otf":
		hint = "opentype"
	default:
		hint = ext
	}

	if _, err := fmt.Fprintf(
		&b,
		"<style id=%q type=\"text/css\">@font-face {\n\tfont-family: %q;\n\tsrc: url(\"data:%s;base64,%s\") format(%q);\n}</style>",
		idSlug,
		family,
		mediaType,
		source,
		hint,
	); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
