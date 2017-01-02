# go-powershell

This package is inspired by jPowerShell and allows one to run and remote-control a
PowerShell session. Use this if you don't have a static script that you want to
execute, bur rather run dynamic commands.

## Installation

    go get github.com/gorillalabs/go-powershell

## Usage

The package was originally written to use remote powershell sessions, so a few API
methods are geared towards that usecase.

```go
package main

import (
	"fmt"

	ps "github.com/gorillalabs/go-powershell"
)

func main() {
	config := ps.NewDefaultConfig()
	config.ComputerName = "remote-pc-1"

	session, err := ps.EnterSession(config)
	if err != nil {
		panic(err)
	}
	defer session.Exit()

	stdout, stderr, err := session.Execute("Get-WmiObject -Class Win32_Processor")
	if err != nil {
		panic(err)
	}

	fmt.Println(stdout)
}
```

## License

MIT, see LICENSE file.
