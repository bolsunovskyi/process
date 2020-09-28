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
	logsFolder, processesFile string
	signals                   chan os.Signal
	commands                  map[string]Command
	renewOldProcesses         bool
)

func init() {
	flag.StringVar(&logsFolder, "l", "logs", "logs folder location (default: './logs')")
	flag.BoolVar(&renewOldProcesses, "r", true, "restart old processes after manager start")
	flag.StringVar(&processesFile, "p", "processes.txt", "file to store processes to be restored")
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	signals = make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)

	commands = map[string]Command{
		"h": help{},
		"q": quit{},
		"r": run{},
		"l": list{},
		"k": kill{},
	}
}

func exitWait(manager *process.Manager) {
	<-signals
	log.Println("force exit")
	if err := manager.ShutDown(); err != nil {
		log.Println(err)
	}
}

type Command interface {
	Exec(scanner *bufio.Scanner, manager *process.Manager) error
}

func scanInput(manager *process.Manager) {
	fmt.Print("Process Manager, type 'h' for help\nCommand: ")
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		s := strings.TrimSpace(scanner.Text())

		if cmd, ok := commands[s]; ok {
			if err := cmd.Exec(scanner, manager); err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Command not found, type 'h' for help")
		}

		fmt.Print("Command: ")
	}
}

func main() {
	if _, err := os.Stat(logsFolder); err != nil {
		if err := os.MkdirAll(logsFolder, 0775); err != nil {
			log.Fatalln(err)
		}
	}

	manager, err := process.CreateManager(process.ManagerConfig{
		LogsFolder:        logsFolder,
		RenewOldProcesses: renewOldProcesses,
		ProcessesListFile: processesFile,
	})

	if err != nil {
		log.Fatalln(err)
	}

	go scanInput(manager)

	exitWait(manager)
}
