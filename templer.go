package main

import (
	"os"
	"syscall"
	"unsafe"

	"./lib/login"
)

func main() {
	original_term, _ := login.Tcgetattr(syscall.Stdin)
	original_winsize := login.Winsize{}
	_ = login.Ioctl(syscall.Stdin, syscall.TIOCGWINSZ, uintptr(unsafe.Pointer(&original_winsize)))

	pid, master_fd, _ := login.Forkpty(*original_term, original_winsize)

	switch pid {
	case 0:
		syscall.Exec("/bin/bash", []string{""}, os.Environ())
	default:
		new_term, _ := login.Tcgetattr(syscall.Stdin)
		login.Cfmakeraw(new_term)
		login.Tcsetattr(syscall.Stdin, new_term)
		go func() {
			var buf = make([]byte, 256)

			for {
				nin, _ := syscall.Read(syscall.Stdin, buf)
				if nin <= 0 {
					break
				}
				syscall.Write(master_fd, buf[:nin])
			}
		}()

		var buf = make([]byte, 256)

		for {
			nin, _ := syscall.Read(master_fd, buf)
			if nin <= 0 {
				break
			}
			syscall.Write(syscall.Stdout, buf[:nin])
		}
	}

	defer login.Tcsetattr(syscall.Stdin, original_term)

}
