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

func scanPid(scanner *bufio.Scanner) (int, error) {
	fmt.Print("Enter process PID: ")
	if !scanner.Scan() {
		return 0, errors.New("wrong process PID")
	}

	cmd := strings.TrimSpace(scanner.Text())
	if cmd == "" {
		return 0, errors.New("wrong process PID")
	}

	pid, err := strconv.Atoi(cmd)
	if err != nil {
		return 0, errors.New("pid must be a number")
	}

	return pid, nil
}

func (kill) Exec(scanner *bufio.Scanner, manager *process.Manager) error {
	pid, err := scanPid(scanner)
	if err != nil {
		return err
	}

	if err := manager.TerminateProcess(pid); err != nil {
		return err
	}

	fmt.Println("process killed")

	return nil
}
