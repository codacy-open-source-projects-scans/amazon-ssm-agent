// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// IDownloadManager is an autogenerated mock type for the IDownloadManager type
type IDownloadManager struct {
	mock.Mock
}

// DownloadArtifacts provides a mock function with given fields: version, manifestUrl, folderPath
func (_m *IDownloadManager) DownloadArtifacts(version string, manifestUrl string, folderPath string) error {
	ret := _m.Called(version, manifestUrl, folderPath)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string) error); ok {
		r0 = rf(version, manifestUrl, folderPath)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DownloadLatestSSMSetupCLI provides a mock function with given fields: artifactsStorePath, expectedCheckSum
func (_m *IDownloadManager) DownloadLatestSSMSetupCLI(artifactsStorePath string, expectedCheckSum string) error {
	ret := _m.Called(artifactsStorePath, expectedCheckSum)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(artifactsStorePath, expectedCheckSum)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DownloadSignatureFile provides a mock function with given fields: version, artifactsStorePath, extension
func (_m *IDownloadManager) DownloadSignatureFile(version string, artifactsStorePath string, extension string) (string, error) {
	ret := _m.Called(version, artifactsStorePath, extension)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string, string) string); ok {
		r0 = rf(version, artifactsStorePath, extension)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string) error); ok {
		r1 = rf(version, artifactsStorePath, extension)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLatestVersion provides a mock function with given fields:
func (_m *IDownloadManager) GetLatestVersion() (string, error) {
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

// GetStableVersion provides a mock function with given fields:
func (_m *IDownloadManager) GetStableVersion() (string, error) {
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