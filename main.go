// Command that implements the main executable.
package main

import (
	"github.com/distribworks/dkron/v2/cmd"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cmd.Execute()
}
