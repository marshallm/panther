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

import dayjs from 'dayjs';
import * as React from 'react';
import * as Yup from 'yup';
import {
  ActiveSuppressCount,
  ComplianceItem,
  ComplianceStatusCounts,
  Integration,
  OrganizationReportBySeverity,
  ResourceDetails,
  ResourceSummary,
  ScannedResources,
} from 'Generated/schema';
import {
  INCLUDE_DIGITS_REGEX,
  INCLUDE_LOWERCASE_REGEX,
  INCLUDE_SPECIAL_CHAR_REGEX,
  INCLUDE_UPPERCASE_REGEX,
} from 'Source/constants';
import mapValues from 'lodash-es/mapValues';
import sum from 'lodash-es/sum';
import { Box, ColumnProps, Label } from 'pouncejs';
import { ErrorResponse } from 'apollo-link-error';
import { ApolloError } from '@apollo/client';

// Generate a new secret code that contains metadata of issuer and user email
export const formatSecretCode = (code: string, email: string): string => {
  const issuer = 'Panther';
  return `otpauth://totp/${email}?secret=${code}&issuer=${issuer}`;
};

export const getArnRegexForService = (awsService: string) => {
  return new RegExp(`arn:aws:${awsService.toLowerCase()}:([a-z]){2}-([a-z])+-[0-9]:\\d{12}:.+`);
};

export const createYupPasswordValidationSchema = () =>
  Yup.string()
    .required()
    .min(14)
    .matches(INCLUDE_DIGITS_REGEX, 'Include at least 1 digit')
    .matches(INCLUDE_LOWERCASE_REGEX, 'Include at least 1 lowercase character')
    .matches(INCLUDE_UPPERCASE_REGEX, 'Include at least 1 uppercase character')
    .matches(INCLUDE_SPECIAL_CHAR_REGEX, 'Include at least 1 special character');

/**
 * checks whether the input is a valid UUID
 */
export const isGuid = (str: string) =>
  /^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/.test(str);

/**
 * caps the first letter of a string
 */
export const capitalize = (str: string) => str.charAt(0).toUpperCase() + str.slice(1);

/* eslint-disable react/display-name */
export const generateEnumerationColumn = (startIndex = 0) => {
  const enumerationColumn: ColumnProps<{}> = {
    key: 'enumeration',
    flex: '0 1 auto',
    renderColumnHeader: () => <Box ml={2} width={20} />,
    renderCell: (item: any, index: number) => (
      <Label size="medium" ml={2} minWidth={20}>
        {startIndex + index + 1}
      </Label>
    ),
  };

  return enumerationColumn;
};
/* eslint-enable react/display-name */

/**
 * Given a server-received DateTime string, creates a proper display text for it. We manually
 * calculate the offset cause there is no available format-string that can display the UTC offset
 * as a single digit (all of them display it either as 03:00 or as 0300) and require string
 * manipulation which is harder
 * */
export const formatDatetime = (datetime: string) => {
  // get the offset minutes and calculate the hours from them
  const utcOffset = dayjs(datetime).utcOffset() / 60;

  // properly format the date
  return dayjs(datetime).format(
    `YYYY-MM-DD HH:mm G[M]T${utcOffset > 0 ? '+' : ''}${utcOffset !== 0 ? utcOffset : ''}`
  );
};

/** Converts any value of the object that is an array to a comma-separated string */
export const convertObjArrayValuesToCsv = (obj: { [key: string]: any }) =>
  mapValues(obj, v => (Array.isArray(v) ? v.join(',') : v));

/**
 * makes sure that it properly formats a JSON struct in order to be properly displayed within the
 * editor
 * @param code valid JSON
 * @returns String
 */
export const formatJSON = (code: { [key: string]: number | string }) =>
  JSON.stringify(code, null, '\t');

/**
 * Extends the resource by adding an `integrationLabel` field. We define two overloads for this
 * function
 * @param resource A resource
 * @param integrations A list of integrations with at least (integrationId & integrationType)
 */

function extendResourceWithIntLabel(
  resource: ResourceSummary,
  integrations: (Partial<Integration> & Pick<Integration, 'integrationId' | 'integrationLabel'>)[]
): ResourceSummary & Pick<Integration, 'integrationLabel'>;

function extendResourceWithIntLabel(
  resource: ResourceDetails,
  integrations: (Partial<Integration> & Pick<Integration, 'integrationId' | 'integrationLabel'>)[]
): ResourceDetails & Pick<Integration, 'integrationLabel'>;

