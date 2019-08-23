// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package spdxify

import (
	"fmt"

	"github.com/spdx/tools-golang/v0/idsearcher"
	"github.com/spdx/tools-golang/v0/spdx"
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
		fis[f.FileName] = f.LicenseConcluded
	}

	return fis, nil
}
