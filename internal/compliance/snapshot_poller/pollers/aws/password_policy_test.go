package aws

/**
 * Panther is a scalable, powerful, cloud-native SIEM written in Golang/React.
 * Copyright (C) 2020 Panther Labs Inc
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	awsmodels "github.com/panther-labs/panther/internal/compliance/snapshot_poller/models/aws"
	"github.com/panther-labs/panther/internal/compliance/snapshot_poller/pollers/aws/awstest"
)

func TestGetPasswordPolicy(t *testing.T) {
	mockSvc := awstest.BuildMockIAMSvc([]string{"GetAccountPasswordPolicy"})

	out, err := getPasswordPolicy(mockSvc)
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}

func TestGetPasswordPolicyError(t *testing.T) {
	mockSvc := awstest.BuildMockIAMSvcError([]string{"GetAccountPasswordPolicy"})

	out, err := getPasswordPolicy(mockSvc)
	require.NotNil(t, err)
	assert.Nil(t, out)
}

func TestPasswordPolicyPoller(t *testing.T) {
	awstest.MockIAMForSetup = awstest.BuildMockIAMSvc([]string{"GetAccountPasswordPolicy"})

	AssumeRoleFunc = awstest.AssumeRoleMock
	IAMClientFunc = awstest.SetupMockIAM

	resources, err := PollPasswordPolicy(&awsmodels.ResourcePollerInput{
		AuthSource:          &awstest.ExampleAuthSource,
		AuthSourceParsedARN: awstest.ExampleAuthSourceParsedARN,
		IntegrationID:       awstest.ExampleIntegrationID,
		Timestamp:           &awstest.ExampleTime,
	})

	require.NoError(t, err)
	assert.Len(t, resources, 1)
	assert.Equal(t, "123456789012::AWS.PasswordPolicy", string(resources[0].ID))
}

func TestPasswordPolicyPollerError(t *testing.T) {
	awstest.MockIAMForSetup = awstest.BuildMockIAMSvcError([]string{"GetAccountPasswordPolicy"})

	AssumeRoleFunc = awstest.AssumeRoleMock
	IAMClientFunc = awstest.SetupMockIAM

	resources, err := PollPasswordPolicy(&awsmodels.ResourcePollerInput{
		AuthSource:          &awstest.ExampleAuthSource,
		AuthSourceParsedARN: awstest.ExampleAuthSourceParsedARN,
		IntegrationID:       awstest.ExampleIntegrationID,
		Timestamp:           &awstest.ExampleTime,
	})

	require.NoError(t, err)
	assert.Len(t, resources, 1)
}
