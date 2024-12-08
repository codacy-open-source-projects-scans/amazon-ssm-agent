// Copyright 2016 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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

//go:build windows
// +build windows

// Package platform contains platform specific utilities.
package platform

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/aws/amazon-ssm-agent/agent/log"
	"github.com/aws/amazon-ssm-agent/agent/versionutil"
)

// Win32_OperatingSystems https://msdn.microsoft.com/en-us/library/aa394239%28v=vs.85%29.aspx
const (
	// PRODUCT_DATA_CENTER_NANO_SERVER = 143
	ProductDataCenterNanoServer = "143"

	// PRODUCT_STANDARD_NANO_SERVER = 144
	ProductStandardNanoServer = "144"

	// WindowsServer2016Version represents Win32_OperatingSystemVersion https://learn.microsoft.com/en-us/windows/win32/sysinfo/operating-system-version
	WindowsServer2016Version = 10

	WindowsServer2025Version = "10.0.26100"
)

var (
	getPlatformVersionRef = getPlatformVersion
)

// isPlatformWindowsServer2012OrEarlier returns true if platform is Windows Server 2012 or earlier
func isPlatformWindowsServer2012OrEarlier(log log.T) (bool, error) {
	var platformVersion string
	var platformVersionInt int
	var err error

	if platformVersion, err = getPlatformVersionRef(log); err != nil {
		return false, err
	}
	versionParts := strings.Split(platformVersion, ".")
	if len(versionParts) == 0 {
		return false, fmt.Errorf("could not get the version from versionstring: %v", versionParts)
	}

	if platformVersionInt, err = strconv.Atoi(versionParts[0]); err != nil {
		return false, err
	}
	return platformVersionInt < WindowsServer2016Version, nil
}

// isPlatformWindowsServer2025OrLater returns true if current platform is Windows Server 2025 or later
func isPlatformWindowsServer2025OrLater(log log.T) (bool, error) {
	if platformVersion, err := getPlatformVersionRef(log); err != nil {
		return false, err
	} else {
		return isWindowsServer2025OrLater(platformVersion, log)
	}
}

// isWindowsServer2025OrLater returns true if passed platformVersion is the same as of Windows Server 2025 or later
func isWindowsServer2025OrLater(platformVersion string, log log.T) (bool, error) {
	log.Debugf("Checking if platform version: %s is Windows 2025 or later...", platformVersion)
	if result, err := versionutil.VersionCompare(platformVersion, WindowsServer2025Version); err != nil {
		return false, err
	} else {
		return result >= 0, nil
	}
}

// IsPlatformNanoServer returns true if SKU is 143 or 144
func isPlatformNanoServer(log log.T) (bool, error) {
	// Get platform sku information
	if sku, err := getPlatformSku(log); err != nil {
		log.Infof("Failed to fetch sku - %v", err)
		return false, err
	} else {
		// Return whether sku represents nano server
		return sku == ProductDataCenterNanoServer || sku == ProductStandardNanoServer, nil
	}
}

func getPlatformName(log log.T) (value string, err error) {
	if osData, err := getPlatformDetails(log); err != nil {
		return notAvailableMessage, err
	} else {
		return osData.Caption, nil
	}
}

func getPlatformType(_ log.T) (value string, err error) {
	return "windows", nil
}

func getPlatformVersion(log log.T) (value string, err error) {
	if osData, err := getPlatformDetails(log); err != nil {
		return notAvailableMessage, err
	} else {
		return osData.Version, nil
	}
}

func getPlatformSku(log log.T) (value string, err error) {
	if osData, err := getPlatformDetails(log); err != nil {
		return notAvailableMessage, err
	} else {
		return strconv.FormatUint(uint64(osData.OperatingSystemSKU), 10), nil
	}
}

func getPlatformDetails(log log.T) (osData Win32_OperatingSystem, err error) {
	if osData, err = GetSingleWMIObject(osData); err != nil {
		log.Errorf("Failed to fetch OS details from WMI: %v", err)
	}

	return osData, err
}

// fullyQualifiedDomainName returns the Fully Qualified Domain Name of the instance, otherwise the hostname
func fullyQualifiedDomainName(log log.T) string {
	var csData Win32_ComputerSystem
	var err error
	if csData, err = GetSingleWMIObject(csData); err != nil {
		log.Errorf("Failed to fetch computer system details from WMI: %v", err)
	}

	if csData.DNSHostName == "" || csData.Domain == "" {
		hostName, _ := os.Hostname()
		return hostName
	}

	return csData.DNSHostName + "." + csData.Domain
}
