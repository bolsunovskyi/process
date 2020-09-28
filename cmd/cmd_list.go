package main

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/bolsunovskyi/process"
)

type list struct{}

func (list) Exec(_ *bufio.Scanner, manager *process.Manager) error {
	fmt.Println("PID\tNAME\tARGS\tSTATUS\tSTARTED\tLOG FILE\n")
	for _, proc := range manager.GetProcesses() {
		fmt.Printf("%d\t%s\t%s\t%s\t%s\t%s\n", proc.PID(), proc.Name(), strings.Join(proc.Args(), " "),
			proc.Status(), proc.Started().Format("2006-01-02 15:04:05"), proc.LogPath())
	}

	return nil
}
