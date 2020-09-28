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
	processes  []*ManagerProcess
	logsFolder string
	signals    chan os.Signal
}

type ManagerProcess struct {
	Process
	logPath string
	logFile *os.File
}

func (m *ManagerProcess) LogPath() string {
	return m.logPath
}

func (m *Manager) AddProcess(name string, args ...string) (*ManagerProcess, error) {
	logFilePath := fmt.Sprintf("%s/%s-%s.log", m.logsFolder, time.Now().Format("2006-01-02"), name)
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	p := CreateProcess(logFile, logFile, name, args...)

	mp := ManagerProcess{
		Process: *p,

		logPath: logFilePath,
		logFile: logFile,
	}

	if err := mp.Start(); err != nil {
		return nil, err
	}

	m.processes = append(m.processes, &mp)

	return &mp, nil
}

func (m *Manager) GetProcesses() []*ManagerProcess {
	return m.processes
}

func (m *Manager) TerminateProcess(pid int) error {
	for i, proc := range m.processes {
		if proc.PID() == pid {
			if err := proc.Terminate(); err != nil {
				return err
			}

			if err := proc.logFile.Close(); err != nil {
				return err
			}

			m.processes = append(m.processes[:i], m.processes[i+1:]...)

			return nil
		}
	}

	return errors.New("process with such pid not found")
}

func (m *Manager) ShutDown() error {
	for _, proc := range m.processes {
		proc.Terminate()
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
		processes:  []*ManagerProcess{},
		signals:    make(chan os.Signal, 1),
	}

	signal.Notify(manager.signals, os.Interrupt, os.Kill)

	go manager.listenSignals()

	return &manager, nil
}
