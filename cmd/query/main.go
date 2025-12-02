// DNS query CLI entrypoint - delegates to cli.NewQueryCommand.
package main

import (
	"fmt"
	"os"

	"github.com/sudo-Tiz/DNS-Tester-GO/internal/cli"
)

func main() {
	cmd := cli.NewQueryCommand()
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
