// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package spdxify

import (
	"fmt"

	"github.com/spdx/tools-golang/v0/idsearcher"
	"github.com/spdx/tools-golang/v0/spdx"
)

// ActionType is the action to be taken on a particular file, given its situation.
type ActionType int

const (
	// ActionSkip means we should skip this file
	ActionSkip ActionType = iota

	// ActionAdd means an identifier should be added
	ActionAdd

	// ActionError means there is an error and we should not proceed
	ActionError
)

// FileData tracks, for a given file, what we know about it and what action we
// want to take on it.
type FileData struct {
	// Name is the file's path, relative to the working root. Note that it
	// should NOT be modified after it is added to the working files map, since
	// its file path functions as its key in that map.
	Name string

	// FoundID is the SPDX license identifier that was found by searching the
	// file using spdx/tools-golang, if any.
	FoundID string

	// WantID is the SPDX license identifier that we want to insert for this file
	WantID string

	// Line is the line where we want to insert the SPDX identifier for this file
	Line int

	// Action is the type of action to be taken.
	Action ActionType
}

// SearchSPDXIDs takes a directory path and a list of directories to ignore, and
// searches each (unignored) file in that directory for an SPDX identifier. It
// returns a mapping of relative file path to identifier, with an empty string
// if no identifier is found.
func SearchSPDXIDs(dirPath string, pathsIgnored []string) (map[string]string, error) {
	idsearcherCfg := &idsearcher.Config{
		// a valid namespace prefix is not used, because this SPDX document will not
		// be used or accessible outside this function
		NamespacePrefix:     "N/A-internal-use-only",
		BuilderPathsIgnored: pathsIgnored,
	}

	pkgName := "internal-pkg"
	doc, err := idsearcher.BuildIDsDocument(pkgName, dirPath, idsearcherCfg)
	if err != nil {
		return nil, fmt.Errorf("error searching for existing SPDX IDs: %v", err)
	}

	// received SPDX document; look inside its package for file info
	var pkg *spdx.Package2_1
	for _, checkPkg := range doc.Packages {
		// make sure we got the expected package name -- we certainly should,
		// but just ensure idsearcher isn't doing anything weird on us
		if checkPkg.PackageName == pkgName {
			pkg = checkPkg
			break
		}
	}
	if pkg == nil {
		return nil, fmt.Errorf("error searching for existing SPDX IDs: package not found in document after IDs searched")
	}

	// and build the mapping of filename to detected ID
	fis := map[string]string{}
	for _, f := range pkg.Files {
		if f.LicenseConcluded == "NOASSERTION" {
			fis[f.FileName] = ""
		} else {
			fis[f.FileName] = f.LicenseConcluded
		}
	}

	return fis, nil
}

// ChooseActions walks through the selected files, comparing the requested
// license(s) against the existing license(s) from identifier search, and
// decides what action should be taken for each file. It does not proceed
// to take the action, but records it for future decision-making.
func ChooseActions(cfg *Config, selectedFiles []string, searchedFiles map[string]string, acfg *ArgsConfig) ([]FileData, error) {
	fds := []FileData{}

	for _, sf := range selectedFiles {
		fd := FileData{Name: sf}
		if foundID, ok := searchedFiles[sf]; ok {
			fd.FoundID = foundID
		}

		// determine what license ID we want
		var wantID string
		var err error
		if acfg.LicenseID != "" {
			wantID = acfg.LicenseID
		} else {
			wantID, err = getLicense(cfg, sf)
			if err != nil {
				return nil, fmt.Errorf("NOT YET IMPLEMENTED")
			}
		}
		fd.WantID = wantID

		// now decide what action to take
		if fd.WantID == "" {
			fd.Action = ActionSkip
		} else if fd.FoundID == fd.WantID {
			fd.Action = ActionSkip
		} else if fd.FoundID == "" {
			fd.Action = ActionAdd
		} else {
			fd.Action = ActionError
		}

		fds = append(fds, fd)
	}

	return fds, nil
}

// getLicense checks the Config and returns the license ID corresponding to the
// requested filename, or returns error if not found.
func getLicense(cfg *Config, filePath string) (string, error) {
	// NOT YET IMPLEMENTED
	return "", fmt.Errorf("NOT YET IMPLEMENTED")
}
