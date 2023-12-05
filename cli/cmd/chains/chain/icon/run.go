package icon

import (
	"fmt"

	"github.com/hugobyte/dive-core/cli/common"
)

func RunIconNode(cli *common.Cli) (*common.DiveServiceResponse, error) {

	enclaveContext, err := cli.Context().GetEnclaveContext(common.DiveEnclave)

	if err != nil {
		return nil, err
	}
	var serviceConfig = &IconServiceConfig{}
	err = LoadConfig(cli, serviceConfig, configFilePath)
	if err != nil {
		return nil, err
	}
	cli.Spinner().StartWithMessage("Starting Icon Node", "green")

	genesisHandler, err := genesismanager(enclaveContext)
	if err != nil {
		return nil, common.Errorc(common.FileError, err.Error())
	}
	params := fmt.Sprintf(`{"private_port":%d, "public_port":%d, "p2p_listen_address": %s, "p2p_address":%s, "cid": "%s","uploaded_genesis":%s,"genesis_file_path":"%s","genesis_file_name":"%s"}`, serviceConfig.Port, serviceConfig.PublicPort, serviceConfig.P2PListenAddress, serviceConfig.P2PAddress, serviceConfig.Cid, genesisHandler.uploadedFiles, genesisHandler.genesisPath, genesisHandler.genesisFile)
	starlarkConfig := common.GetStarlarkRunConfig(params, common.DiveIconNodeScript, "start_icon_node")

	iconData, _, err := enclaveContext.RunStarlarkRemotePackage(cli.Context().GetContext(), common.DiveRemotePackagePath, starlarkConfig)

	if err != nil {
		return nil, common.Errorc(common.FileError, err.Error())
	}

	response, services, skippedInstructions, err := common.GetSerializedData(cli, iconData)

	if err != nil {

		err = cli.Context().RemoveServicesByServiceNames(services, common.DiveEnclave)
		if err != nil {
			return nil, common.Errorc(common.InvalidEnclaveContextError, err.Error())
		}

		return nil, common.Errorc(common.KurtosisContextError, err.Error())
	}

	if cli.Context().CheckSkippedInstructions(skippedInstructions) {
		return nil, common.Errorc(common.KurtosisContextError, "Already Running")
	}

	iconResponseData := &common.DiveServiceResponse{}

	result, err := iconResponseData.Decode([]byte(response))

	if err != nil {

		return nil, common.Errorc(common.KurtosisContextError, err.Error())
	}

	return result, nil
}

func RunDecentralization(cli *common.Cli, params string) error {

	cli.Spinner().SetSuffixMessage("Starting Icon Node Decentralization", "green")

	kurtosisEnclaveContext, err := cli.Context().GetEnclaveContext(common.DiveEnclave)

	if err != nil {
		return common.Errorc(common.KurtosisContextError, err.Error())
	}
	starlarkConfig := common.GetStarlarkRunConfig(params, common.DiveIconDecentraliseScript, "configure_node")
	data, _, err := kurtosisEnclaveContext.RunStarlarkRemotePackage(cli.Context().GetContext(), common.DiveRemotePackagePath, starlarkConfig)

	if err != nil {
		return common.Errorc(common.KurtosisContextError, err.Error())
	}

	_, services, skippedInstructions, err := common.GetSerializedData(cli, data)
	if err != nil {

		err = cli.Context().RemoveServicesByServiceNames(services, common.DiveEnclave)
		if err != nil {
			return common.Errorc(common.InvalidEnclaveContextError, err.Error())
		}

		return common.Errorc(common.KurtosisContextError, err.Error())
	}
	if cli.Context().CheckSkippedInstructions(skippedInstructions) {
		return common.Errorc(common.KurtosisContextError, "Decentralization Already  Completed ")
	}

	return nil

}