package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/bolsunovskyi/process"
)

type quit struct{}

func (quit) Exec(_ *bufio.Scanner, m *process.Manager) error {
	if err := m.ShutDown(); err != nil {
		return err
	}

	fmt.Println("exit")
	os.Exit(0)
	return nil
}
