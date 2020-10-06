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
	name        string
	args        []string
	stdErr      io.ReadWriteSeeker
	stdOut      io.ReadWriteSeeker
	status      string
	statusLock  sync.RWMutex
	cmd         *exec.Cmd
	cmdLock     sync.RWMutex
	exit        chan error
	started     time.Time
	startedLock sync.RWMutex
	restart     bool
	restartLock sync.RWMutex
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
	p.startedLock.RLock()
	defer p.startedLock.RUnlock()
	return p.started
}

func (p *Process) StdOut() (string, error) {
	return fileToString(p.stdOut)
}

func (p *Process) StdErr() (string, error) {
	return fileToString(p.stdErr)
}

func (p *Process) watch() {
	p.cmdLock.RLock()
	err := p.cmd.Wait()
	p.cmdLock.RUnlock()

	p.statusLock.Lock()
	p.status = StatusExit
	p.statusLock.Unlock()

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
	p.cmdLock.RLock()
	defer p.cmdLock.RUnlock()
	return p.cmd.Process.Pid
}

func (p *Process) Status() string {
	p.statusLock.RLock()
	s := p.status
	p.statusLock.RUnlock()
	return s
}

func (p *Process) Terminate() error {
	if p.Status() == StatusRunning {
		p.SetRestart(false)
		p.cmdLock.RLock()
		err := p.cmd.Process.Kill()
		p.cmdLock.RUnlock()
		return err
	}

	return errors.New("process is not running")
}

func (p *Process) Start() error {
	p.cmdLock.Lock()
	p.cmd = exec.Command(p.name, p.args...)
	p.cmd.Stderr = p.stdErr
	p.cmd.Stdout = p.stdOut

	if err := p.cmd.Start(); err != nil {
		return err
	}

	p.cmdLock.Unlock()

	p.statusLock.Lock()
	p.status = StatusRunning
	p.statusLock.Unlock()

	p.startedLock.Lock()
	p.started = time.Now()
	p.startedLock.Unlock()

	go p.watch()

	return nil
}

func CreateProcess(stdOut, stdErr io.ReadWriteSeeker, name string, args ...string) *Process {
	return &Process{
		name:    name,
		args:    args,
		status:  StatusCreated,
		restart: true,
		exit:    make(chan error),
		stdErr:  stdErr,
		stdOut:  stdOut,
	}
}
