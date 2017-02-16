package pty

import (
	"syscall"
	"testing"
)

func TestWinsize(t *testing.T) {
	winsize, err := GetWinsize(syscall.Stdout)
	if err != nil {
		t.Error(err)
	}

	err = winsize.SetWinsize(syscall.Stdout)
	if err != nil {
		t.Error(err)
	}

	t.Logf("%#v", winsize)
}

func TestTermios(t *testing.T) {
	termios, err := GetTermios(syscall.Stdout)
	if err != nil {
		t.Error(err)
	}

	err = termios.SetTermios(syscall.Stdout)
	if err != nil {
		t.Error(err)
	}

	t.Logf("%#v", termios)
}

func TestForkpty(t *testing.T) {
	winsize, _ := GetWinsize(syscall.Stdout)
	termios, _ := GetTermios(syscall.Stdout)
	pid, _, err := Forkpty(termios, winsize)
	if err != nil {
		t.Error(err)
	}

	switch pid {
	case 0:
		syscall.Exit(0)
	default:
		t.Logf("Parent %v Child %v", syscall.Getpid(), pid)
	}
}
