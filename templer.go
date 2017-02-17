package main

import (
	"syscall"

	"./lib/pty"
	"./lib/shell"
)

func main() {
	original_term, _ := pty.GetTermios(syscall.Stdin)
	winsize, _ := pty.GetWinsize(syscall.Stdin)

	s, _ := shell.Spawn(original_term, winsize)

	go func() {
		var buf = make([]byte, 256)

		for {
			nin, _ := syscall.Read(syscall.Stdin, buf)
			if nin < 0 {
				break
			}
			if nin == 0 {
			  continue
			}
			s.Write(buf[:nin])
		}
	}()

	var buf = make([]byte, 256)

	for {
		data, err := s.Read(&buf)
		if err != nil {
			break
		}
		syscall.Write(syscall.Stdout, data.Data)
	}

	defer original_term.SetTermios(syscall.Stdin)

}
