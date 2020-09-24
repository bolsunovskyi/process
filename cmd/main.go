package main

import (
	"github.com/bolsunovskyi/process"
	"log"
)

func main() {
	proc, err := process.CreateProcess("./example_process/example", "-ct=10")
	if err != nil {
		log.Fatalln(err)
	}

	if err := proc.Start(); err != nil {
		log.Fatalln(err)
	}

	if err := proc.Wait(); err != nil {
		log.Println(err)
	}
}
