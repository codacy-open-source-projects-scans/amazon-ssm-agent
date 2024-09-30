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

package sharedprovider

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/aws/amazon-ssm-agent/agent/log"
	"github.com/aws/amazon-ssm-agent/common/runtimeconfig"
)

const (
	// refreshMinutesBeforeExpiryMinutes defines the amount of time the client will wait until it tries to read credentials from disk
	refreshBeforeExpiryDuration = 10 * time.Minute
	providerName                = "SharedCredentialsProvider"
)

var (
	emptyCredential      = credentials.Value{ProviderName: providerName}
	newSharedCredentials = credentials.NewSharedCredentials
	newRuntimeConfig     = runtimeconfig.NewIdentityRuntimeConfigClient
)

type ISharedCredentialsProvider interface {
	credentials.Provider
	credentials.Expirer
	SetExpiration(expiration time.Time, window time.Duration)
	RetrieveWithContext(ctx context.Context) (credentials.Value, error)
}

// sharedCredentialsProvider implements the AWS SDK credential provider, and is used to create AWS client.
// It retrieves credentials from the shared credentials on disk, and keeps track if those credentials are expired.
type sharedCredentialsProvider struct {
	credentials.Expiry

	log                   log.T
	runtimeConfigClient   runtimeconfig.IIdentityRuntimeConfigClient
	identityRuntimeConfig runtimeconfig.IdentityRuntimeConfig
	getTimeNow            func() time.Time
}

// NewCredentialsProvider initializes a shared provider that loads credentials that were saved disk
func NewCredentialsProvider(log log.T) ISharedCredentialsProvider {
	log = log.WithContext("[SharedCredentialsProvider]")

	return &sharedCredentialsProvider{
		log:        log,
		getTimeNow: time.Now,
	}
}

// Retrieve retrieves credentials from the shared profile
// Error will be returned if the request fails, or unable to extract
// the desired credentials.
func (s *sharedCredentialsProvider) Retrieve() (credentials.Value, error) {
	return s.RetrieveWithContext(context.Background())
}

// RetrieveWithContext retrieves credentials from the shared credentials file
// Error will be returned if the request fails, or unable to extract
// the desired credentials.
func (s *sharedCredentialsProvider) RetrieveWithContext(ctx context.Context) (credentials.Value, error) {
	runtimeConfigClient := newRuntimeConfig()
	// before sharedCredentialsProvider is initialized, we check if the runtime config exists
	config, err := runtimeConfigClient.GetConfig()
	if err != nil {
		return emptyCredential, fmt.Errorf("unable to read runtime config for ShareFile information. Err: %w", err)
	}

	if config.ShareFile == "" {
		return emptyCredential, fmt.Errorf("runtime config has an empty ShareFile")
	}

	// If credentials are already expired, return error
	if config.CredentialsExpiresAt.Before(s.getTimeNow()) {
		return emptyCredential, fmt.Errorf("shared credentials are already expired, they were retrieved at %v and expired at %v", config.CredentialsRetrievedAt.Format(time.RFC3339), config.CredentialsExpiresAt.Format(time.RFC3339))
	}

	credsProvider := newSharedCredentials(config.ShareFile, config.ShareProfile)
	creds, err := credsProvider.Get()
	if err != nil {
		return emptyCredential, err
	}

	creds.ProviderName = providerName

	expirationWindow := s.getTimeWindow()
	// If credentials currently saved credentials expire in less than 'refreshBeforeExpiryDuration', no expiry window should be set
	if config.CredentialsExpiresAt.Before(s.getTimeNow().Add(refreshBeforeExpiryDuration)) {
		expirationWindow = time.Duration(0)
	}

	s.SetExpiration(config.CredentialsExpiresAt, expirationWindow)

	return creds, err
}

func (s *sharedCredentialsProvider) getTimeWindow() time.Duration {
	// Random jitter of up to 1 minute in case multiple workers are running, we want to spread read of the runtime config and credentials file
	randomJitterDuration := time.Second * time.Duration(rand.Intn(60))
	return refreshBeforeExpiryDuration + randomJitterDuration
}
