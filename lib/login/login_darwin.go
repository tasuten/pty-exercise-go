// +build darwin
package login

import (
	"errors"
	"syscall"
	"unsafe"
)

func unlockpt(fd int) (err error) {
	return ioctl(fd, syscall.TIOCPTYUNLK, 0)
}

func grantpt(fd int) (err error) {
	return ioctl(fd, syscall.TIOCPTYGRANT, 0)
}

func ptsname(fd int) (name string, err error) {
	n := make([]byte, 128) // from apple libc
	err = ioctl(fd, syscall.TIOCPTYGNAME, uintptr(unsafe.Pointer(&n[0])))
	if err != nil {
		return "", err
	}
	for i, char := range n {
		if char == 0 { // NULL char
			return string(n[:i]), nil
		}
	}

	return "", errors.New("TIOCPTYGNAME string doesn't NULL-terminated")
}

func tcsetattr(fd int, termios *syscall.Termios) (err error) {
	return ioctl(fd, syscall.TIOCSETAF, uintptr(unsafe.Pointer(termios)))
}

func tcgetattr(fd int) (*syscall.Termios, error) {
	var termios = &syscall.Termios{}
	err := ioctl(fd, syscall.TIOCGETA, uintptr(unsafe.Pointer(termios)))
	if err != nil {
		return termios, err
	}
	return termios, nil
}

func fork() int {
	var r1, r2 uintptr
	var e syscall.Errno

	r1, r2, e = syscall.RawSyscall(syscall.SYS_FORK, 0, 0, 0)

	if e != 0 || r2 < 0 {
		return -1
	}

	var pid = int(r1) // ?
	if r2 == 1 {
		pid = 0
	}
	return pid
}
