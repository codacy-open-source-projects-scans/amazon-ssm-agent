// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	log "github.com/aws/amazon-ssm-agent/agent/log"
	mock "github.com/stretchr/testify/mock"

	packagemanagers "github.com/aws/amazon-ssm-agent/agent/setupcli/managers/packagemanagers"

	servicemanagers "github.com/aws/amazon-ssm-agent/agent/setupcli/managers/servicemanagers"

	verificationmanagers "github.com/aws/amazon-ssm-agent/agent/setupcli/managers/verificationmanagers"
)

// IPackageManager is an autogenerated mock type for the IPackageManager type
type IPackageManager struct {
	mock.Mock
}

// GetFileExtension provides a mock function with given fields:
func (_m *IPackageManager) GetFileExtension() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetFilesReqForInstall provides a mock function with given fields: _a0
func (_m *IPackageManager) GetFilesReqForInstall(_a0 log.T) []string {
	ret := _m.Called(_a0)

	var r0 []string
	if rf, ok := ret.Get(0).(func(log.T) []string); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// GetInstalledAgentVersion provides a mock function with given fields:
func (_m *IPackageManager) GetInstalledAgentVersion() (string, error) {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetName provides a mock function with given fields:
func (_m *IPackageManager) GetName() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetSupportedServiceManagers provides a mock function with given fields:
func (_m *IPackageManager) GetSupportedServiceManagers() []servicemanagers.ServiceManager {
	ret := _m.Called()

	var r0 []servicemanagers.ServiceManager
	if rf, ok := ret.Get(0).(func() []servicemanagers.ServiceManager); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]servicemanagers.ServiceManager)
		}
	}

	return r0
}

// GetSupportedVerificationManager provides a mock function with given fields:
func (_m *IPackageManager) GetSupportedVerificationManager() verificationmanagers.VerificationManager {
	ret := _m.Called()

	var r0 verificationmanagers.VerificationManager
	if rf, ok := ret.Get(0).(func() verificationmanagers.VerificationManager); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(verificationmanagers.VerificationManager)
	}

	return r0
}

// GetType provides a mock function with given fields:
func (_m *IPackageManager) GetType() packagemanagers.PackageManager {
	ret := _m.Called()

	var r0 packagemanagers.PackageManager
	if rf, ok := ret.Get(0).(func() packagemanagers.PackageManager); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(packagemanagers.PackageManager)
	}

	return r0
}

// InstallAgent provides a mock function with given fields: _a0, folderPath
func (_m *IPackageManager) InstallAgent(_a0 log.T, folderPath string) error {
	ret := _m.Called(_a0, folderPath)

	var r0 error
	if rf, ok := ret.Get(0).(func(log.T, string) error); ok {
		r0 = rf(_a0, folderPath)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// IsAgentInstalled provides a mock function with given fields:
func (_m *IPackageManager) IsAgentInstalled() (bool, error) {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsManagerEnvironment provides a mock function with given fields:
func (_m *IPackageManager) IsManagerEnvironment() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// UninstallAgent provides a mock function with given fields: _a0, _a1
func (_m *IPackageManager) UninstallAgent(_a0 log.T, _a1 string) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(log.T, string) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
