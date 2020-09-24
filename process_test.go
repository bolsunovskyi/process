package process

import (
	"testing"
	"time"
)

func TestCreateProcessNoRestart(t *testing.T) {
	p, err := CreateProcess("uname", "-r")
	if err != nil {
		t.Fatal(err)
	}
	p.Restart = false

	if err := p.Start(); err != nil {
		t.Fatal(err)
	}

	if err := p.Wait(); err != nil {
		t.Fatal(err)
	}

	out := p.StdOut()
	t.Log(out)
	if out == "" {
		t.Fatal("empty stdout")
	}
}

func TestCreateProcessRestart(t *testing.T) {
	p, err := CreateProcess("uname", "-r")
	if err != nil {
		t.Fatal(err)
	}

	if err := p.Start(); err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second)

	if err := p.Kill(); err != nil {
		t.Fatal(err)
	}

	out := p.StdOut()
	t.Log(out)
	if out == "" {
		t.Fatal("empty stdout")
	}
}

