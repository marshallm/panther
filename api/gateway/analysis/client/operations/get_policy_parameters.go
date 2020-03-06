// Code generated by go-swagger; DO NOT EDIT.

package operations

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

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	strfmt "github.com/go-openapi/strfmt"
)

// NewGetPolicyParams creates a new GetPolicyParams object
// with the default values initialized.
func NewGetPolicyParams() *GetPolicyParams {
	var ()
	return &GetPolicyParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetPolicyParamsWithTimeout creates a new GetPolicyParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetPolicyParamsWithTimeout(timeout time.Duration) *GetPolicyParams {
	var ()
	return &GetPolicyParams{

		timeout: timeout,
	}
}

// NewGetPolicyParamsWithContext creates a new GetPolicyParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetPolicyParamsWithContext(ctx context.Context) *GetPolicyParams {
	var ()
	return &GetPolicyParams{

		Context: ctx,
	}
}

// NewGetPolicyParamsWithHTTPClient creates a new GetPolicyParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetPolicyParamsWithHTTPClient(client *http.Client) *GetPolicyParams {
	var ()
	return &GetPolicyParams{
		HTTPClient: client,
	}
}

/*GetPolicyParams contains all the parameters to send to the API endpoint
for the get policy operation typically these are written to a http.Request
*/
type GetPolicyParams struct {

	/*ID
	  User-specified unique rule/policy ID

	*/
	ID string
	/*VersionID
	  The version of the analysis to retrieve

	*/
	VersionID *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get policy params
func (o *GetPolicyParams) WithTimeout(timeout time.Duration) *GetPolicyParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get policy params
func (o *GetPolicyParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get policy params
func (o *GetPolicyParams) WithContext(ctx context.Context) *GetPolicyParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get policy params
func (o *GetPolicyParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get policy params
func (o *GetPolicyParams) WithHTTPClient(client *http.Client) *GetPolicyParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get policy params
func (o *GetPolicyParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithID adds the id to the get policy params
func (o *GetPolicyParams) WithID(id string) *GetPolicyParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the get policy params
func (o *GetPolicyParams) SetID(id string) {
	o.ID = id
}

// WithVersionID adds the versionID to the get policy params
func (o *GetPolicyParams) WithVersionID(versionID *string) *GetPolicyParams {
	o.SetVersionID(versionID)
	return o
}

// SetVersionID adds the versionId to the get policy params
func (o *GetPolicyParams) SetVersionID(versionID *string) {
	o.VersionID = versionID
}

// WriteToRequest writes these params to a swagger request
func (o *GetPolicyParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// query param id
	qrID := o.ID
	qID := qrID
	if qID != "" {
		if err := r.SetQueryParam("id", qID); err != nil {
			return err
		}
	}

	if o.VersionID != nil {

		// query param versionId
		var qrVersionID string
		if o.VersionID != nil {
			qrVersionID = *o.VersionID
		}
		qVersionID := qrVersionID
		if qVersionID != "" {
			if err := r.SetQueryParam("versionId", qVersionID); err != nil {
				return err
			}
		}

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
