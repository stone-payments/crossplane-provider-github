/*
Copyright 2022 The Crossplane Authors.

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

// OrphanRefParameters defines the desired state of a GitHub Repository OrphanRef.
type OrphanRefParameters struct {
	// The name of the Repository owner.
	// The owner can be an organization or an user.
	Owner string `json:"owner"`

	// The name of the Repository.
	//
	// +optional
	Repository *string `json:"repository,omitempty"`

	// RepositoryRef references a Repository and retrieves its name.
	//
	// +optional
	RepositoryRef *xpv1.Reference `json:"repositoryRef,omitempty"`

	// RepositorySelector selects a reference to a Repository.
	//
	// +optional
	RepositorySelector *xpv1.Selector `json:"repositorySelector,omitempty"`

	// The dummy file path.
	// +immutable
	Path string `json:"path"`

	// The initial commit message.
	// +optional
	// +kubebuilder:default="Initial Commit"
	// +immutable
	Message *string `json:"message,omitempty"`
}

// OrphanRefSpec defines the desired state of a OrphanRef.
type OrphanRefSpec struct {
	xpv1.ResourceSpec `json:",inline"`

	ForProvider OrphanRefParameters `json:"forProvider"`
}

// OrphanRefObservation is the representation of the current state that is observed
type OrphanRefObservation struct {
	URL string `json:"url,omitempty"`
	Ref string `json:"ref,omitempty"`
}

// OrphanRefStatus represents the observed state of a OrphanRef.
type OrphanRefStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          OrphanRefObservation `json:"atProvider"`
}

// +kubebuilder:object:root=true

// A OrphanRef is a managed resource that represents a GitHub Repository OrphanRef
// +kubebuilder:printcolumn:name="URL",type="string",JSONPath=".status.atProvider.url"
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,github}
type OrphanRef struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OrphanRefSpec   `json:"spec"`
	Status OrphanRefStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// OrphanRefList contains a list of OrphanRef
type OrphanRefList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OrphanRef `json:"items"`
}
