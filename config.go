/*
	Copyright © 2014–2021 Thomas Michael Edwards. All rights reserved.
	Use of this source code is governed by a Simplified BSD License which
	can be found in the LICENSE file.
*/

package main

import (
	// standard packages
	"log"
	"os"
	"path/filepath"

	// internal packages
	"github.com/tmedwards/tweego/internal/option"
	// external packages
	"github.com/paulrosania/go-charset/charset"
)

type outputMode int

const (
	outModeHTML outputMode = iota
	outModeJSON
	outModeTwee3
	outModeTwee1
	outModeTwine2Archive
	outModeTwine1Archive
)

type common struct {
	formatID  string // ID of the story format to use
	startName string // name of the starting passage
}

type config struct {
	cmdline common
	common

	encoding    string     // input encoding
	sourcePaths []string   // slice of paths to seach for source files
	modulePaths []string   // slice of paths to seach for module files
	headFile    string     // name of the head file
	outFile     string     // name of the output file
	outMode     outputMode // output mode

	formats     storyFormatsMap // map of all enumerated story formats
	logFiles    bool            // log input files
	logStats    bool            // log story statistics
	testMode    bool            // enable test mode
	trim        bool            // enable passage trimming
	twee2Compat bool            // enable Twee2 header extension compatibility mode
	watchFiles  bool            // enable filesystem watching
}

const (
	defaultFormatID  = "sugarcube-2"
	defaultOutFile   = "-" // <stdout>
	defaultOutMode   = outModeHTML
	defaultStartName = "Start"
	defaultTrimState = true
)

