package main

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/bolsunovskyi/process"
)

type kill struct{}

func (kill) Exec(scanner *bufio.Scanner, manager *process.Manager) error {
	fmt.Print("Enter process PID: ")
	if !scanner.Scan() {
		return errors.New("wrong process PID")
	}

	cmd := strings.TrimSpace(scanner.Text())
	if cmd == "" {
		return errors.New("wrong process PID")
	}

	pid, err := strconv.Atoi(cmd)
	if err != nil {
		return errors.New("pid must be a number")
	}

	if err := manager.TerminateProcess(pid); err != nil {
		return err
	}

	fmt.Println("process killed")

	return nil
}
