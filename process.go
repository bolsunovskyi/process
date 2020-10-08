package process

import (
	"io"
	"io/ioutil"
	"os/exec"
	"sync"
	"time"

	"errors"
)

const (
	StatusCreated = "created"
	StatusRunning = "running"
	StatusExit    = "exit"
)

type Process struct {
	name   string
	args   []string
	stdErr io.ReadWriteSeeker
	stdOut io.ReadWriteSeeker

	command struct {
		status  string
		cmd     *exec.Cmd
		started time.Time
	}
	commandLock sync.RWMutex

	restart     bool
	restartLock sync.RWMutex

	exit chan error
}

func (p *Process) SetRestart(restart bool) {
	p.restartLock.Lock()
	p.restart = restart
	p.restartLock.Unlock()
}

func (p *Process) Restart() bool {
	p.restartLock.RLock()
	defer p.restartLock.RUnlock()
	return p.restart
}

func fileToString(f io.ReadWriteSeeker) (string, error) {
	f.Seek(0, 0)
	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}
	f.Seek(0, io.SeekEnd)

	return string(bytes), nil
}

func (p *Process) Name() string {
	return p.name
}

func (p *Process) Args() []string {
	return p.args
}

func (p *Process) Started() time.Time {
	p.commandLock.RLock()
	defer p.commandLock.RUnlock()
	return p.command.started
}

func (p *Process) StdOut() (string, error) {
	return fileToString(p.stdOut)
}

func (p *Process) StdErr() (string, error) {
	return fileToString(p.stdErr)
}

func (p *Process) watch() {
	p.commandLock.RLock()
	err := p.command.cmd.Wait()
	p.commandLock.RUnlock()

	p.commandLock.Lock()
	p.command.status = StatusExit
	p.commandLock.Unlock()

	if exitErr, ok := err.(*exec.ExitError); ok {
		if !exitErr.Exited() {
			p.exit <- err
			return
		}
	}

	if !p.Restart() {
		p.exit <- err
		return
	}

	if err = p.Start(); err != nil {
		p.exit <- err
	}
}

func (p *Process) Wait() error {
	return <-p.exit
}

func (p *Process) PID() int {
	p.commandLock.RLock()
	defer p.commandLock.RUnlock()
	return p.command.cmd.Process.Pid
}

func (p *Process) Status() string {
	p.commandLock.RLock()
	defer p.commandLock.RUnlock()
	return p.command.status
}

func (p *Process) Terminate() error {
	if p.Status() == StatusRunning {
		p.SetRestart(false)
		p.commandLock.RLock()
		err := p.command.cmd.Process.Kill()
		p.commandLock.RUnlock()

		return err
	}

	return errors.New("process is not running")
}

func (p *Process) Start() error {
	p.commandLock.Lock()
	p.command.cmd = exec.Command(p.name, p.args...)
	p.command.cmd.Stderr = p.stdErr
	p.command.cmd.Stdout = p.stdOut

	if err := p.command.cmd.Start(); err != nil {
		return err
	}

	p.command.status = StatusRunning
	p.command.started = time.Now()
	p.commandLock.Unlock()

	go p.watch()

	return nil
}

func CreateProcess(stdOut, stdErr io.ReadWriteSeeker, name string, args ...string) *Process {
	return &Process{
		name: name,
		args: args,
		command: struct {
			status  string
			cmd     *exec.Cmd
			started time.Time
		}{status: StatusCreated, cmd: nil, started: time.Now()},
		exit:    make(chan error),
		restart: true,
		stdErr:  stdErr,
		stdOut:  stdOut,
	}
}
