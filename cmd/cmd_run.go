package main

import (
	"bufio"
	"errors"
	"fmt"
	"strings"

	"github.com/bolsunovskyi/process"
)

type run struct{}

func (run) Exec(scanner *bufio.Scanner, m *process.Manager) error {
	fmt.Print("Enter process name and arguments: ")
	if !scanner.Scan() {
		return errors.New("wrong process name")
	}

	cmd := strings.TrimSpace(scanner.Text())
	if cmd == "" {
		return errors.New("wrong process name")
	}

	cmdParts := strings.Split(cmd, " ")
	var args []string
	if len(cmdParts) > 1 {
		args = cmdParts[1:]
	}

	if _, err := m.AddProcess(cmdParts[0], args...); err != nil {
		return err
	}

	fmt.Println("process started")

	return nil
}
