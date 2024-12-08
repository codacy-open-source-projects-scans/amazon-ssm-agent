// Copyright 2021 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may not
// use this file except in compliance with the License. A copy of the
// License is located at
//
// http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
// either express or implied. See the License for the specific language governing
// permissions and limitations under the License.

//go:build freebsd || linux || netbsd || openbsd
// +build freebsd linux netbsd openbsd

package main

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/aws/amazon-ssm-agent/agent/log"
	"github.com/aws/amazon-ssm-agent/agent/managedInstances/registration"
	rMock "github.com/aws/amazon-ssm-agent/agent/managedInstances/registration/mocks"
	logmocks "github.com/aws/amazon-ssm-agent/agent/mocks/log"
	"github.com/aws/amazon-ssm-agent/agent/setupcli/managers/common"
	"github.com/aws/amazon-ssm-agent/agent/setupcli/managers/configurationmanager"
	cmMock "github.com/aws/amazon-ssm-agent/agent/setupcli/managers/configurationmanager/mocks"
	"github.com/aws/amazon-ssm-agent/agent/setupcli/managers/downloadmanager"
	dmMock "github.com/aws/amazon-ssm-agent/agent/setupcli/managers/downloadmanager/mocks"
	"github.com/aws/amazon-ssm-agent/agent/setupcli/managers/packagemanagers"
	pmMock "github.com/aws/amazon-ssm-agent/agent/setupcli/managers/packagemanagers/mocks"
	"github.com/aws/amazon-ssm-agent/agent/setupcli/managers/registermanager"
	rmMock "github.com/aws/amazon-ssm-agent/agent/setupcli/managers/registermanager/mocks"
	"github.com/aws/amazon-ssm-agent/agent/setupcli/managers/servicemanagers"
	smMock "github.com/aws/amazon-ssm-agent/agent/setupcli/managers/servicemanagers/mocks"
	"github.com/aws/amazon-ssm-agent/agent/setupcli/managers/verificationmanagers"
	vmMock "github.com/aws/amazon-ssm-agent/agent/setupcli/managers/verificationmanagers/mocks"
	"github.com/aws/amazon-ssm-agent/agent/setupcli/utility"
	"github.com/aws/amazon-ssm-agent/agent/updateutil/updateinfo"
	agentVersioning "github.com/aws/amazon-ssm-agent/agent/version"
	"github.com/aws/amazon-ssm-agent/core/executor"
	"github.com/aws/amazon-ssm-agent/core/executor/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const breakOutWithPanicMessageOnprem = "BREAKOUT_WITH_PANIC"

func storeMockedFunctionsOnprem() func() {
	getPackageManagerStorage := getPackageManager
	getConfigurationManagerStorage := getConfigurationManager
	getServiceManagerStorage := getServiceManager
	getRegisterManagerStorage := getRegisterManager
	getRegistrationInfoStorage := getRegistrationInfo
	hasElevatedPermissions = func() error {
		return nil
	}
	return func() {
		getPackageManager = getPackageManagerStorage
		getConfigurationManager = getConfigurationManagerStorage
		getServiceManager = getServiceManagerStorage
		getRegisterManager = getRegisterManagerStorage
		getRegistrationInfo = getRegistrationInfoStorage
	}
}

func setArgsAndRestoreOnprem(args ...string) func() {
	var oldArgs = make([]string, len(os.Args))
	copy(oldArgs, os.Args)
	os.Args = args

	hasElevatedPermissions = func() error {
		return nil
	}
	return func() {
		os.Args = oldArgs
	}
}

func TestMain_ErrorGetPackageManager(t *testing.T) {
	initializeArgs()
	defer storeMockedFunctions()()

	defer setArgsAndRestore("/some/path/setupcli", "-shutdown", "-env", "greengrass")()

	getPackageManager = func(log.T) (packagemanagers.IPackageManager, error) {
		return nil, fmt.Errorf("SomeError")
	}

	osExit = func(exitCode int, log log.T, message string, args ...interface{}) {
		assert.Equal(t, 1, exitCode)
		assert.Contains(t, message, "Failed to determine package manager")

		panic(breakOutWithPanicMessage)
	}

	defer func() {
		if errInterface := recover(); errInterface != nil {
			assert.Equal(t, breakOutWithPanicMessage, errInterface)
		}
	}()
	main()
	assert.True(t, false, "Should never reach here because of exit")
}

