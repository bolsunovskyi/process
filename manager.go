package process

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path"
	"time"
)

type Manager struct {
	processes  map[int]ManagerProcess
	logsFolder string
	signals    chan os.Signal
}

type ManagerProcess struct {
	Process

	logPath string
	logFile *os.File
	PID     int
}

func (m *Manager) AddProcess(name string, args ...string) (*ManagerProcess, error) {
	logFilePath := fmt.Sprintf("%s/%s-%s.log", m.logsFolder, time.Now().Format("2006-01-02"), name)
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	p := CreateProcess(logFile, logFile, name, args...)
	if err := p.Start(); err != nil {
		return nil, err
	}

	mp := ManagerProcess{
		Process: p,
		logPath: logFilePath,
		logFile: logFile,
		PID:     p.cmd.Process.Pid,
	}

	m.processes[mp.PID] = mp

	return &mp, nil
}

func (m *Manager) GetProcesses() map[int]ManagerProcess {
	return m.processes
}

func (m *Manager) TerminateProcess(pid int) error {
	proc, ok := m.processes[pid]
	if !ok {
		return errors.New("process with such pid not found")
	}

	if err := proc.Terminate(); err != nil {
		if err.Error() != "os: process already finished" {
			return err
		}
	}

	if err := proc.logFile.Close(); err != nil {
		return err
	}

	delete(m.processes, pid)

	return nil
}

func (m *Manager) ShutDown() error {
	for pid := range m.processes {
		if err := m.TerminateProcess(pid); err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) listenSignals() {
	<-m.signals
	m.ShutDown()
}

func CreateManager(logsFolder string) (*Manager, error) {
	logsFolder = path.Clean(logsFolder)

	if _, err := os.Stat(logsFolder); err != nil {
		return nil, err
	}

	manager := Manager{
		logsFolder: logsFolder,
		processes:  make(map[int]ManagerProcess),
		signals:    make(chan os.Signal, 1),
	}

	signal.Notify(manager.signals, os.Interrupt)

	go manager.listenSignals()

	return &manager, nil
}
