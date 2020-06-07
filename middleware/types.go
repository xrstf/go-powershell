// Copyright (c) 2017 Gorillalabs. All rights reserved.
// Copyright (c) 2020 xrstf.

package middleware

type Middleware interface {
	Execute(cmd string) (string, string, error)
	Exit()
}
