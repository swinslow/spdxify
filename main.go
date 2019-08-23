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

	// load configuration file
	cfg, err := spdxify.LoadConfig(*cfgPtr)
	if err != nil {
		log.Fatalf("error loading config file: %v", err)
	}

	log.Printf("config: %#v", cfg)
}
