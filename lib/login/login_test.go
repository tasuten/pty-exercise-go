package login

import (
	"syscall"
	"testing"
)

func TestOpenpt(t *testing.T) {
	fd, err := openpt()

	if err != nil {
		t.Error(err)
	}
	if fd < 3 {
		t.Error("fd must be greater than 2")
	}
}

func TestPtsname(t *testing.T) {
	master_fd, err := openpt()
	err = grantpt(master_fd)
	if err != nil {
		t.Error(err)
	}
	err = unlockpt(master_fd)
	if err != nil {
		t.Error(err)
	}

	_, err = ptsname(master_fd)
	if err != nil {
		t.Error(err)
	}
}

func TestTermios(t *testing.T) {
	master_fd, _ := openpt()
	_ = grantpt(master_fd)
	_ = unlockpt(master_fd)
	slave_name, _ := ptsname(master_fd)
	slave_fd, err := syscall.Open(slave_name, syscall.O_RDWR|syscall.O_NOCTTY, 0)
	if err != nil {
		t.Error(err)
	}
	termios1, err := Tcgetattr(slave_fd)
	if err != nil {
		t.Error(err)
	}
	err = Tcsetattr(slave_fd, termios1)
	if err != nil {
		t.Error(err)
	}
	termios2, err := Tcgetattr(slave_fd)
	if err != nil {
		t.Error(err)
	}
	if *termios1 != *termios2 {
		t.Error("tcset/getattr is not symmetry")
	}

}

func TestFork(t *testing.T) {
	pid := fork()

	switch pid {
	case -1:
		t.Error("fork failed")
	case 0:
		// child
		syscall.Exit(0)
	default:
		// parent
	}

}
