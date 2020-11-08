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

// ExternalUpdate is an external update.
type ExternalUpdate struct {
	ConnectionDetails ConnectionDetails `json:"connectionDetails,omitempty"`
}

// UpdateRequest is a request to create an external resource.
type UpdateRequest struct {
	Managed resource.Managed `json:"managed"`
}

// UpdateResponse is the response to a create request.
type UpdateResponse struct {
	External   ExternalUpdate `json:"externalUpdate,omitempty"`
	ErrMessage string         `json:"errMessage,omitempty"`
}
