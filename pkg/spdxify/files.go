// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package spdxify

import (
	"os"
	"path/filepath"
	"strings"
)

// SelectFiles takes a directory path and a Config, and returns a slice of file
// paths for the files that are selected to check. It will skip files that match
// the "skip" patterns in the Config, or that are not regular files (e.g.
// symbolic links).
func SelectFiles(cfg *Config, dirPath string) ([]string, error) {
	allFiles, err := getAllFilePaths(dirPath, cfg.Skip.Filetypes, cfg.Skip.Dirs)
	return allFiles, err
}

// The unexported function getAllFilePaths takes a path to a directory
// (including an optional slice of path patterns to ignore), and returns a slice
// of relative paths to all files in that directory and its subdirectories
// (excluding those that are ignored).
func getAllFilePaths(dirRoot string, suffixesIgnored []string, pathsIgnored []string) ([]string, error) {
	// paths is a _pointer_ to a slice -- not just a slice.
	// this is so that it can be appropriately modified by append
	// in the sub-function.
	paths := &[]string{}
	prefix := strings.TrimSuffix(dirRoot, "/")

	err := filepath.Walk(dirRoot, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// don't include path if it's excluded by suffix
		if isSuffixMatch(suffixesIgnored, path) {
			return nil
		}
		shortPath := strings.TrimPrefix(path, prefix)

		// don't include path if it should be ignored
		if pathsIgnored != nil && shouldIgnore(shortPath, pathsIgnored) {
			return nil
		}

		// don't include path if it's a directory
		if fi.IsDir() {
			return nil
		}
		// don't include path if it's a symbolic link
		if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
			return nil
		}
		// don't include path if it's an empty file
		if fi.Size() == 0 {
			return nil
		}

		// if we got here, record the path
		*paths = append(*paths, shortPath)
		return nil
	})

	return *paths, err
}

// The unexported function isSuffixMatch takes a slice of suffixes to be
// skipped, and a file path, and returns true if the file path ends with one of
// the suffixes.
func isSuffixMatch(suffixes []string, filePath string) bool {
	for _, suffix := range suffixes {
		if strings.HasSuffix(filePath, suffix) {
			return true
		}
	}
	return false
}

// The unexported function shouldIgnore compares a file path to a slice of file
// path patterns, and determines whether that file should be ignored because it
// matches any of those patterns.
func shouldIgnore(fileName string, pathsIgnored []string) bool {
	fDirs, fFile := filepath.Split(fileName)

	for _, pattern := range pathsIgnored {
		// split into dir(s) and filename
		patternDirs, patternFile := filepath.Split(pattern)
		patternDirStars := strings.HasPrefix(patternDirs, "**")
		if patternDirStars {
			patternDirs = patternDirs[2:]
		}

		// case 1: specific file
		if !patternDirStars && patternDirs == fDirs && patternFile != "" && patternFile == fFile {
			return true
		}

		// case 2: all files in specific directory
		if !patternDirStars && strings.HasPrefix(fDirs, patternDirs) && patternFile == "" {
			return true
		}

		// case 3: specific file in any dir
		if patternDirStars && patternDirs == "/" && patternFile != "" && patternFile == fFile {
			return true
		}

		// case 4: specific file in any matching subdir
		if patternDirStars && strings.Contains(fDirs, patternDirs) && patternFile != "" && patternFile == fFile {
			return true
		}

		// case 5: any file in any matching subdir
		if patternDirStars && strings.Contains(fDirs, patternDirs) && patternFile == "" {
			return true
		}

	}

	// if no match, don't ignore
	return false
}
