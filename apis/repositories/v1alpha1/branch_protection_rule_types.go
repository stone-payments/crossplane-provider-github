/*
Copyright 2022 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

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

// BranchProtectionRuleParameters defines the desired state of a GitHub Repository BranchProtectionRule.
type BranchProtectionRuleParameters struct {
	// The pattern to be protected.
	//
	// +immutable
	Pattern string `json:"pattern"`

	// The name of the Repository owner.
	// The owner can be an organization or an user.
	//
	// +immutable
	Owner string `json:"owner"`

	// The name of the Repository.
	//
	// +immutable
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

	// Actors who may force push to the protected branch. User, app, and team restrictions are only
	// available for organization-owned repositories. Defaults to disabled.
	//
	// Users should be specified in the format: /{username}.
	// Teams should be specified in the format: /{organization}/{team-slug}
	//
	// +optional
	BypassForcePushAllowances []string `json:"bypassForcePushAllowances,omitempty"`

	// A list of actors able to bypass PRs for this branch protection rule. Defaults to disabled.
	//
	// Users should be specified in the format: /{username}.
	// Teams should be specified in the format: /{organization}/{team-slug}
	//
	// +optional
	BypassPullRequestAllowances []string `json:"bypassPullRequestAllowances,omitempty"`

	// Whether new commits pushed to matching branches dismiss pull request review approvals.
	//
	// +optional
	DismissesStaleReviews *bool `json:"dismissesStaleReviews,omitempty"`

	// Whether admins can bypass branch protection rules.
	//
	// +optional
	IsAdminEnforced *bool `json:"isAdminEnforced,omitempty"`

	// Actors who may push to the protected branch. Defaults to disabled.
	//
	// Users should be specified in the format: /{username}.
	// Teams should be specified in the format: /{organization}/{team-slug}
	//
	// +optional
	PushAllowances []string `json:"pushAllowances,omitempty"`

	// Number of approving reviews required in the pull request.
	// Must be a number between 0-6.
	//
	// +optional
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=6
	// +kubebuilder:default=0
	RequiredApprovingReviewCount *int32 `json:"requiredApprovingReviewCount,omitempty"`

	// List of required status check contexts that must pass for commits to be
	// accepted to matching branches.
	//
	// +optional
	RequiredStatusCheckContexts []string `json:"requiredStatusCheckContexts,omitempty"`

	// Whether reviews from code owners are required to update matching branches.
	//
	// +optional
	RequiresCodeOwnerReviews *bool `json:"requiresCodeOwnerReviews,omitempty"`

	// Whether commits are required to be signed.
	//
	// +optional
	RequiresCommitSignatures *bool `json:"requiresCommitSignatures,omitempty"`

	// Whether conversations are required to be resolved before merging.
	//
	// +optional
	RequiresConversationResolution *bool `json:"requiresConversationResolution,omitempty"`

	// Whether merge commits are prohibited from being pushed to this branch.
	//
	// +optional
	RequiresLinearHistory *bool `json:"requiresLinearHistory,omitempty"`

	// Whether branches are required to be up to date before merging.
	//
	// +optional
	// +kubebuilder:default=false
	RequiresStrictStatusChecks *bool `json:"requiresStrictStatusChecks,omitempty"`
}

// BranchProtectionRuleSpec defines the desired state of a BranchProtectionRule.
type BranchProtectionRuleSpec struct {
	xpv1.ResourceSpec `json:",inline"`

	ForProvider BranchProtectionRuleParameters `json:"forProvider"`
}

// BranchProtectionRuleObservation is the representation of the current state that is observed
type BranchProtectionRuleObservation struct {
	// Global ID that represents this BranchProtectionRule
	ID string `json:"id,omitempty"`

	// Repository global node ID. Will be inferred by the
	// spec name and owner field.
	RepositoryID string `json:"repositoryId,omitempty"`
}

// BranchProtectionRuleStatus represents the observed state of a BranchProtectionRule.
type BranchProtectionRuleStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          BranchProtectionRuleObservation `json:"atProvider"`
}

// +kubebuilder:object:root=true

// A BranchProtectionRule is a managed resource that represents a GitHub Repository BranchProtectionRule configuration
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,github}
type BranchProtectionRule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BranchProtectionRuleSpec   `json:"spec"`
	Status BranchProtectionRuleStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// BranchProtectionRuleList contains a list of BranchProtectionRule
type BranchProtectionRuleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BranchProtectionRule `json:"items"`
}
