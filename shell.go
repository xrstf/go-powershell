// Copyright (c) 2017 Gorillalabs. All rights reserved.
// Copyright (c) 2020 xrstf.

package powershell

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"

	"go.xrstf.de/go-powershell/backend"
	"go.xrstf.de/go-powershell/utils"
)

const newline = "\r\n"

var (
	ShellClosedErr = errors.New("shell is closed")
)

type Shell interface {
	Execute(cmd string) (string, string, error)
	Exit()
}

type shell struct {
	handle backend.Waiter
	stdin  io.Writer
	stdout io.Reader
	stderr io.Reader
}

func New(backend backend.Starter) (Shell, error) {
	handle, stdin, stdout, stderr, err := backend.StartProcess("powershell.exe", "-NoExit", "-Command", "-")
	if err != nil {
		return nil, err
	}

	return &shell{handle, stdin, stdout, stderr}, nil
}

func (s *shell) Execute(cmd string) (string, string, error) {
	if s.handle == nil {
		return "", "", ShellClosedErr
	}

	outBoundary := createBoundary()
	errBoundary := createBoundary()

	// wrap the command in special markers so we know when to stop reading from the pipes
	full := fmt.Sprintf("%s; echo '%s'; [Console]::Error.WriteLine('%s')%s", cmd, outBoundary, errBoundary, newline)

	_, err := s.stdin.Write([]byte(full))
	if err != nil {
		return "", "", fmt.Errorf("failed to send PowerShell command: %w", err)
	}

	// read stdout and stderr
	sout := ""
	serr := ""

	waiter := &sync.WaitGroup{}
	waiter.Add(2)

	go streamReader(s.stdout, outBoundary, &sout, waiter)
	go streamReader(s.stderr, errBoundary, &serr, waiter)

	waiter.Wait()

	if len(serr) > 0 {
		return sout, serr, errors.New(serr)
	}

	return sout, serr, nil
}

func (s *shell) Exit() {
	_, _ = s.stdin.Write([]byte("exit" + newline))

	// if it's possible to close stdin, do so (some backends, like the local one,
	// do support it)
	closer, ok := s.stdin.(io.Closer)
	if ok {
		closer.Close()
	}

	_ = s.handle.Wait()

	s.handle = nil
	s.stdin = nil
	s.stdout = nil
	s.stderr = nil
}

func streamReader(stream io.Reader, boundary string, buffer *string, signal *sync.WaitGroup) {
	// read all output until we have found our boundary token
	output := ""
	bufsize := 64
	marker := boundary + newline

	for {
		buf := make([]byte, bufsize)
		read, err := stream.Read(buf)
		if err != nil {
			return
		}

		output = output + string(buf[:read])

		if strings.HasSuffix(output, marker) {
			break
		}
	}

	*buffer = strings.TrimSuffix(output, marker)
	signal.Done()
}

func createBoundary() string {
	return "$gorilla" + utils.CreateRandomString(12) + "$"
}
