// Copyright (c) 2017 Gorillalabs. All rights reserved.
// Copyright (c) 2020 xrstf.

package backend

import "io"

type Waiter interface {
	Wait() error
}

type Starter interface {
	StartProcess(cmd string, args ...string) (Waiter, io.Writer, io.Reader, io.Reader, error)
}
