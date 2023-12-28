package foundry

import (
	"os"
	"slices"	
	"github.com/hugobyte/dive-core/cli/cmd/foundry/build"
	"github.com/hugobyte/dive-core/cli/common"
	"github.com/spf13/cobra"
)

var FondryCmd = common.NewDiveCommandBuilder().
	SetUse("foundry").
	SetShort("Build solidity contracts").
	AddCommand(build.BuildCmd).
	SetRun(chains).
	Build()

func chains(cmd *cobra.Command, args []string) {

	validArgs := cmd.ValidArgs
	for _, c := range cmd.Commands() {
		validArgs = append(validArgs, c.Name())
	}
	cmd.ValidArgs = validArgs

	if len(args) == 0 {
		cmd.Help()

	} else if !slices.Contains(cmd.ValidArgs, args[0]) {

		cmd.Usage()
		os.Exit(1)
	}
}
