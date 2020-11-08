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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/crossplane/crossplane-runtime/apis/core/v1alpha1"
)

// TypeReference is used to refer to a type for declaring compatibility.
type TypeReference struct {
	// APIVersion of the type.
	APIVersion string `json:"apiVersion"`

	// Kind of the type.
	Kind string `json:"kind"`
}

// A RegistrationSpec defines the desired state of a Registration.
type RegistrationSpec struct {
	TypeRef  TypeReference `json:"typeRef"`
	Endpoint string        `json:"endpoint"`
}

// A RegistrationStatus reflects the observed state of a Registration.
type RegistrationStatus struct {
	v1alpha1.ConditionedStatus `json:",inline"`
}

// +kubebuilder:object:root=true

// A Registration registers a new Crisscross controller.
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:scope=Cluster
type Registration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RegistrationSpec   `json:"spec"`
	Status RegistrationStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// RegistrationList contains a list of Registration.
type RegistrationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Registration `json:"items"`
}
