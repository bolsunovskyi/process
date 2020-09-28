package main

import (
	"bufio"
	"fmt"

	"github.com/bolsunovskyi/process"
)

type help struct{}

func (help) Exec(_ *bufio.Scanner, _ *process.Manager) error {
	fmt.Println(`	r - run process
	k - kill process
	l - list process
	q - kill all processes and quit`)

	return nil
}