function extendResourceWithIntLabel(
  resource: ComplianceItem,
  integrations: (Partial<Integration> & Pick<Integration, 'integrationId' | 'integrationLabel'>)[]
): ComplianceItem & Pick<Integration, 'integrationLabel'>;

function extendResourceWithIntLabel(
  resource: any,
  integrations: (Partial<Integration> & Pick<Integration, 'integrationId' | 'integrationLabel'>)[]
) {
  const matchingIntegration = integrations.find(i => i.integrationId === resource.integrationId);
  return {
    ...resource,
    integrationLabel: matchingIntegration?.integrationLabel || 'Cannot find account',
  };
}

export const extendResourceWithIntegrationLabel = extendResourceWithIntLabel;

/**
 * sums up the total number of items based on the active/suppresed count breakdown that the API
 * exposes
 */
export const getComplianceItemsTotalCount = (totals: ActiveSuppressCount) => {
  return (
    totals.active.pass +
    totals.active.fail +
    totals.active.error +
    totals.suppressed.pass +
    totals.suppressed.fail +
    totals.suppressed.error
  );
};

/**
 * sums up the total number of policies based on the severity and compliance status count breakdown
 * that the API exposes. With this function we can choose to aggregate only the failing policies
 * for a severity or even all of them, simply by passing the corresponding array of statuses to
 * aggregate.
 *
 * For example:
 * countPoliciesBySeverityAndStatus([], 'critical', ['fail', 'error']) would count the critical
 * policies that are either failing or erroring
 */
export const countPoliciesBySeverityAndStatus = (
  data: OrganizationReportBySeverity,
  severity: keyof OrganizationReportBySeverity,
  complianceStatuses: (keyof ComplianceStatusCounts)[]
) => {
  return sum(complianceStatuses.map(complianceStatus => data[severity][complianceStatus]));
};

/**
 * sums up the total number of resources based on the compliance status count breakdown
 * that the API exposes. With this function we can choose to aggregate only the failing resources
 * or even all of them, simply by passing the corresponding array of statuses to
 * aggregate.
 *
 * For example:
 * countResourcesByStatus([], ['fail', 'error']) would count the resources that are either failing
 * or erroring
 */
export const countResourcesByStatus = (
  data: ScannedResources,
  complianceStatuses: (keyof ComplianceStatusCounts)[]
) => {
  // aggregates the list of "totals" for each resourceType. The "total" for a resource type is the
  // aggregation of ['fail', 'error', ...] according to the parameter passed by the user
  return sum(
    data.byType.map(({ count }) =>
      sum(complianceStatuses.map(complianceStatus => count[complianceStatus]))
    )
  );
};

/**
 * A function that takes the whole GraphQL error as a payload and returns the message that should
 * be shown to the user
 */
export const extractErrorMessage = (error: ApolloError | ErrorResponse) => {
  // If there is a network error show something (we are already showing the network-error-modal though)
  if (error.networkError) {
    return "Can't perform any action because of a problem with your network";
  }

  // If there are no networkErrors or graphQL errors, then show the fallback
  if (!error.graphQLErrors || !error.graphQLErrors.length) {
    return 'A unpredicted server error has occured';
  }

  // isolate the first GraphQL error. Currently all of our APIs return a single error. If we ever
  // return multiple, we should handle that for all items within the `graphQLErrors` key
  const { errorType, message } = error.graphQLErrors[0];
  switch (errorType) {
    case '401':
    case '403':
      return message || 'You are not authorized to perform this request';
    case '404':
      return message || "The resource you requested couldn't be found on our servers";
    default:
      return message;
  }
};

// Copies a text to clipboard, with fallback for Safari and old-Edge
export const copyTextToClipboard = (text: string) => {
  if (navigator.clipboard) {
    navigator.clipboard.writeText(text);
  } else {
    const container = document.querySelector('[role="dialog"] [role="document"]') || document.body;
    const textArea = document.createElement('textarea');
    textArea.innerHTML = text;
    textArea.style.position = 'fixed'; // avoid scrolling to bottom
    container.appendChild(textArea);
    textArea.focus();
    textArea.select();
    document.execCommand('copy');
    container.removeChild(textArea);
  }
};

export const isNumber = (value: string) => /^-{0,1}\d+$/.test(value);
