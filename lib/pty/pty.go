package pty

import (
	"syscall"
	"unsafe"

	"../pty_low"
)

type Termios syscall.Termios

func GetTermios(fd int) (termios Termios, err error) {
	ptr, err := pty_low.Tcgetattr(fd)
	return Termios(*ptr), err
}

func (self Termios) SetTermios(fd int) (err error) {
	syster := syscall.Termios(self)
	return pty_low.Tcsetattr(fd, &syster)
}

func (self Termios) Rawmode() Termios {
	var rawterm = syscall.Termios(self)
	pty_low.Cfmakeraw(&rawterm)
	return Termios(rawterm)
}

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
