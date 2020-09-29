package main

import (
	"bufio"
	"fmt"

	"github.com/bolsunovskyi/process"
)

type printLog struct{}

func (printLog) Exec(scanner *bufio.Scanner, m *process.Manager) error {
	pid, err := scanPid(scanner)
	if err != nil {
		return err
	}

	p, err := m.GetProcess(pid)
	if err != nil {
		return err
	}

	out, err := p.StdOut()
	if err != nil {
		return err
	}

	fmt.Println(out)
	return nil
}
