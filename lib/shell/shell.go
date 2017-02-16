package shell

import (
	"os"
	"syscall"

	"../pty"
)

type Shell struct {
	Master_fd int
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
		new_t.Cc[syscall.VMIN] = 1
		new_t.SetTermios(syscall.Stdin)
		result.Master_fd = master_fd
	}

	return result, nil
}
