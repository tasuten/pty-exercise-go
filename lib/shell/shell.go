package shell

import (
	"errors"
	"os"
	"syscall"

	"../pty"
)

type Shell struct {
	Master_fd int
}

type ReadData struct {
	Data []byte
	Fd   int
}

func Spawn(t pty.Termios, w pty.Winsize) (Shell, error) {
	var result = Shell{}
	pid, master_fd, err := pty.Forkpty(t, w)
	if err != nil {
		return result, err
	}

	switch pid {
	case 0:
		syscall.Exec("/bin/bash", []string{""}, os.Environ())
		syscall.Exit(0)
	default:
		// 親プロセスでのみ実行されるべき処理っぽい
		// うまくproxyとかに分離できないか
		var new_t = t.Rawmode()
		new_t.Cc[syscall.VTIME] = 0
		new_t.Cc[syscall.VMIN] = 0
		new_t.SetTermios(syscall.Stdin)
		result.Master_fd = master_fd
	}

	return result, nil
}

func (self Shell) Read(buf *[]byte) (ReadData, error) {
	result := ReadData{}
	n_read, err := syscall.Read(self.Master_fd, *buf)
	if err != nil {
		return result, err
	}

	if n_read <= 0 {
		return result, errors.New("Read error")
	}

	result.Fd = self.Master_fd
	result.Data = (*buf)[:n_read]

	return result, nil
}

func (self Shell) Write(data []byte) (err error) {
	n_write, err := syscall.Write(self.Master_fd, data)

	if err != nil {
		return err
	}

	if n_write != len(data) {
		return errors.New("Write error")
	}

	return nil
}
