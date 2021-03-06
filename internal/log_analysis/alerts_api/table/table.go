// Package table manages all of the Dynamo calls (query, scan, get, write, etc).
package table

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
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

const (
	RuleIDKey          = "ruleId"
	AlertIDKey         = "id"
	TimePartitionKey   = "timePartition"
	TimePartitionValue = "defaultPartition"
)

// API defines the interface for the alerts table which can be used for mocking.
type API interface {
	GetAlert(*string) (*AlertItem, error)
	ListByRule(string, *string, *int) ([]*AlertItem, *string, error)
	ListAll(*string, *int) ([]*AlertItem, *string, error)
}

// AlertsTable encapsulates a connection to the Dynamo alerts table.
type AlertsTable struct {
	AlertsTableName                    string
	RuleIDCreationTimeIndexName        string
	TimePartitionCreationTimeIndexName string
	Client                             dynamodbiface.DynamoDBAPI
}

// The AlertsTable must satisfy the API interface.
var _ API = (*AlertsTable)(nil)

// DynamoItem is a type alias for the item format expected by the Dynamo SDK.
type DynamoItem = map[string]*dynamodb.AttributeValue

// AlertItem is a DDB representation of an Alert
type AlertItem struct {
	AlertID      string    `json:"id"`
	RuleID       string    `json:"ruleId"`
	CreationTime time.Time `json:"creationTime"`
	UpdateTime   time.Time `json:"updateTime"`
	Severity     string    `json:"severity"`
	EventCount   int       `json:"eventCount"`
	LogTypes     []string  `json:"logTypes"`
}
