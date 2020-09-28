package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/bolsunovskyi/process"
)

var (
	logsFolder string
	signals    chan os.Signal
	commands   map[string]Command
)

func init() {
	flag.StringVar(&logsFolder, "l", "logs", "logs folder location (default: './logs')")
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	signals = make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)
	go exitWait()

	commands = map[string]Command{
		"h": help{},
		"q": quit{},
		"r": run{},
		"l": list{},
		"k": kill{},
	}
}

func exitWait() {
	<-signals
	log.Println("force exit")
	os.Exit(0)
}

type Command interface {
	Exec(scanner *bufio.Scanner, manager *process.Manager) error
}

func main() {
	if _, err := os.Stat(logsFolder); err != nil {
		if err := os.MkdirAll(logsFolder, 0775); err != nil {
			log.Fatalln(err)
		}
	}

	p, err := process.CreateManager("logs")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Print("Process Manager, type 'h' for help\nCommand: ")
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		s := strings.TrimSpace(scanner.Text())

		if cmd, ok := commands[s]; ok {
			if err := cmd.Exec(scanner, p); err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Command not found, type 'h' for help")
		}

		fmt.Print("Command: ")
	}
}
