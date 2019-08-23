// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package main

import (
	"flag"
	"log"

	"github.com/swinslow/spdxify/pkg/spdxify"
)

func main() {
	// set up command line flags
	cfgPtr := flag.String("c", "", "path to configuration file (defaults to ~/.spdxify.json")

	// parse command line
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		log.Fatalf("Usage: spdxify [-c] REPOPATH")
	}

	// load configuration file
	cfg, err := spdxify.LoadConfig(*cfgPtr)
	if err != nil {
		log.Fatalf("error loading config file: %v", err)
	}

	// get slice of files to test
	sel, err := spdxify.SelectFiles(cfg, args[0])
	if err != nil {
		log.Fatalf("error preparing files for analysis: %v", err)
	}

	log.Printf("selected files: %v", sel)
}
