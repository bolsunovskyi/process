package process

import (
	"bytes"
	"os/exec"
)

const (
	StatusCreated = "created"
	StatusRunning = "running"
	StatusExit = "exit"
)

type Process struct {
	name string
	args []string
	stdErr bytes.Buffer
	stdOut bytes.Buffer
	status string
	cmd *exec.Cmd
	Restart bool
	exit chan error
}

func (p *Process) StdOut() string {
	return p.stdOut.String()
}

func (p *Process) StdErr() string {
	return p.stdErr.String()
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

func (p *Process) Kill() error {
	p.Restart = false

	return p.cmd.Process.Kill()
}

func (p *Process) Start() error {
	p.cmd = exec.Command(p.name, p.args...)
	p.cmd.Stderr = &p.stdErr
	p.cmd.Stdout = &p.stdOut

	if err := p.cmd.Start(); err != nil {
		return err
	}

	p.status = StatusRunning

	go p.watch()

	return nil
}

func CreateProcess(name string, args ...string) (*Process, error) {
	p := Process{
		name: name,
		args: args,
		status: StatusCreated,
		Restart: true,
		exit: make(chan error),
	}

	return &p, nil
}