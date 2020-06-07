// Copyright (c) 2017 Gorillalabs. All rights reserved.
// Copyright (c) 2020 xrstf.

package backend

import (
	"fmt"
	"io"
	"os/exec"
)

type Local struct{}

func (b *Local) StartProcess(cmd string, args ...string) (Waiter, io.Writer, io.Reader, io.Reader, error) {
	command := exec.Command(cmd, args...)

	stdin, err := command.StdinPipe()
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to get PowerShell's stdin stream: %w", err)
	}

	stdout, err := command.StdoutPipe()
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to get PowerShell's stdout stream: %w", err)
	}

	stderr, err := command.StderrPipe()
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to get PowerShell's stderr stream: %w", err)
	}

	err = command.Start()
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to spawn PowerShell process: %w", err)
	}

	return command, stdin, stdout, stderr, nil
}
