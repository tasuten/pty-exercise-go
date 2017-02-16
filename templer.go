package main

import (
	"os"
	"syscall"

	"./lib/pty"
)

func main() {
	original_term, _ := pty.GetTermios(syscall.Stdin)
	winsize, _ := pty.GetWinsize(syscall.Stdin)

	pid, master_fd, _ := pty.Forkpty(original_term, winsize)

	switch pid {
	case 0:
		syscall.Exec("/bin/bash", []string{""}, os.Environ())
	default:
		var newterm = original_term.Rawmode()
		newterm.Cc[syscall.VTIME] = 0
		newterm.Cc[syscall.VMIN] = 1

		newterm.SetTermios(syscall.Stdin)
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

	defer original_term.SetTermios(syscall.Stdin)

}
