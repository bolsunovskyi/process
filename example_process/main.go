package main

import (
	"flag"
	"log"
	"os"
	"time"
)

var crashTimeout int

func init() {
	flag.IntVar(&crashTimeout, "ct", -1, "crash timeout (seconds)")
	flag.Parse()
}

func main() {
	stdErrLog := log.New(os.Stderr, "out", log.LstdFlags)
	stdOutLog := log.New(os.Stdout, "err", log.LstdFlags)

	if crashTimeout > 0 {
		go func() {
			time.Sleep(time.Second * time.Duration(crashTimeout))
			stdErrLog.Println("crash timeout reach")
			os.Exit(1)
		}()
	}

	for i:=0;;i++ {
		time.Sleep(time.Second * 2)
		stdOutLog.Printf("stdout %d\n", i)
		if i % 5 == 0 {
			stdErrLog.Printf("stderr %d\n", i)
		}
	}
}
