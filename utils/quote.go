// Copyright (c) 2017 Gorillalabs. All rights reserved.
// Copyright (c) 2020 xrstf.

package utils

import "strings"

func QuoteArg(s string) string {
	return "'" + strings.Replace(s, "'", "\"", -1) + "'"
}
