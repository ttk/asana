package commands

import (
	"fmt"

	"github.com/ttk/asana/api"
	"github.com/urfave/cli/v2"
)

func Workspaces(c *cli.Context) {
	for _, w := range api.Me().Workspaces {
		fmt.Printf("%s %s\n", w.Gid, w.Name)
	}
}
