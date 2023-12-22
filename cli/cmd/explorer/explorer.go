package explorer

import (
	"os"
	"slices"

	"github.com/hugobyte/dive-core/cli/cmd/explorer/polkadot"
	"github.com/hugobyte/dive-core/cli/common"
	"github.com/spf13/cobra"
)

var ExplorerCmd = common.NewDiveCommandBuilder().
	SetUse("explorer").
	SetShort("Command for to run expolores").
	AddCommand(polkadot.PolkadotJsCMD).
	SetRun(explorer).
	Build()

func explorer(cmd *cobra.Command, args []string) {
	cli := common.GetCli()
	validArgs := cmd.ValidArgs
	for _, c := range cmd.Commands() {
		validArgs = append(validArgs, c.Name())
	}
	cmd.ValidArgs = validArgs

	if len(args) == 0 {
		cmd.Help()

	} else if !slices.Contains(cmd.ValidArgs, args[0]) {
		cli.Error(common.WrapMessageToErrorf(common.ErrInvalidCommand, "%s", cmd.UsageString()))
		os.Exit(1)
	}
}
