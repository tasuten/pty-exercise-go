package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"./lib/pty"
	"./lib/shell"
)

func main() {
	original_term, _ := pty.GetTermios(syscall.Stdin)
	winsize, _ := pty.GetWinsize(syscall.Stdin)

	s, _ := shell.Spawn(original_term, winsize)

	wg := sync.WaitGroup{}
	wg.Add(2)
	go input_handler(&s, &wg)
	go output_handler(&s, &wg)
	wg.Wait()

	defer original_term.SetTermios(syscall.Stdin)

}

func input_handler(s *shell.Shell, wg *sync.WaitGroup) {
	defer wg.Done()
	var buf = make([]byte, 256)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGCHLD)

	for {
		n_read, _ := syscall.Read(syscall.Stdin, buf)
		if n_read < 0 {
			break
		}

		if n_read == 0 {
			select {
			case <-c:
				return
			default:
				continue
			}

		}

		s.Write(buf[:n_read])
	}

}

func output_handler(s *shell.Shell, wg *sync.WaitGroup) {
	defer wg.Done()
	var buf = make([]byte, 256)
	for {
		data, err := s.Read(&buf)
		if err != nil {
			break
		}
		syscall.Write(syscall.Stdout, data.Data)
	}
}
