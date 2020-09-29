package main

import (
	"flag"
	"log"
	"os"
	"time"
)

var (
	crashTimeout, messageInterval           uint
	enableStdOut, enableStdErr, successExit bool
)

func init() {
	flag.UintVar(&crashTimeout, "ct", 0, "crash timeout (seconds)")
	flag.UintVar(&messageInterval, "i", 1, "message interval")
	flag.BoolVar(&enableStdOut, "o", true, "enable stdout messages (default true)")
	flag.BoolVar(&enableStdErr, "e", true, "enable stderr messages (default true)")
	flag.BoolVar(&successExit, "s", false, "perform one cycle iteration and exit")

	flag.Parse()
}

func main() {
	stdErrLog := log.New(os.Stderr, "out", log.LstdFlags)
	stdOutLog := log.New(os.Stdout, "err", log.LstdFlags)

	if crashTimeout > 0 && !successExit {
		go func() {
			time.Sleep(time.Second * time.Duration(crashTimeout))
			stdErrLog.Println("crash timeout reach")
			os.Exit(1)
		}()
	}

	for i := 0; ; i++ {
		if enableStdOut {
			stdOutLog.Printf("stdout %d\n", i)
		}

		if enableStdErr {
			stdErrLog.Printf("stderr %d\n", i)
		}

		if successExit {
			os.Exit(0)
		}

		time.Sleep(time.Second * time.Duration(messageInterval))
	}
}
