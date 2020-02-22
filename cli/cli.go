package cli

import (
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "space",
		Usage: "",
		Action: func(c *cli.Context) error {
		},
	}
	app.Run(os.Args)
}
