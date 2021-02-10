/*
	Copyright © 2014–2021 Thomas Michael Edwards. All rights reserved.
	Use of this source code is governed by a Simplified BSD License which
	can be found in the LICENSE file.
*/

package main

import (
	// standard packages
	"fmt"
	"math"
	"os"
	"sort"

	// external packages
	"github.com/paulrosania/go-charset/charset"
)

// print basic help
func usage() {
	outFile := defaultOutFile
	if outFile == "-" {
		outFile = "<stdout>"
	}

	fmt.Fprintf(os.Stderr, `
Usage: %s [options] sources...

  sources                  Input sources (repeatable); may consist of supported
                             files and/or directories to recursively search for
                             such files.

Options:
  -a, --archive-twine2     Output Twine 2 archive, instead of compiled HTML.
      --archive-twine1     Output Twine 1 archive, instead of compiled HTML.
  -c SET, --charset=SET    Name of the input character set (default: "utf-8",
                             fallback: %q).
  -d, --decompile-twee3    Output Twee 3 source code, instead of compiled HTML.
      --decompile-twee1    Output Twee 1 source code, instead of compiled HTML.
  -f NAME, --format=NAME   ID of the story format (default: %q).
  -h, --help               Print this help, then exit.
      --head=FILE          Name of the file whose contents will be appended
                             as-is to the <head> element of the compiled HTML.
  -j, --json               Output JSON, instead of compiled HTML.
      --list-charsets      List the supported input character sets, then exit.
      --list-formats       List the available story formats, then exit.
      --log-files          Log the processed input files.
  -l, --log-stats          Log various story statistics.
  -m SRC, --module=SRC     Module sources (repeatable); may consist of supported
                             files and/or directories to recursively search for
                             such files.
      --no-trim            Do not trim whitespace surrounding passages.
  -o FILE, --output=FILE   Name of the output file (default: %q).
  -s NAME, --start=NAME    Name of the starting passage (default: the passage
                             set by the story data, elsewise %q).
  -t, --test               Compile in test mode; only for story formats in the
                             Twine 2 style.
      --twee2-compat       Enable Twee2 source compatibility mode; files with
                             the .tw2 or .twee2 extensions automatically have
                             compatibility mode enabled.
  -v, --version            Print version information, then exit.
  -w, --watch              Start watch mode; watch input sources for changes,
                             rebuilding the output as necessary.

`, tweegoName, fallbackCharset, defaultFormatID, outFile, defaultStartName)
	os.Exit(1)
}

// formats the list of supported character sets/encodings somewhat nicely for the user
func usageCharsets() {
	charsets := charset.Names()
	sort.Strings(charsets)

	cols := 4
	rows := int(math.Ceil(float64(len(charsets)) / float64(cols)))

	fmt.Fprintln(os.Stderr, "\nSupported input charsets:")
	for i, cnt := 0, len(charsets); i < rows; i++ {
		fmt.Fprintf(os.Stderr, "  %-18s", charsets[i])
		offset := rows
		for j := 0; j < cols; j++ {
			if i+offset < cnt {
				fmt.Fprintf(os.Stderr, " %-18s", charsets[i+offset])
				offset += rows
			}
		}
		fmt.Fprintln(os.Stderr)
	}
	fmt.Fprintln(os.Stderr)
	os.Exit(1)
}

// formats the list of supported story formats somewhat nicely for the user
func usageFormats(formats storyFormatsMap) {
	fmt.Fprintln(os.Stderr)
	if formats.isEmpty() {
		fmt.Fprintln(os.Stderr, "Story formats not found.")
	} else {
		ids := formats.ids()
		sort.Sort(StringsInsensitively(ids))
		fmt.Fprintln(os.Stderr, "Available formats:")
		fmt.Fprintln(os.Stderr, "  ID                     Name (Version) [Details]")
		fmt.Fprintln(os.Stderr, "  --------------------   ------------------------------")
		for _, id := range ids {
			f := formats[id]
			fmt.Fprintf(os.Stderr, "  %-20s", f.id)
			if f.isTwine2Style() {
				fmt.Fprintf(os.Stderr, "   %s (%s)", f.name, f.version)
				if f.proofing {
					fmt.Fprint(os.Stderr, " [proofing]")
				}
			}
			fmt.Fprintln(os.Stderr)
		}
	}
	fmt.Fprintln(os.Stderr)
	os.Exit(1)
}

func usageVersion() {
	fmt.Fprintf(os.Stderr, "\n%s, %s\n", tweegoName, tweegoVersion)
	fmt.Fprint(os.Stderr, `
Tweego (a Twee compiler in Go) [http://www.motoslave.net/tweego/]
Copyright (c) 2014-2021 Thomas Michael Edwards. All rights reserved.

`)
	os.Exit(1)
}
