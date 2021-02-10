/*
	Copyright © 2014–2021 Thomas Michael Edwards. All rights reserved.
	Use of this source code is governed by a Simplified BSD License which
	can be found in the LICENSE file.
*/

package main

import (
	"fmt"
	"runtime"
)

// versionInfo contains build version information.
type versionInfo struct {
	major uint64
	minor uint64
	patch uint64
	pre   string
}

var (
	// tweegoVersion holds the current version info.
	tweegoVersion = versionInfo{
		major: 2,
		minor: 1,
		patch: 1,
		pre:   "",
	}
	// tweegoVersion holds the build ID.
	tweegoBuild = ""
	// tweegoVersion holds the build date.
	tweegoDate = ""
)

// String returns the full version string (version, date, and platform).
func (v versionInfo) String() string {
	date := tweegoDate
	if date != "" {
		date = " (" + date + ")"
	}
	return fmt.Sprintf("version %s%s [%s]", v.Version(), date, v.Platform())
}

// Version returns the SemVer version string.
func (v versionInfo) Version() string {
	pre := v.pre
	if pre != "" {
		pre = "-" + pre
	}
	build := tweegoBuild
	if build != "" {
		build = "+" + build
	}
	return fmt.Sprintf(
		"%d.%d.%d%s%s",
		v.major,
		v.minor,
		v.patch,
		pre,
		build,
	)
}

// Date returns the build date string.
func (v versionInfo) Date() string {
	return tweegoDate
}

// Platform returns the OS/Arch pair.
func (v versionInfo) Platform() string {
	return runtime.GOOS + "/" + runtime.GOARCH
}
