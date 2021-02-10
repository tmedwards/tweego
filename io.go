/*
	Copyright © 2014–2021 Thomas Michael Edwards. All rights reserved.
	Use of this source code is governed by a Simplified BSD License which
	can be found in the LICENSE file.
*/

package main

import (
	// standard packages
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"unicode/utf8"

	// external packages
	"github.com/paulrosania/go-charset/charset"
	_ "github.com/paulrosania/go-charset/data" // import the charset data
)

const (
	// Assumed encoding of input files which are not UTF-8 encoded.
	fallbackCharset = "windows-1252" // match case from "charset" packages

	// Record separators.
	recordSeparatorLF   = "\n"   // I.e., UNIX-y OSes.
	recordSeparatorCRLF = "\r\n" // I.e., DOS/Windows.
	recordSeparatorCR   = "\r"   // I.e., MacOS ≤9.

	utfBOM = "\uFEFF"
)

func fileReadAllAsBase64(filename string) ([]byte, error) {
	var (
		r    io.Reader
		data []byte
		err  error
	)
	if filename == "-" {
		r = os.Stdin
	} else {
		var f *os.File
		if f, err = os.Open(filename); err != nil {
			return nil, err
		}
		defer f.Close()
		r = f
	}
	if data, err = ioutil.ReadAll(r); err != nil {
		return nil, err
	}
	buf := make([]byte, base64.StdEncoding.EncodedLen(len(data))) // try to avoid additional allocations
	base64.StdEncoding.Encode(buf, data)
	return buf, nil
}

func fileReadAllAsUTF8(filename string) ([]byte, error) {
	return fileReadAllWithEncoding(filename, "utf-8")
}

func fileReadAllWithEncoding(filename, encoding string) ([]byte, error) {
	var (
		r      io.Reader
		data   []byte
		rsLF   = []byte(recordSeparatorLF)
		rsCRLF = []byte(recordSeparatorCRLF)
		rsCR   = []byte(recordSeparatorCR)
		err    error
	)

	// Read in the entire file.
	if filename == "-" {
		r = os.Stdin
	} else {
		var f *os.File
		if f, err = os.Open(filename); err != nil {
			return nil, err
		}
		defer f.Close()
		r = f
	}
	if data, err = ioutil.ReadAll(r); err != nil {
		return nil, err
	}

	// Convert the charset to UTF-8, if necessary.
	encoding = charset.NormalizedName(encoding)
	if utf8.Valid(data) {
		switch encoding {
		case "", "utf-8", "utf8", "ascii", "us-ascii":
			// no-op
		default:
			log.Printf("warning: read %s: Already valid UTF-8; skipping charset conversion.", filename)
		}
	} else {
		switch encoding {
		case "utf-8", "utf8", "ascii", "us-ascii":
			log.Printf("warning: read %s: Invalid UTF-8; assuming charset is %s.", filename, fallbackCharset)
			fallthrough
		case "":
			encoding = charset.NormalizedName(fallbackCharset)
		}
		if r, err = charset.NewReader(encoding, bytes.NewReader(data)); err != nil {
			return nil, err
		}
		if data, err = ioutil.ReadAll(r); err != nil {
			return nil, err
		}
		if !utf8.Valid(data) {
			return nil, fmt.Errorf("read %s: Charset conversion yielded invalid UTF-8.", filename)
		}
	}

	// Strip the UTF BOM (\uFEFF), if it exists.
	if bytes.Equal(data[:3], []byte(utfBOM)) {
		data = data[3:]
	}

	// Normalize record separators.
	data = bytes.Replace(data, rsCRLF, rsLF, -1)
	data = bytes.Replace(data, rsCR, rsLF, -1)

	return data, nil
}

func alignRecordSeparators(data []byte) []byte {
	switch runtime.GOOS {
	case "windows":
		return bytes.Replace(data, []byte(recordSeparatorLF), []byte(recordSeparatorCRLF), -1)
	default:
		return data
	}
}

func modifyHead(data []byte, modulePaths []string, headFile, encoding string) []byte {
	var headTags [][]byte

	if len(modulePaths) > 0 {
		source := bytes.TrimSpace(loadModules(modulePaths, encoding))
		if len(source) > 0 {
			headTags = append(headTags, source)
		}
	}

	if headFile != "" {
		if source, err := fileReadAllWithEncoding(headFile, encoding); err == nil {
			source = bytes.TrimSpace(source)
			if len(source) > 0 {
				headTags = append(headTags, source)
			}
			statsAddExternalFile(headFile)
		} else {
			log.Fatalf("error: load %s: %s", headFile, err.Error())
		}
	}

	if len(headTags) > 0 {
		headTags = append(headTags, []byte("</head>"))
		return bytes.Replace(data, []byte("</head>"), bytes.Join(headTags, []byte("\n")), 1)
	}
	return data
}

func fileWriteAll(filename string, data []byte) (int, error) {
	var (
		w   io.Writer
		err error
	)
	if filename == "-" {
		w = os.Stdout
	} else {
		var f *os.File
		if f, err = os.Create(filename); err != nil {
			return 0, err
		}
		defer f.Close()
		w = f
	}
	return w.Write(data)
}
