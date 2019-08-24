// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package main

import (
	"log"

	"github.com/swinslow/spdxify/pkg/spdxify"
)

func main() {
	acfg, args, err := spdxify.ParseArgs()
	if err != nil {
		log.Fatalf("error parsing command-line arguments: %v", err)
	}

	repoPath := args[0]

	// load configuration file
	cfg, err := spdxify.LoadConfig(acfg.CfgPath)
	if err != nil {
		log.Fatalf("error loading config file: %v", err)
	}

	// get slice of files to test
	sel, err := spdxify.SelectFiles(cfg, repoPath)
	if err != nil {
		log.Fatalf("error preparing files for analysis: %v", err)
	}

	// search for files with existing IDs
	fis, err := spdxify.SearchSPDXIDs(repoPath, cfg.Skip.Dirs)
	if err != nil {
		log.Fatalf("error searching for existing SPDX IDs: %v", err)
	}

	// set which licenses to apply to which files
	var fds []spdxify.FileData
	if acfg.LicenseID != "" {
		// license was specified on command line
		fds, err = spdxify.ChooseActions(cfg, sel, fis, acfg)
		if err != nil {
			log.Fatalf("error determining actions: %v", err)
		}
	} else {
		// NOT YET IMPLEMENTED
		log.Fatalf("license to apply was not specified on command line (-l flag); license choice from config file or SPDX file is not yet implemented")
	}

	log.Printf("selected files: %v", sel)
	log.Printf("detected SPDX IDs: %v", fis)
	log.Printf("file data and actions: %#v", fds)

}