func TestMain_OnPrem_GetExecutingDirectoryPath_Failed(t *testing.T) {
	evalSymLinks = func(path string) (string, error) {
		return "", nil
	}
	osExecutable = func() (string, error) {
		return "", fmt.Errorf("os executable")
	}
	filePathDir = func(path string) string {
		return "sample"
	}
	pkgManagerMock := &pmMock.IPackageManager{}
	svcManagerMock := &smMock.IServiceManager{}
	verificationManager := &vmMock.IVerificationManager{}
	logMock := logmocks.NewMockLog()
	err := performOnpremSteps(logMock, pkgManagerMock, verificationManager, svcManagerMock)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not get the ssm-setup-cli executable path")
}

func TestMain_OnPrem_GetExecutingDirectoryPath_Success(t *testing.T) {
	evalSymLinks = func(path string) (string, error) {
		return "test", nil
	}
	osExecutable = func() (string, error) {
		return "", fmt.Errorf("err1")
	}
	filePathDir = func(path string) string {
		return "sample"
	}
	pkgManagerMock := &pmMock.IPackageManager{}
	svcManagerMock := &smMock.IServiceManager{}
	verificationManager := &vmMock.IVerificationManager{}
	logMock := logmocks.NewMockLog()
	err := performOnpremSteps(logMock, pkgManagerMock, verificationManager, svcManagerMock)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not get the ssm-setup-cli executable path")
}

