//
// wsterm/main.go
// wsterm
//
// Created by steve on 2019-08-29.
// Copyright © 2019 Klassen Software Solutions. All rights reserved.
// Permission is hereby granted for use under the MIT License (https://opensource.org/licenses/MIT).
//

package main

import (
	"github.com/docopt/docopt-go"
	"github.com/klassen-software-solutions/wsterm"
	"github.com/klassen-software-solutions/wsterm/internal/pkg/terminal"
	"log"
	"os"
)

const (
	// note that usage is interpreted by docopt - i.e. this is code not just a string
	usage = `wsterm - web socket terminal

Usage:
  wsterm --help
  wsterm --version
  wsterm URI [--quiet] [--retry] [--pretty]

Options:
 -h, --help    print this message, then exit.
 -p, --pretty  if the received content is JSON, pretty print it.
 -q, --quiet   only display what is received, no status messages.
 -r, --retry   when connection fails, keep retrying to connect.
 --version     print the version, then exit.

Examples:
  wsterm ws://echo.websocket.org
`
)

func main() {
	opts, err := docopt.ParseArgs(usage, os.Args[1:], wsterm.Version)
	if err != nil {
		log.Fatalf("error parsing arguments: %v", err)
	}

	t := terminal.NewTerminal(opts)
	err = t.Run()
	if err != nil {
		_ = os.Stdout.Sync()
		log.Fatalf("error running: %v", err)
	}
}