// newConfig creates a new config instance
func newConfig() *config {
	// Get the base paths to search for story formats.
	formatDirs := (func() []string {
		var (
			baseDirnames = []string{
				"storyformats",
				".storyformats",
				"story-formats", // DEPRECATED
				"storyFormats",  // DEPRECATED
				"targets",       // DEPRECATED
			}
			basePaths      = []string{programDir}
			searchDirnames []string
		)
		if homeDir, err := userHomeDir(); err == nil {
			if !stringSliceContains(basePaths, homeDir) {
				basePaths = append(basePaths, homeDir)
			}
		}
		if !stringSliceContains(basePaths, workingDir) {
			basePaths = append(basePaths, workingDir)
		}
		for _, basePath := range basePaths {
			for _, baseDirname := range baseDirnames {
				searchDirname := filepath.Join(basePath, baseDirname)
				if info, err := os.Stat(searchDirname); err == nil && info.IsDir() {
					searchDirnames = append(searchDirnames, searchDirname)
				}
			}
		}
		return searchDirnames
	})()

	// Create a new instance of `config` and assign defaults.
	c := &config{
		common:  common{formatID: defaultFormatID, startName: defaultStartName},
		outFile: defaultOutFile,
		outMode: defaultOutMode,
		trim:    defaultTrimState,
	}

	// Merge values from the environment variables.
	if env := os.Getenv("TWEEGO_PATH"); env != "" {
		formatDirs = append(formatDirs, filepath.SplitList(env)...)
	}

	// TODO: Move story formats out of the config?
	// Enumerate story formats.
	if len(formatDirs) == 0 {
		log.Fatal("error: Story format search directories not found.")
	}
	c.formats = newStoryFormatsMap(formatDirs)
	if c.formats.isEmpty() {
		log.Print("error: Story formats not found within the search directories: (in order)")
		for i, path := range formatDirs {
			log.Printf("  %2d. %s", i+1, path)
		}
		os.Exit(1)
	}

	// Merge values from the command line.
	options := option.NewParser()
	options.Add("archive_twine2", "-a|--archive-twine2")
	options.Add("archive_twine1", "--archive-twine1")
	options.Add("decompile_twee3", "-d|--decompile-twee3|--decompile") // NOTE: "--decompile" is deprecated.
	options.Add("decompile_twee1", "--decompile-twee1")
	options.Add("encoding", "-c=s|--charset=s")
	options.Add("format", "-f=s|--format=s")
	options.Add("head", "--head=s")
	options.Add("help", "-h|--help")
	options.Add("json", "-j|--json")
	options.Add("listcharsets", "--list-charsets")
	options.Add("listformats", "--list-formats")
	options.Add("logfiles", "--log-files")
	options.Add("logstats", "-l|--log-stats")
	options.Add("module", "-m=s+|--module=s+")
	options.Add("no_trim", "--no-trim")
	options.Add("output", "-o=s|--output=s")
	options.Add("start", "-s=s|--start=s")
	options.Add("test", "-t|--test")
	options.Add("twee2_compat", "--twee2-compat")
	options.Add("version", "-v|--version")
	options.Add("watch", "-w|--watch")
	if opts, sources, err := options.ParseCommandLine(); err == nil {
		for opt, val := range opts {
			switch opt {
			case "archive_twine2":
				c.outMode = outModeTwine2Archive
			case "archive_twine1":
				c.outMode = outModeTwine1Archive
			case "decompile_twee3":
				c.outMode = outModeTwee3
			case "decompile_twee1":
				c.outMode = outModeTwee1
			case "encoding":
				c.encoding = val.(string)
			case "format":
				c.cmdline.formatID = val.(string)
				c.formatID = c.cmdline.formatID
			case "head":
				c.headFile = val.(string)
			case "help":
				usage()
			case "json":
				c.outMode = outModeJSON
			case "listcharsets":
				usageCharsets()
			case "listformats":
				usageFormats(c.formats)
			case "logfiles":
				c.logFiles = true
			case "logstats":
				c.logStats = true
			case "module":
				c.modulePaths = append(c.modulePaths, val.([]string)...)
			case "no_trim":
				c.trim = false
			case "output":
				c.outFile = val.(string)
			case "start":
				c.cmdline.startName = val.(string)
				c.startName = c.cmdline.startName
			case "test":
				c.testMode = true
			case "twee2_compat":
				c.twee2Compat = true
			case "version":
				usageVersion()
			case "watch":
				c.watchFiles = true
			}
		}
		if len(sources) > 0 {
			c.sourcePaths = append(c.sourcePaths, sources...)
		}
	} else {
		log.Printf("error: %s", err.Error())
		usage()
	}

	// Basic sanity checks.
	if c.encoding != "" {
		if cs := charset.Info(c.encoding); cs == nil {
			log.Printf("error: Charset %q is unsupported.", c.encoding)
			usageCharsets()
		}
	}
	if len(c.sourcePaths) == 0 {
		log.Print("error: Input sources not specified.")
		usage()
	}
	if c.watchFiles {
		if c.outFile == "-" {
			log.Fatal("error: Writing to standard output is unsupported in watch mode.")
		}
		// if c.logFiles {
		// 	log.Print("warning: File logging is unsupported in watch mode.")
		// }
		// if c.logStats {
		// 	log.Print("warning: Statistic logging is unsupported in watch mode.")
		// }
	}

	// Return the base configuration.
	return c
}

func (c *config) mergeStoryConfig(s *story) {
	if c.cmdline.formatID != "" {
		c.formatID = c.cmdline.formatID
	} else if s.twine2.format != "" {
		c.formatID = c.formats.getIDFromTwine2NameAndVersion(s.twine2.format, s.twine2.formatVersion)
		if c.formatID == "" {
			log.Printf("error: Story format named %q at version %q is not available.", s.twine2.format, s.twine2.formatVersion)
			usageFormats(c.formats)
		}
	} else {
		c.formatID = defaultFormatID
	}
	if !c.formats.hasByID(c.formatID) {
		log.Printf("error: Story format %q is not available.", c.formatID)
		usageFormats(c.formats)
	}

	if c.cmdline.startName != "" {
		c.startName = c.cmdline.startName
	} else if s.twine2.start != "" {
		c.startName = s.twine2.start
	} else {
		c.startName = defaultStartName
	}

	// Finalize the story setup.
	s.format = c.formats.getByID(c.formatID)
	s.twine2.options["debug"] = s.twine2.options["debug"] || c.testMode
}
