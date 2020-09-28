package process

import (
	"testing"
	"time"
)

func TestCreateManager(t *testing.T) {
	manager, err := CreateManager(ManagerConfig{LogsFolder: "logs"})
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

func TestManagerProcess_Restart(t *testing.T) {
	manager, err := CreateManager(ManagerConfig{LogsFolder: "logs"})
	if err != nil {
		t.Fatal(err)
	}

	proc, err := manager.AddProcess("sleep", "2")
	if err != nil {
		t.Fatal(err)
	}
	pid1 := proc.PID()
	time.Sleep(time.Millisecond * 2500)

	pid2 := proc.PID()
	if pid1 == pid2 {
		t.Logf("p1: %d, pd2: %d", pid1, pid2)
		t.Fatal("pid is not changed")
	}
}

func TestManager_RenewProcesses(t *testing.T) {
	manager, err := CreateManager(ManagerConfig{LogsFolder: "logs", RenewOldProcesses: true})
	if err != nil {
		t.Fatal(err)
	}

	_, err = manager.AddProcess("sleep", "10")
	if err != nil {
		t.Fatal(err)
	}

	if err := manager.ShutDown(); err != nil {
		t.Fatal(err)
	}

	manager, err = CreateManager(ManagerConfig{LogsFolder: "logs", RenewOldProcesses: true})
	if err != nil {
		t.Fatal(err)
	}

	processes := manager.GetProcesses()
	if len(processes) == 0 {
		t.Fatal("no renewed processes")
	}

	if err := manager.ShutDown(); err != nil {
		t.Fatal(err)
	}
}

func TestManager_GetProcessesTerminateProcess(t *testing.T) {
	manager, err := CreateManager(ManagerConfig{LogsFolder: "logs"})
	if err != nil {
		t.Fatal(err)
	}

	proc, err := manager.AddProcess("sleep", "10")
	if err != nil {
		t.Fatal(err)
	}

	allProcs := manager.GetProcesses()
	found := false
	for _, p := range allProcs {
		if p.PID() == proc.PID() {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("processes not found in all processes list")
	}

	if err := manager.TerminateProcess(allProcs[0].PID()); err != nil {
		t.Fatal(err)
	}
}

func TestCreateManagerLogs(t *testing.T) {
	manager, err := CreateManager(ManagerConfig{LogsFolder: "logs"})
	if err != nil {
		t.Fatal(err)
	}

	proc, err := manager.AddProcess("uname", "-r")
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Millisecond * 500)

	if proc.LogPath() == "" {
		t.Fatal("no log path was set")
	}

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
	_, err := CreateManager(ManagerConfig{LogsFolder: "/mike/"})
	if err == nil {
		t.Fatal("no error on wrong folder")
	}

	manager, err := CreateManager(ManagerConfig{LogsFolder: "/proc"})
	if err != nil {
		t.Fatal(err)
	}

	if _, err := manager.AddProcess("echo", "john"); err == nil {
		t.Fatal("no error on wrong folder path")
	}

	manager, err = CreateManager(ManagerConfig{LogsFolder: "logs"})
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
