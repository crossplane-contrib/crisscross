/*
Copyright 2020 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package models

import (
	"github.com/crossplane/crossplane-runtime/pkg/resource"
)

// ExternalObservation is an external observation.
type ExternalObservation struct {
	ResourceExists          bool              `json:"resourceExists,omitempty"`
	ResourceUpToDate        bool              `json:"resourceUpToDate,omitempty"`
	ResourceLateInitialized bool              `json:"resourceLateInitialized,omitempty"`
	ConnectionDetails       ConnectionDetails `json:"connectionDetails,omitempty"`
}

// ConnectionDetails created or updated during an operation on an external
// resource, for example usernames, passwords, endpoints, ports, etc.
type ConnectionDetails map[string][]byte

// ObservationRequest is a request to create an external resource.
type ObservationRequest struct {
	Managed resource.Managed `json:"managed"`
}

// ObservationResponse is the response to a create request.
type ObservationResponse struct {
	External   ExternalObservation `json:"externalObservation,omitempty"`
	ErrMessage string              `json:"errMessage,omitempty"`
}
