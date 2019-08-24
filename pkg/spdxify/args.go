// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package spdxify

import (
	"flag"
	"fmt"
)

// ArgsConfig defines the command-line configuration options at run time.
type ArgsConfig struct {
	// CfgPath is the path to the configuration file.
	CfgPath string

	// LicenseID is the chosen license ID, if one is declared on command line.
	LicenseID string

	// Extension is the chosen file extension for search+apply, if one is
	// declared on command line for a one-shot.
	Extension string

	// Comment is the format to use for the chosen extension, if one is declared
	// on command line for a one-shot. If no format is identified and no
	// corresponding format appears in the config file, an error will be
	// returned when trying to apply to a file.
	Comment string

	// PrefixToSkip is the prefix to check on the first line of a file, if one
	// is declared on command line for a one-shot.
	Prefix string

	// DryRun is whether this should be just a dry-run (e.g., test for no errors
	// but don't actually make changes to files).
	DryRun bool

	// Verbose is whether verbose log messages should be returned.
	Verbose bool
}

// ParseArgs parses the command-line arguments and returns an ArgsConfig struct
// and location-based args slice, or error if parsing is incorrect.
func ParseArgs() (*ArgsConfig, []string, error) {
	// set up command line flags
	cfgFlag := flag.String("c", "", "path to configuration file (defaults to ~/.spdxify.json")
	licFlag := flag.String("l", "", "license to apply, if same for all")
	extFlag := flag.String("e", "", "file extensions to apply (for one-shot)")
	cmtFlag := flag.String("t", "", "comment format (for one-shot)")
	prefixFlag := flag.String("p", "", "prefix for skipping first line (for one-shot)")
	dryFlag := flag.Bool("d", false, "dry run -- do not actually add identifiers")
	verboseFlag := flag.Bool("v", false, "verbose")

	// parse command line
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		return nil, nil, fmt.Errorf("Usage: spdxify [-cletpdv] PATH-TO-REPO")
	}

	return &ArgsConfig{
		CfgPath:   *cfgFlag,
		LicenseID: *licFlag,
		Extension: *extFlag,
		Comment:   *cmtFlag,
		Prefix:    *prefixFlag,
		DryRun:    *dryFlag,
		Verbose:   *verboseFlag,
	}, args, nil
}
