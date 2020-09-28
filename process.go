package process

import (
	"io"
	"io/ioutil"
	"os/exec"
	"time"

	"errors"
)

const (
	StatusCreated = "created"
	StatusRunning = "running"
	StatusExit    = "exit"
)

type Process struct {
	name    string
	args    []string
	stdErr  io.ReadWriteSeeker
	stdOut  io.ReadWriteSeeker
	status  string
	cmd     *exec.Cmd
	exit    chan error
	started time.Time
	Restart bool
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
	return p.started
}

func (p *Process) StdOut() (string, error) {
	return fileToString(p.stdOut)
}

func (p *Process) StdErr() (string, error) {
	return fileToString(p.stdErr)
}

func (p *Process) watch() {
	err := p.cmd.Wait()
	p.status = StatusExit

	if exitErr, ok := err.(*exec.ExitError); ok {
		if !exitErr.Exited() {
			p.exit <- err
			return
		}
	}

	if !p.Restart {
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
	return p.cmd.Process.Pid
}

func (p *Process) Status() string {
	return p.status
}

func (p *Process) Terminate() error {
	if p.status == StatusRunning {
		p.Restart = false
		return p.cmd.Process.Kill()
	}

	return errors.New("process is not running")
}

func (p *Process) Start() error {
	p.cmd = exec.Command(p.name, p.args...)
	p.cmd.Stderr = p.stdErr
	p.cmd.Stdout = p.stdOut

	if err := p.cmd.Start(); err != nil {
		return err
	}

	p.status = StatusRunning
	p.started = time.Now()

	go p.watch()

	return nil
}

func CreateProcess(stdOut, stdErr io.ReadWriteSeeker, name string, args ...string) *Process {
	return &Process{
		name:    name,
		args:    args,
		status:  StatusCreated,
		Restart: true,
		exit:    make(chan error),
		stdErr:  stdErr,
		stdOut:  stdOut,
	}
}
