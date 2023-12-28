package build

import (
	"crypto/rand"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/hugobyte/dive-core/cli/common"
	"github.com/mholt/archiver"
)

const (
	filesArtifactPermission = 0o744
	defaultTmpDir           = ""
	tmpDirPattern           = "tmp-dir-for-download-*"
	filesArtifactExtension  = ".tgz"
)

func RunBuild(cli *common.Cli) error {

	enclaveContext, err := cli.Context().GetEnclaveContext(common.EnclaveName)
	if err != nil {
		return common.WrapMessageToError(err, "Foundry build service failed while getting Enclave Context")
	}

	hash, err := generateRandomHash()

	if err != nil {
		return common.WrapMessageToError(err, "Error while generating the random hash for contract artifact name.")
	}

	contractArtifactName := fmt.Sprint("contract", "_", hash)

	_, _, err = enclaveContext.UploadFiles(contractPath, contractArtifactName)

	if err != nil {
		return common.WrapMessageToError(err, "Error while uploading the contract folder")
	}

	if err != nil {
		return common.WrapMessageToError(common.ErrInvalidEnclaveContext, err.Error())
	}

	param := fmt.Sprintf(`{"contract_artifact": %s}`, contractArtifactName)
	runConfig := common.GetStarlarkRunConfig(param, common.DiveFoundryScript, "run_build")

	response, _, err := enclaveContext.RunStarlarkRemotePackage(cli.Context().GetContext(), common.DiveRemotePackagePath, runConfig)

	if err != nil {
		return common.WrapMessageToErrorf(common.ErrStarlarkRunFailed, "%s. %s", err, "Foundry build failed")
	}

	responseData, services, _, err := common.GetSerializedData(cli, response)

	if err != nil {
		errRemove := cli.Context().RemoveServicesByServiceNames(services, common.EnclaveName)
		if err != nil {
			return common.WrapMessageToError(errRemove, "Foundry build service failed")
		}

		return common.WrapMessageToError(err, "Foundry build service failed")

	}

	//downlaod contract file artifact in current directory
	responseData = strings.Trim(responseData, "\"")
	contracts, err := enclaveContext.DownloadFilesArtifact(cli.Context().GetContext(), responseData)

	if err != nil {
		errRemove := cli.Context().RemoveServicesByServiceNames(services, common.EnclaveName)
		if err != nil {
			return common.WrapMessageToError(errRemove, "Foundry build service failed")
		}
		return common.WrapMessageToError(err, "Error while downloading the contract artifact")
	}

	tmpDirPath, err := os.MkdirTemp(defaultTmpDir, tmpDirPattern)

	if err != nil {
		return common.WrapMessageToError(err, "An error occurred while creating a temporary directory to download the files artifact with identifier")
	}

	fileNameToWriteTo := fmt.Sprintf("%v%v", responseData, filesArtifactExtension)

	shouldCleanupTmpDir := false
	defer func() {
		if shouldCleanupTmpDir {
			os.RemoveAll(tmpDirPath)
		}
	}()

	tmpFileToWriteTo := path.Join(tmpDirPath, fileNameToWriteTo)
	err = os.WriteFile(tmpFileToWriteTo, contracts, filesArtifactPermission)

	if err != nil {
		return common.WrapMessageToError(err, "An error occurred while writing bytes to file")
	}
	err = archiver.Unarchive(tmpFileToWriteTo, ".")
	if err != nil {
		return common.WrapMessageToError(err, "An error occurred while extracting")

	}
	//remove service after building contracts
	cli.Context().RemoveServicesByServiceNames(services, common.EnclaveName)

	return nil
}

// generateRandomHash(), This function will generate random hash of leanth 6
func generateRandomHash() (string, error) {
	b := make([]byte, 3) // 3 bytes = 6 hexadecimal digits
	_, err := rand.Read(b)
	if err != nil {
		return "", common.WrapMessageToError(err, "An error occurred while extracting")
	}
	return fmt.Sprintf("%x", b), nil
}
