package pty

import (
	"syscall"
	"unsafe"

	"../pty_low"
)

type Termios syscall.Termios

type Winsize pty_low.Winsize

func GetWinsize(fd int) (winsize Winsize, err error) {
	winsize = Winsize{}
	err = pty_low.Ioctl(fd, syscall.TIOCGWINSZ, uintptr(unsafe.Pointer(&winsize)))
	return winsize, err
}

func (self Winsize) SetWinsize(fd int) (err error) {
	err = pty_low.Ioctl(fd, syscall.TIOCSWINSZ, uintptr(unsafe.Pointer(&self)))
	return err
}
