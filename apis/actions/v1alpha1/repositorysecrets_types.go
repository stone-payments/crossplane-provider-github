/*
Copyright 2021 The Crossplane Authors.

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

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
)

// RepositorySecretParameters defines the desired state of a GitHub Repository Secrets.
type RepositorySecretParameters struct {
	// The name of the Repository owner.
	Owner string `json:"owner"`
	// The name of the repository.
	Repository string `json:"repository"`
	// The value of the secret
	Value xpv1.SecretKeySelector `json:"value"`
}

// RepositorySecretObservation are the observable fields of a Repository Secrets.
type RepositorySecretObservation struct {
	// The encrypted value stored in K8s Secret
	// +optional
	EncryptValue string `json:"encrypt_value,omitempty"`
	// Last updated time in Repository Secret GitHub
	// +optional
	LastUpdate string `json:"last_update,omitempty"`
}

// A RepositorySecretSpec defines the desired state of a Repository Secrets.
type RepositorySecretSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       RepositorySecretParameters `json:"forProvider"`
}

// A RepositorySecretStatus represents the observed state of a Repository Secrets.
type RepositorySecretStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          RepositorySecretObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A RepositorySecret is a managed resource that represents a GitHub Secrets
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,github}
type RepositorySecret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RepositorySecretSpec   `json:"spec"`
	Status RepositorySecretStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// RepositorySecretList contains a list of Secrets
type RepositorySecretList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RepositorySecret `json:"items"`
}
