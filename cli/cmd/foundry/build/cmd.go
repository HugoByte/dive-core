package build

import (
	"fmt"

	"github.com/hugobyte/dive-core/cli/common"
	"github.com/spf13/cobra"
)

var (
	contractPath string
)

var	BuildCmd = common.NewDiveCommandBuilder().
	SetUse("build").
	SetShort("Build solidity contract").
	AddStringFlagWithShortHand(&contractPath, "path", "p", "", "root project path of contracts").
	SetRun(build).
	Build()


func build(cmd *cobra.Command, args []string) {

	cliContext := common.GetCliWithKurtosisContext()

	err := common.ValidateArgs(args)

	if err != nil {
		cliContext.Fatalf("Error %s. %s", err, cmd.UsageString())
	}

	cliContext.Spinner().StartWithMessage("Building contracts", "green")

	err = RunBuild(cliContext)

	if err != nil {
		cliContext.Fatal(err)
	}

	stopMessage := fmt.Sprintf("Contract]. Please find the artifact in current working directory")
	cliContext.Spinner().StopWithMessage(stopMessage)

}
             