package process

import (
	"testing"
	"time"
)

func TestCreateManager(t *testing.T) {
	manager, err := CreateManager("logs")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := manager.AddProcess("uname", "-r"); err != nil {
		t.Fatal(err)
	}

	if err := manager.ShutDown(); err != nil {
		t.Fatal(err)
	}
}

func TestManager_GetProcesses(t *testing.T) {
	manager, err := CreateManager("logs")
	if err != nil {
		t.Fatal(err)
	}

	proc, err := manager.AddProcess("sleep", "10")
	if err != nil {
		t.Fatal(err)
	}

	allProcs := manager.GetProcesses()
	if _, ok := allProcs[proc.PID]; !ok {
		t.Fatal("processes not found in all processes map")
	}
}

func TestCreateManagerLogs(t *testing.T) {
	manager, err := CreateManager("logs")
	if err != nil {
		t.Fatal(err)
	}

	proc, err := manager.AddProcess("uname", "-r")
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Millisecond * 500)

	output, err := proc.StdOut()
	if err != nil {
		t.Fatal(err)
	}

	if output == "" {
		t.Fatal("output is empty")
	}

	if err := manager.ShutDown(); err != nil {
		t.Fatal(err)
	}
}

func TestManagerFailures(t *testing.T) {
	_, err := CreateManager("/mike/")
	if err == nil {
		t.Fatal("no error on wrong folder")
	}

	manager, err := CreateManager("/proc")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := manager.AddProcess("echo", "john"); err == nil {
		t.Fatal("no error on wrong folder path")
	}

	manager, err = CreateManager("logs")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := manager.AddProcess("John", "Doe"); err == nil {
		t.Fatal("no error on wrong process name")
	}

	if err := manager.TerminateProcess(-1); err == nil {
		t.Fatal("no error on termination by fake pid")
	}
}
