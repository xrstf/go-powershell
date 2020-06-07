// Copyright (c) 2017 Gorillalabs. All rights reserved.
// Copyright (c) 2020 xrstf.

package backend

import (
	"fmt"
	"io"
	"regexp"
	"strings"
)

// sshSession exists so we don't create a hard dependency on crypto/ssh.
type sshSession interface {
	Waiter

	StdinPipe() (io.WriteCloser, error)
	StdoutPipe() (io.Reader, error)
	StderrPipe() (io.Reader, error)
	Start(string) error
}

type SSH struct {
	Session sshSession
}

func (b *SSH) StartProcess(cmd string, args ...string) (Waiter, io.Writer, io.Reader, io.Reader, error) {
	stdin, err := b.Session.StdinPipe()
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to get SSH session's stdin stream: %w", err)
	}

	stdout, err := b.Session.StdoutPipe()
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to get SSH session's stdout stream: %w", err)
	}

	stderr, err := b.Session.StderrPipe()
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to get SSH session's stderr stream: %w", err)
	}

	err = b.Session.Start(b.createCmd(cmd, args))
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to spawn process via SSH: %w", err)
	}

	return b.Session, stdin, stdout, stderr, nil
}

func (b *SSH) createCmd(cmd string, args []string) string {
	parts := []string{cmd}
	simple := regexp.MustCompile(`^[a-z0-9_/.~+-]+$`)

	for _, arg := range args {
		if !simple.MatchString(arg) {
			arg = b.quote(arg)
		}

		parts = append(parts, arg)
	}

	return strings.Join(parts, " ")
}

func (b *SSH) quote(s string) string {
	return fmt.Sprintf(`"%s"`, s)
}
