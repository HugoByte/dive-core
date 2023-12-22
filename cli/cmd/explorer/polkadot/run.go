package polkadot

import (
	"fmt"

	"github.com/hugobyte/dive-core/cli/common"
)

func RunPolkadotJs(cli *common.Cli) (*common.DiveMultipleServiceResponse, error) {
	enclaveContext, err := cli.Context().GetEnclaveContext(common.EnclaveName)

	if err != nil {
		return nil, common.WrapMessageToError(err, "Polkadot js Run Failed")
	}

	para := fmt.Sprintf(`{"args": %s}`)

	runConfig := common.GetStarlarkRunConfig(para, common.DivePolkadotExplorerPath, runPlkadotjsFunction)
	response, _, err := enclaveContext.RunStarlarkPackage(cli.Context().GetContext(), common.PolkadotRemotePackagePath, runConfig)
	if err != nil {
		return nil, common.WrapMessageToError(common.ErrStarlarkRunFailed, err.Error())
	}

	responseData, services, skippedInstructions, err := common.GetSerializedData(cli, response)

	if err != nil {

		errRemove := cli.Context().RemoveServicesByServiceNames(services, common.DiveEnclave)
		if errRemove != nil {
			return nil, common.WrapMessageToError(errRemove, "Polkadot js Run Failed ")
		}

		return nil, common.WrapMessageToError(err, "Polkadot js Run Failed ")
	}

	if cli.Context().CheckSkippedInstructions(skippedInstructions) {
		return nil, common.WrapMessageToError(common.ErrStarlarkResponse, "Polkadot js is already Running")
	}

	polkadotResponseData := &common.DiveMultipleServiceResponse{}
	result, err := polkadotResponseData.Decode([]byte(responseData))

	if err != nil {

		errRemove := cli.Context().RemoveServicesByServiceNames(services, common.DiveEnclave)
		if errRemove != nil {
			return nil, common.WrapMessageToError(errRemove, "Polkadot js Run Failed ")
		}

		return nil, common.WrapMessageToErrorf(common.ErrDataUnMarshall, "%s.%s", err, "Polkadot js Run Failed ")

	}

	return result, nil

}
