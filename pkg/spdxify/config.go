// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

// Package spdxify contains data structures and functions to add SPDX short-form
// license identifiers (https://spdx.org/ids) to source code files in an
// automated manner.
package spdxify

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

// Config defines the configuration for a run of spdxify.
type Config struct {
	// Filetypes is a mapping of file suffixes to ConfigFiletype entries. The
	// suffix is compared as-is, without evaluating whether it comprises a true
	// "file extension", so typically each key in the Filetypes map should
	// begin with a period.
	Filetypes map[string]ConfigFiletype `json:"filetypes,omitempty"`

	// Skip is a ConfigSkip section as defined below.
	Skip ConfigSkip `json:"skip,omitempty"`
}

// ConfigFiletype defines the configuration for each value in the "filetypes"
// mapping section of Config.
type ConfigFiletype struct {
	// Comment defines the comment format for this file type. It should include
	// the string "SPDX" where the "SPDX-License-Identifier: ABC" data should be
	// placed, and may include arbitrary characters before or after it. Comment
	// may not include line breaks as the license identifier must be on a single
	// line. If the Comment field is omitted or empty, the SPDX identifier line
	// will be added as plain text without any comment delimiters.
	Comment string `json:"comment,omitempty"`

	// SkipFirstIfPrefix may optionally be included to indicate that for this
	// file type, the SPDX identifier should NOT be added in the first line if
	// the first line begins with the indicated prefix. For example, shell
	// scripts should likely have a value of "#!" for SkipFirstIfPrefix. If
	// omitted or empty, the SPDX identifier will be added in the first line.
	SkipFirstIfPrefix string `json:"skipFirstIfPrefix,omitempty"`
}

// ConfigSkip defines the configuration for the "skip" section of Config.
type ConfigSkip struct {
	// Filetypes is a slice of strings for files that should be skipped if they
	// end with the indicated suffix. The suffix is compared as-is, without
	// evaluating whether it comprises a true "file extension", so typically
	// each value in the Filetypes slice should begin with a period.
	Filetypes []string `json:"filetypes,omitempty"`

	// Dirs is a slice of string glob patterns for directories or files that
	// should be skipped, regardless of their file names. It uses the same path
	// patterns as spdx/tools-golang (see, e.g.,
	// https://github.com/spdx/tools-golang/blob/master/examples/3-build/example_build.go#L61)
	Dirs []string `json:"dirs,omitempty"`
}

// LoadConfig attempts to load the specified config file; or if an empty string
// is passed, it attempts to load config from ~/.spdxify.json. It returns the
// parsed Config if successful or error if unable to load.
func LoadConfig(configFile string) (*Config, error) {
	// get default config file from user's home directory if none specified
	if configFile == "" {
		usr, err := user.Current()
		if err != nil {
			return nil, fmt.Errorf("cannot get current user for home directory config file: %v", err)
		}
		configFile = filepath.Join(usr.HomeDir, ".spdxify.json")
	}

	// load config file from disk into byte slice
	js, err := os.Open(configFile)
	if err != nil {
		return nil, fmt.Errorf("error opening config file %s: %v", configFile, err)
	}
	defer js.Close()
	b, err := ioutil.ReadAll(js)
	if err != nil {
		return nil, fmt.Errorf("error reading config file %s: %v", configFile, err)
	}

	// and parse as JSON
	var cfg Config
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON from config file %s: %v", configFile, err)
	}

	return &cfg, nil
}
