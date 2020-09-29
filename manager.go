package process

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"time"
)

type ManagerConfig struct {
	LogsFolder        string
	RenewOldProcesses bool
	ProcessesListFile string
}

type Manager struct {
	ManagerConfig

	processes         []*ManagerProcess
	signals           chan os.Signal
	renewOldProcesses bool
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
	logFilePath := fmt.Sprintf("%s/%s-%s.log", m.LogsFolder, time.Now().Format("2006-01-02"),
		path.Base(name))
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

//GetProcess get process by pid
func (m *Manager) GetProcess(pid int) (*ManagerProcess, error) {
	for _, p := range m.processes {
		if pid == p.PID() {
			return p, nil
		}
	}

	return nil, errors.New("process with such pid not found")
}

//ShutDown - kill all processes, flush logs and save process state for renewal
func (m *Manager) ShutDown() error {
	if m.RenewOldProcesses {
		procFile, err := os.Create(m.ProcessesListFile)
		if err != nil {
			return err
		}

		for _, proc := range m.processes {
			if _, err := procFile.WriteString(
				fmt.Sprintf("%s %s\n", proc.Name(), strings.Join(proc.Args(), " "))); err != nil {
				return err
			}
		}

		if err := procFile.Close(); err != nil {
			return err
		}
	}

	for _, proc := range m.processes {
		proc.Terminate()
	}

	return nil
}

func (m *Manager) renewProcesses() error {
	if m.RenewOldProcesses && m.ProcessesListFile != "" {
		f, err := os.Open(m.ProcessesListFile)
		if err != nil {
			return nil
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			cmd := strings.Split(scanner.Text(), " ")
			var args []string
			if len(cmd) > 1 {
				args = cmd[1:]
			}

			if _, err := m.AddProcess(cmd[0], args...); err != nil {
				return err
			}
		}
	}

	return nil
}

//CreateManager init main library structure
func CreateManager(config ManagerConfig) (*Manager, error) {
	config.LogsFolder = path.Clean(config.LogsFolder)
	config.ProcessesListFile = path.Clean(config.ProcessesListFile)

	if config.ProcessesListFile == "" || config.ProcessesListFile == "." {
		config.ProcessesListFile = "processes.txt"
	}

	if _, err := os.Stat(config.LogsFolder); err != nil {
		return nil, err
	}

	manager := Manager{
		ManagerConfig: config,
		processes:     []*ManagerProcess{},
		signals:       make(chan os.Signal, 1),
	}

	if err := manager.renewProcesses(); err != nil {
		return nil, err
	}

	return &manager, nil
}
