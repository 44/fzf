package fzf

import (
	"github.com/junegunn/fzf/src/util"
	"os/exec"
	"io"
	"bufio"
)

type Formatter struct {
	cmd      *exec.Cmd
	stdin    io.WriteCloser
	stdout   io.ReadCloser
	stderr   io.ReadCloser
	scanner *bufio.Scanner
}

func NewFormatter(command string) *Formatter {
	cmd := util.ExecCommand(command, true)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil
	}

	err = cmd.Start()
	if err != nil {
		return nil
	}

	return &Formatter{cmd, stdin, stdout, stderr, nil}
}

func (f *Formatter) Format(data []byte) []byte {
	if f.scanner == nil {
		f.scanner = bufio.NewScanner(f.stdout)
	}
	f.stdin.Write(data)
	f.stdin.Write([]byte("\n"))

	if f.scanner.Scan() {
		return f.scanner.Bytes()
	} else {
		if f.scanner.Err() != nil {
			return append([]byte(f.scanner.Err().Error()), data...)
		}
	}

	return data
}
