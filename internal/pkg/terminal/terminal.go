//
// terminal/terminal.go
// terminal
//
// Created by steve on 2019-08-30.
// Copyright © 2019 Klassen Software Solutions. All rights reserved.
// Permission is hereby granted for use under the MIT License (https://opensource.org/licenses/MIT).
//

package terminal

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/docopt/docopt-go"
	"github.com/gorilla/websocket"
	"github.com/klassen-software-solutions/gocontract/contract"
	"os"
	"strings"
	"time"
)

const (
	retryPause = time.Second * 2
	prompt     = ""
)

// Terminal provides the ability to run a web socket terminal application
type Terminal struct {
	uri           string
	isQuiet       bool
	isAutoConnect bool
	isPrettyPrint bool
}

// NewTerminal creates a new terminal from the command line arguments
func NewTerminal(opts docopt.Opts) *Terminal {
	t := Terminal{}
	t.uri = opts["URI"].(string)
	t.isQuiet = opts["--quiet"].(bool)
	t.isAutoConnect = opts["--retry"].(bool)
	t.isPrettyPrint = opts["--pretty"].(bool)

	contract.Postconditions(t.uri != "")
	return &t
}

// Run runs the terminal
func (t *Terminal) Run() error {
	contract.Preconditions(t.uri != "")

	for {
		err := t.runSession()
		if err != nil {
			if !t.isAutoConnect {
				return err
			}

			t.displayStatus(fmt.Sprintf("Connection failed, retry in %vs", retryPause.Seconds()))
			time.Sleep(retryPause)
			continue
		}
	}
}

func (t *Terminal) runSession() error {
	contract.Preconditions(t.uri != "")

	t.displayStatus(fmt.Sprintf("Connecting to %s...", t.uri))
	conn, err := connect(t.uri)
	if err != nil {
		return err
	}

	contract.Conditions(conn != nil)
	defer CloseAndIgnore(conn)

	t.displayStatus("Connected")

	shouldStop := false
	conn.SetCloseHandler(func(code int, text string) error {
		shouldStop = true
		return nil
	})

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Printf("!! Error reading: %v\n", err)
				shouldStop = true
			}

			t.displayMessage(string(message))
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	for shouldStop == false {
		fmt.Print(prompt)
		message, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		err = sendMessage(conn, &message)
		if err != nil {
			return err
		}
	}

	return nil
}

func connect(uri string) (*websocket.Conn, error) {
	if uri == "" {
		return nil, fmt.Errorf("missing uri")
	}

	conn, _, err := websocket.DefaultDialer.Dial(uri, nil)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func sendMessage(conn *websocket.Conn, message *string) error {
	return conn.WriteMessage(websocket.TextMessage, []byte(strings.TrimRight(*message, "\n")))
}

func (t *Terminal) displayStatus(message string) {
	if !t.isQuiet {
		fmt.Printf("! %s\n", message)
	}
}

func (t *Terminal) displayMessage(message string) {
	if t.isPrettyPrint {
		message = prettyPrintJSON(message)
	}

	fmt.Println(message)
}

func prettyPrintJSON(message string) string {
	bytes := []byte(message)
	if !json.Valid(bytes) {
		return message
	}

	m := map[string]interface{}{}
	err := json.Unmarshal([]byte(message), &m)
	if err != nil {
		return message
	}

	newMessage, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return message
	}

	return string(newMessage)
}
