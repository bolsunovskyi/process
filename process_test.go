package process

import (
	"bytes"
	"os"
	"testing"
	"time"
)

type BufferSeeker struct {
	bytes.Buffer
}

func (BufferSeeker) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func TestCreateProcessNoRestart(t *testing.T) {
	var stdOut, stdErr BufferSeeker
	p := CreateProcess(&stdOut, &stdErr, "uname", "-r")
	p.Restart = false

	if p.Name() != "uname" {
		t.Fatal("wrong name in getter")
	}

	if len(p.Args()) == 0 || p.Args()[0] != "-r" {
		t.Fatal("wrong args in getter")
	}

	if err := p.Start(); err != nil {
		t.Fatal(err)
	}

	if err := p.Wait(); err != nil {
		t.Fatal(err)
	}

	output, err := p.StdOut()
	if err != nil {
		t.Fatal("unable to get stdout")
	}
	t.Log(output)

	if output == "" {
		t.Fatal("empty stdout")
	}

	if err := p.Terminate(); err == nil {
		t.Fatal("no error on trying to kill stopped process")
	}
}

func TestCreateProcessStdErr(t *testing.T) {
	var stdOut, stdErr BufferSeeker
	p := CreateProcess(&stdOut, &stdErr, "date", "-x")
	p.Restart = false

	if err := p.Start(); err != nil {
		t.Fatal(err)
	}

	if err := p.Wait(); err == nil {
		t.Fatal("no error on exit status 1")
	}

	output, err := p.StdErr()
	if err != nil {
		t.Fatal("unable to get stderr")
	}
	t.Log(output)
	if output == "" {
		t.Fatal("empty stderr")
	}
}

func TestCreateProcessRestart(t *testing.T) {
	var stdOut, stdErr BufferSeeker
	p := CreateProcess(&stdOut, &stdErr, "sleep", "1")
	if err := p.Start(); err != nil {
		t.Fatal(err)
	}

	started := p.Started()
	pid1 := p.PID()

	time.Sleep(time.Millisecond * 2500)

	if p.Status() != StatusRunning {
		t.Fatal("process is not running")
	}

	if started == p.Started() {
		t.Fatal("started time not changed")
	}

	pid2 := p.PID()

	if pid1 == pid2 {
		t.Fatal("pid is not changed")
	}

	if err := p.Terminate(); err != nil {
		t.Fatal(err)
	}
}

func TestCreateProcessFailure(t *testing.T) {
	p := CreateProcess(nil, nil, "mike", "bolsunovskyi")

	if err := p.Start(); err == nil {
		t.Fatal("no error on trying to start no existing binary")
	}
}

func TestCreateProcessFiles(t *testing.T) {
	stdErr, err := os.Create("err101.txt")
	if err != nil {
		t.Fatal(err)
	}
	stdOut, err := os.Create("out101.txt")
	if err != nil {
		t.Fatal(err)
	}

	p := CreateProcess(stdOut, stdErr, "uname", "-r")
	p.Restart = false

	if err := p.Start(); err != nil {
		t.Fatal(err)
	}

	if err := p.Wait(); err != nil {
		t.Fatal(err)
	}

	output, err := p.StdOut()
	if err != nil {
		t.Fatal("unable to get stdout")
	}
	t.Log(output)

	if output == "" {
		t.Fatal("empty stdout")
	}

	p = CreateProcess(stdOut, stdErr, "date", "-x")
	p.Restart = false

	if err := p.Start(); err != nil {
		t.Fatal(err)
	}

	if err := p.Wait(); err == nil {
		t.Fatal("no error on exit status 1")
	}

	output, err = p.StdErr()
	if err != nil {
		t.Fatal("unable to get stderr")
	}
	t.Log(output)
	if output == "" {
		t.Fatal("empty stderr")
	}

	if err := stdErr.Close(); err != nil {
		t.Fatal(err)
	}

	if err := stdOut.Close(); err != nil {
		t.Fatal(err)
	}

	if err := os.Remove("err101.txt"); err != nil {
		t.Fatal(err)
	}

	if err := os.Remove("out101.txt"); err != nil {
		t.Fatal(err)
	}
}