func TestMain_OnPrem_Register_Success(t *testing.T) {
	evalSymLinks = func(path string) (string, error) {
		return "test", nil
	}
	osExecutable = func() (string, error) {
		return "", nil
	}
	filePathDir = func(path string) string {
		return "sample"
	}
	signatureFile := "sign1"
	pkgManagerMock := &pmMock.IPackageManager{}
	pkgManagerMock.On("IsAgentInstalled").Return(false, nil).Once()
	pkgManagerMock.On("IsAgentInstalled").Return(true, nil).Once()
	pkgManagerMock.On("GetFileExtension").Return("rpm")
	pkgManagerMock.On("GetInstalledAgentVersion").Return("", nil)
	svcManagerMock := &smMock.IServiceManager{}
	verificationManager := &vmMock.IVerificationManager{}
	verificationManager.On("VerifySignature", mock.Anything, signatureFile, mock.Anything, mock.Anything).Return(nil).Once()
	logMock := logmocks.NewMockLog()

	tempDir := "temp1"
	fileUtilCreateTemp = func(dir, prefix string) (name string, err error) {
		return tempDir, nil
	}
	fileUtilMakeDirs = func(destinationDir string) (err error) {
		return nil
	}
	isPlatformNano = func(log log.T) (bool, error) {
		return false, nil
	}

	getDownloadManager = func(log log.T, region string, manifestUrl string, updateInfo updateinfo.T, setupCLIArtifactsPath string, isNano bool) downloadmanager.IDownloadManager {
		managerMock := &dmMock.IDownloadManager{}
		managerMock.On("DownloadLatestSSMSetupCLI", mock.Anything, mock.Anything).Return(nil).Once()
		managerMock.On("GetLatestVersion").Return(agentVersioning.Version, nil).Once()
		managerMock.On("DownloadArtifacts", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		managerMock.On("DownloadSignatureFile", mock.Anything, mock.Anything, mock.Anything).Return(signatureFile, nil).Once()
		return managerMock
	}
	getConfigurationManager = func() configurationmanager.IConfigurationManager {
		managerMock := &cmMock.IConfigurationManager{}
		managerMock.On("CreateUpdateAgentConfigWithOnPremIdentity").Return(nil)
		return managerMock
	}

	utilityCheckSum = func(filePath string) (hash string, err error) {
		return "", nil
	}
	version = utility.LatestVersionString
	helperInstallAgent = func(log log.T, pManager packagemanagers.IPackageManager, sManager servicemanagers.IServiceManager, folderPath string) error {
		return nil
	}

	newProcessExecutor = func(log log.T) executor.IExecutor {
		executorMock := &mocks.IExecutor{}
		executorMock.On("Processes").Return([]executor.OsProcess{executor.OsProcess{
			Executable: utility.AgentBinary,
		}}, nil)
		return executorMock
	}
	timeSleep = func(d time.Duration) {
		return
	}
	getRegistrationInfo = func() registration.IOnpremRegistrationInfo {
		registrationMock := &rMock.IOnpremRegistrationInfo{}
		registrationMock.On("InstanceID", mock.Anything, "", mock.Anything).Return("").Once()
		registrationMock.On("ReloadInstanceInfo", mock.Anything, "", mock.Anything).Return("")
		registrationMock.On("InstanceID", mock.Anything, "", mock.Anything).Return("i-temp")
		return registrationMock
	}
	getRegisterManager = func() registermanager.IRegisterManager {
		managerMock := &rmMock.IRegisterManager{}
		managerMock.On("RegisterAgent", mock.Anything).Return(nil)
		return managerMock
	}
	svcMgrStopAgent = func(manager servicemanagers.IServiceManager, log log.T) error {
		return nil
	}
	startAgent = func(manager servicemanagers.IServiceManager, log log.T) error {
		return nil
	}
	err := performOnpremSteps(logMock, pkgManagerMock, verificationManager, svcManagerMock)
	assert.Nil(t, err)
}

func TestMain_Register_InvalidActivationId_Failed(t *testing.T) {
	defer storeMockedFunctionsOnprem()()
	defer setArgsAndRestoreOnprem("/some/path/setupcli", "-env", "onprem", "-register", "-activation-id", "", "-activation-code", "test")()
	hasElevatedPermissions = func() error {
		return nil
	}
	getServiceManager = func(log.T) (servicemanagers.IServiceManager, error) {
		managerMock := &smMock.IServiceManager{}
		return managerMock, nil
	}

	osExit = func(exitCode int, log log.T, message string, args ...interface{}) {
		assert.Equal(t, 1, exitCode)
		message = fmt.Sprintf(message, args)
		assert.Contains(t, message, "Activation id required for on-prem registration.")
		panic(breakOutWithPanicMessageOnprem)
	}

	defer func() {
		if errInterface := recover(); errInterface != nil {
			fmt.Println(errInterface)
			assert.Equal(t, breakOutWithPanicMessageOnprem, errInterface)
		}
	}()
	main()
	assert.True(t, false, "Should never reach here because of exit")
}

func TestMain_Register_InvalidActivationCode_Failed(t *testing.T) {
	defer storeMockedFunctionsOnprem()()
	defer setArgsAndRestoreOnprem("/some/path/setupcli", "-env", "onprem", "-register", "-activation-id", "test", "-activation-code", "")()
	hasElevatedPermissions = func() error {
		return nil
	}
	getServiceManager = func(log.T) (servicemanagers.IServiceManager, error) {
		managerMock := &smMock.IServiceManager{}
		return managerMock, nil
	}

	osExit = func(exitCode int, log log.T, message string, args ...interface{}) {
		assert.Equal(t, 1, exitCode)
		message = fmt.Sprintf(message, args)
		assert.Contains(t, message, "Activation code required for on-prem registration.")
		panic(breakOutWithPanicMessageOnprem)
	}

	defer func() {
		if errInterface := recover(); errInterface != nil {
			assert.Equal(t, breakOutWithPanicMessageOnprem, errInterface)
		}
	}()
	main()
	assert.True(t, false, "Should never reach here because of exit")
}

func TestMain_Register_InvalidEnvironment_Failed(t *testing.T) {
	defer storeMockedFunctionsOnprem()()
	defer setArgsAndRestoreOnprem("/some/path/setupcli", "-env", "dummyEnv", "-register")()

	hasElevatedPermissions = func() error {
		return nil
	}

	osExit = func(exitCode int, log log.T, message string, args ...interface{}) {
		assert.Equal(t, 1, exitCode)
		message = fmt.Sprintf(message, args)
		assert.Contains(t, message, "Invalid environment.")
		panic(breakOutWithPanicMessageOnprem)
	}

	defer func() {
		if errInterface := recover(); errInterface != nil {
			assert.Equal(t, breakOutWithPanicMessageOnprem, errInterface)
		}
	}()
	main()
	assert.True(t, false, "Should never reach here because of exit")
}

func TestMain_SSMSetupCLI_NoInstallNoReg_Failed(t *testing.T) {
	defer storeMockedFunctionsOnprem()()
	defer setArgsAndRestoreOnprem("/some/path/setupcli", "-version", "dummy")()

	hasElevatedPermissions = func() error {
		return nil
	}

	osExit = func(exitCode int, log log.T, message string, args ...interface{}) {
		assert.Equal(t, 1, exitCode)
		message = fmt.Sprintf(message, args)
		fmt.Print(message)
		fmt.Print(args)
		assert.Contains(t, message, "Action required (-register or -install flag required). ")
		panic(breakOutWithPanicMessageOnprem)
	}

	defer func() {
		if errInterface := recover(); errInterface != nil {
			assert.Equal(t, breakOutWithPanicMessageOnprem, errInterface)
		}
	}()
	main()
	assert.True(t, false, "Should never reach here because of exit")
}

func TestMain_InstallAgent_StableVersion_Onprem_Success(t *testing.T) {
	defer storeMockedFunctionsOnprem()()
	defer setArgsAndRestoreOnprem("/some/path/setupcli", "-env", "onprem", "-install", "--version", "stable", "--region", "us-east-1")()

	hasElevatedPermissions = func() error {
		return nil
	}

	getPackageManager = func(log.T) (packagemanagers.IPackageManager, error) {
		managerMock := &pmMock.IPackageManager{}
		managerMock.On("GetInstalledAgentVersion").Return("2.1.2.2", nil)
		managerMock.On("IsAgentInstalled").Return(true, nil)
		managerMock.On("GetFileExtension").Return("test")
		return managerMock, nil
	}

	getServiceManager = func(log.T) (servicemanagers.IServiceManager, error) {
		managerMock := &smMock.IServiceManager{}
		managerMock.On("GetName").Return("ServiceManagerName")
		managerMock.On("GetAgentStatus").Return(common.Stopped, nil).Times(1)
		managerMock.On("StopAgent").Return(nil)
		managerMock.On("StartAgent").Return(nil)
		managerMock.On("GetAgentStatus").Return(common.Running, nil).Times(1)
		return managerMock, nil
	}
	helperInstallAgent = func(log log.T, pManager packagemanagers.IPackageManager, sManager servicemanagers.IServiceManager, folderPath string) error {
		return nil
	}
	helperUnInstallAgent = func(log log.T, pkgManager packagemanagers.IPackageManager, sManager servicemanagers.IServiceManager, installedVersionPath string) error {
		return nil
	}
	utilityCheckSum = func(filePath string) (hash string, err error) {
		return "", nil
	}
	getConfigurationManager = func() configurationmanager.IConfigurationManager {
		cfgManagerMock := &cmMock.IConfigurationManager{}
		cfgManagerMock.On("CreateUpdateAgentConfigWithOnPremIdentity").Return(nil)
		cfgManagerMock.On("ConfigureAgent", mock.Anything).Return(nil)
		cfgManagerMock.On("IsConfigAvailable", mock.Anything).Return(nil)
		return cfgManagerMock
	}
	stableVersion := "3.2.0.0"
	getDownloadManager = func(log log.T, region string, manifestUrl string, updateInfo updateinfo.T, setupCLIArtifactsPath string, isNano bool) downloadmanager.IDownloadManager {
		managerMock := &dmMock.IDownloadManager{}
		managerMock.On("DownloadLatestSSMSetupCLI", mock.Anything, mock.Anything).Return(nil).Once()
		managerMock.On("GetStableVersion").Return(stableVersion, nil).Once()
		managerMock.On("GetLatestVersion").Return(agentVersioning.Version, nil).Once()

		// this mocks stable version
		managerMock.On("DownloadArtifacts", stableVersion, mock.Anything, mock.Anything).Return(nil).Once()
		managerMock.On("DownloadSignatureFile", mock.Anything, mock.Anything, mock.Anything).Return("sign1", nil).Once()
		return managerMock
	}
	getVerificationManager = func() (verificationmanagers.IVerificationManager, error) {
		verificationManager := &vmMock.IVerificationManager{}
		verificationManager.On("VerifySignature", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
		return verificationManager, nil
	}
	osExit = func(exitCode int, log log.T, message string, args ...interface{}) {
		if exitCode == 0 {
			return
		}
		panic(fmt.Sprintf("should not receive non zero exit code. Err: %v; Msg: %v; Args: %v", exitCode, message, args))
	}

	main()
	assert.Equal(t, true, install)                        // install flag should be set
	assert.Equal(t, utility.StableVersionString, version) // version should be set as Stable for default
}

func TestMain_InstallAgent_LatestVersion_Onprem_Success(t *testing.T) {
	defer storeMockedFunctionsOnprem()()
	defer setArgsAndRestoreOnprem("/some/path/setupcli", "-env", "onprem", "-install", "--version", "latest", "--region", "us-east-1")()

	hasElevatedPermissions = func() error {
		return nil
	}

	getPackageManager = func(log.T) (packagemanagers.IPackageManager, error) {
		managerMock := &pmMock.IPackageManager{}
		managerMock.On("GetInstalledAgentVersion").Return("2.1.2.2", nil)
		managerMock.On("IsAgentInstalled").Return(true, nil)
		managerMock.On("GetFileExtension").Return("test")
		return managerMock, nil
	}

	getServiceManager = func(log.T) (servicemanagers.IServiceManager, error) {
		managerMock := &smMock.IServiceManager{}
		managerMock.On("GetName").Return("ServiceManagerName")
		managerMock.On("GetAgentStatus").Return(common.Stopped, nil).Times(1)
		managerMock.On("StopAgent").Return(nil)
		managerMock.On("StartAgent").Return(nil)
		managerMock.On("GetAgentStatus").Return(common.Running, nil).Times(1)
		return managerMock, nil
	}
	helperInstallAgent = func(log log.T, pManager packagemanagers.IPackageManager, sManager servicemanagers.IServiceManager, folderPath string) error {
		return nil
	}
	helperUnInstallAgent = func(log log.T, pkgManager packagemanagers.IPackageManager, sManager servicemanagers.IServiceManager, installedVersionPath string) error {
		return nil
	}
	utilityCheckSum = func(filePath string) (hash string, err error) {
		return "", nil
	}
	getConfigurationManager = func() configurationmanager.IConfigurationManager {
		cfgManagerMock := &cmMock.IConfigurationManager{}
		cfgManagerMock.On("CreateUpdateAgentConfigWithOnPremIdentity").Return(nil)
		cfgManagerMock.On("ConfigureAgent", mock.Anything).Return(nil)
		cfgManagerMock.On("IsConfigAvailable", mock.Anything).Return(nil)
		return cfgManagerMock
	}
	latestVersion := "3.0.0.0"
	getDownloadManager = func(log log.T, region string, manifestUrl string, updateInfo updateinfo.T, setupCLIArtifactsPath string, isNano bool) downloadmanager.IDownloadManager {
		managerMock := &dmMock.IDownloadManager{}
		managerMock.On("DownloadLatestSSMSetupCLI", mock.Anything, mock.Anything).Return(nil).Once()
		managerMock.On("GetLatestVersion").Return(latestVersion, nil).Once()

		// this mocks stable version
		managerMock.On("DownloadArtifacts", latestVersion, mock.Anything, mock.Anything).Return(nil).Once()
		managerMock.On("DownloadSignatureFile", mock.Anything, mock.Anything, mock.Anything).Return("sign1", nil).Once()
		return managerMock
	}
	getVerificationManager = func() (verificationmanagers.IVerificationManager, error) {
		verificationManager := &vmMock.IVerificationManager{}
		verificationManager.On("VerifySignature", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
		return verificationManager, nil
	}
	osExit = func(exitCode int, log log.T, message string, args ...interface{}) {
		if exitCode == 0 {
			return
		}
		panic(fmt.Sprintf("should not receive non zero exit code. Err: %v; Msg: %v; Args: %v", exitCode, message, args))
	}

	main()
	assert.Equal(t, true, install)                        // install flag should be set
	assert.Equal(t, utility.LatestVersionString, version) // version should be set as latest
}

func TestMain_InstallAgent_AlreadyInstalledVersion_Onprem_Success(t *testing.T) {
	defer storeMockedFunctionsOnprem()()
	defer setArgsAndRestoreOnprem("/some/path/setupcli", "-env", "onprem", "-install", "--region", "us-east-1")()

	hasElevatedPermissions = func() error {
		return nil
	}

	getPackageManager = func(log.T) (packagemanagers.IPackageManager, error) {
		managerMock := &pmMock.IPackageManager{}
		managerMock.On("GetInstalledAgentVersion").Return("2.1.2.2", nil)
		managerMock.On("IsAgentInstalled").Return(true, nil)
		managerMock.On("GetFileExtension").Return("test")
		return managerMock, nil
	}

	getServiceManager = func(log.T) (servicemanagers.IServiceManager, error) {
		managerMock := &smMock.IServiceManager{}
		managerMock.On("GetName").Return("ServiceManagerName")
		managerMock.On("GetAgentStatus").Return(common.Stopped, nil).Times(1)
		managerMock.On("StopAgent").Return(nil)
		managerMock.On("StartAgent").Return(nil)
		managerMock.On("GetAgentStatus").Return(common.Running, nil).Times(1)
		return managerMock, nil
	}
	installInitiated := false
	helperInstallAgent = func(log log.T, pManager packagemanagers.IPackageManager, sManager servicemanagers.IServiceManager, folderPath string) error {
		installInitiated = true
		return nil
	}
	helperUnInstallAgent = func(log log.T, pkgManager packagemanagers.IPackageManager, sManager servicemanagers.IServiceManager, installedVersionPath string) error {
		return nil
	}
	utilityCheckSum = func(filePath string) (hash string, err error) {
		return "", nil
	}
	getConfigurationManager = func() configurationmanager.IConfigurationManager {
		cfgManagerMock := &cmMock.IConfigurationManager{}
		cfgManagerMock.On("CreateUpdateAgentConfigWithOnPremIdentity").Return(nil)
		cfgManagerMock.On("ConfigureAgent", mock.Anything).Return(nil)
		cfgManagerMock.On("IsConfigAvailable", mock.Anything).Return(nil)
		return cfgManagerMock
	}
	latestVersion := "3.0.0.0"
	getDownloadManager = func(log log.T, region string, manifestUrl string, updateInfo updateinfo.T, setupCLIArtifactsPath string, isNano bool) downloadmanager.IDownloadManager {
		managerMock := &dmMock.IDownloadManager{}
		managerMock.On("DownloadLatestSSMSetupCLI", mock.Anything, mock.Anything).Return(nil).Once()
		managerMock.On("GetLatestVersion").Return(latestVersion, nil).Once()

		// this mocks stable version
		managerMock.On("DownloadArtifacts", latestVersion, mock.Anything, mock.Anything).Return(nil).Once()
		managerMock.On("DownloadSignatureFile", mock.Anything, mock.Anything, mock.Anything).Return("sign1", nil).Once()
		return managerMock
	}
	getVerificationManager = func() (verificationmanagers.IVerificationManager, error) {
		verificationManager := &vmMock.IVerificationManager{}
		verificationManager.On("VerifySignature", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
		return verificationManager, nil
	}
	osExit = func(exitCode int, log log.T, message string, args ...interface{}) {
		if exitCode == 0 {
			return
		}
		panic(fmt.Sprintf("should not receive non zero exit code. Err: %v; Msg: %v; Args: %v", exitCode, message, args))
	}

	main()
	assert.Equal(t, true, install)
	// install should not happen when version flag not passed
	assert.Equal(t, false, installInitiated)
	// version should be blank when version flag not passed
	assert.Equal(t, "", version)
}
