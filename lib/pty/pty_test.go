package pty

import (
	"syscall"
	"testing"
)

func TestWinsize(t *testing.T) {
	winsize, err := GetWinsize(syscall.Stdout)
	if err != nil {
		t.Errorf("Fail to get winsize: %v", err)
	}

	err = winsize.SetWinsize(syscall.Stdout)
	if err != nil {
	  t.Errorf("Fail to set winsize: %v", err)
	}

	t.Log(winsize)
}
