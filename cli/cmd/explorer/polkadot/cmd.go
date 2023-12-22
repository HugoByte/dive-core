package polkadot

import (
	"fmt"

	"github.com/hugobyte/dive-core/cli/common"
	"github.com/spf13/cobra"
)

const runPlkadotjsFunction = "run_pokadot_js_app"

var PolkadotJsCMD = common.NewDiveCommandBuilder().
	SetUse("polkadot").
	SetShort("Starts polkadot js explorer").
	SetLong("Starts polkadot js explorer").
	SetRun(runPolkadotJs).
	Build()

func runPolkadotJs(cmd *cobra.Command, args []string) {

	cliContext := common.GetCli()

	err := common.ValidateArgs(args)
	if err != nil {
		cliContext.Fatal(err)
	}
	cliContext.StartSpinnerIfNotVerbose("Starting polkadot js explore", common.DiveLogs)

	response, err := RunPolkadotJs(cliContext)

	if err != nil {
		cliContext.Fatal(err)
	}

	serviceFileName := fmt.Sprintf(common.ServiceFilePath, common.EnclaveName)

	fmt.Print(response.Dive)
	for serviceName := range response.Dive {
		err = common.WriteServiceResponseData(response.Dive[serviceName].ServiceName, *response.Dive[serviceName], cliContext, serviceFileName)

		if err != nil {
			cliContext.Fatal(err)
		}
	}
	stopMessage := fmt.Sprintf("Polkadot Node Started. Please find service details in current working directory(%s)", serviceFileName)
	cliContext.StopSpinnerIfNotVerbose(stopMessage, common.DiveLogs)

}
