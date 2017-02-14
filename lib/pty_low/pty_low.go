package pty_low

import (
	"syscall"
	"unsafe"
)

func Openpty(term syscall.Termios, win Winsize) (master_fd int, slave_fd int, err error) {
	master_fd, err = openpt()
	if err != nil {
		return -1, -1, err
	}
	err = grantpt(master_fd)
	if err != nil {
		return -1, -1, err
	}
	err = unlockpt(master_fd)
	if err != nil {
		return -1, -1, err
	}

	slave_name, err := ptsname(master_fd)
	if err != nil {
		return -1, -1, err
	}
	slave_fd, err = syscall.Open(slave_name, syscall.O_RDWR|syscall.O_NOCTTY, 0)
	if err != nil {
		return -1, -1, err
	}

	Tcsetattr(slave_fd, &term)

	Ioctl(slave_fd, syscall.TIOCSWINSZ, uintptr(unsafe.Pointer(&win)))
	return master_fd, slave_fd, nil
}

func Forkpty(term syscall.Termios, win Winsize) (pid int, master_fd int, err error) {
	master_fd, slave_fd, err := Openpty(term, win)
	if err != nil {
		return -1, -1, err
	}

	pid = fork()

	switch pid {
	case -1:
		syscall.Close(master_fd)
		syscall.Close(slave_fd)
		return -1, -1, nil
	case 0:
		// child
		syscall.Close(master_fd)
		Login_tty(slave_fd)
		return 0, -1, nil
	default:
		// parent
		syscall.Close(slave_fd)
		return pid, master_fd, nil
	}

}

func Login_tty(fd int) {
	syscall.Setsid()
	Ioctl(fd, syscall.TIOCSCTTY, 0)

	syscall.Dup2(fd, syscall.Stdin)
	syscall.Dup2(fd, syscall.Stdout)
	syscall.Dup2(fd, syscall.Stderr)

	if fd > syscall.Stderr {
		syscall.Close(fd)
	}
}

func openpt() (fd int, err error) {
	fd, err = syscall.Open("/dev/ptmx", syscall.O_RDWR, 0)
	if err != nil {
		return -1, err
	}
	return fd, nil
}

func Ioctl(fd int, command uint, ptr uintptr) (err error) {
	_, _, err = syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(command), ptr)
	if err != nil && err.Error() != "errno 0" {
		return err
	}
	return nil
}

func Cfmakeraw(termios *syscall.Termios) {
	termios.Iflag &^= (syscall.IGNBRK | syscall.BRKINT | syscall.PARMRK | syscall.ISTRIP | syscall.INLCR | syscall.IGNCR | syscall.ICRNL | syscall.IXON)
	termios.Oflag &^= syscall.OPOST
	termios.Lflag &^= (syscall.ECHO | syscall.ECHONL | syscall.ICANON | syscall.ISIG | syscall.IEXTEN)
	termios.Cflag &^= (syscall.CSIZE | syscall.PARENB)
	termios.Cflag |= syscall.CS8
	termios.Cc[syscall.VMIN] = 1
	termios.Cc[syscall.VTIME] = 0
}

type Winsize struct {
	// unsigned short int
	Height uint16
	Width  uint16
	x      uint16
	y      uint16
}
