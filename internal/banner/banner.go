// Package banner prints the program identification, attribution, and
// disclaimer text that Inspired Trek73 displays at startup and on exit.
package banner

import (
	"fmt"
	"io"
)

// text is the shared attribution and disclaimer block shown both at
// program start and program exit, per project requirements.
const text = `Inspired Trek73

A modern remake of the classic Trek73 game

Updated to modern operating systems and hardware by
Peter S. Lee
AppliedInspiration.com

The original Trek73 game was created by William K. Char, Perry Lee, and
Dan Gee. It was later corrected, completed, and enhanced by Jeff Okamoto
and Peter Yee, whose FreeBSD-era revision served as the reference source
for this modernization.

This implementation Copyright 2026 Peter S. Lee (AppliedInspiration.com)
All Rights Reserved.

This software is made available under the permissive MIT License.

This program is provided "as is" without warranties or guarantees of any kind.`

// Print writes the startup/exit banner to w.
func Print(w io.Writer) {
	fmt.Fprintln(w, text)
}
