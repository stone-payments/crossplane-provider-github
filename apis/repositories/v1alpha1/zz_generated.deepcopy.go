//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2020 The Crossplane Authors.

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"github.com/crossplane/crossplane-runtime/apis/common/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BranchProtectionRule) DeepCopyInto(out *BranchProtectionRule) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BranchProtectionRule.
func (in *BranchProtectionRule) DeepCopy() *BranchProtectionRule {
	if in == nil {
		return nil
	}
	out := new(BranchProtectionRule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *BranchProtectionRule) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BranchProtectionRuleList) DeepCopyInto(out *BranchProtectionRuleList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]BranchProtectionRule, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BranchProtectionRuleList.
func (in *BranchProtectionRuleList) DeepCopy() *BranchProtectionRuleList {
	if in == nil {
		return nil
	}
	out := new(BranchProtectionRuleList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *BranchProtectionRuleList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BranchProtectionRuleObservation) DeepCopyInto(out *BranchProtectionRuleObservation) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BranchProtectionRuleObservation.
func (in *BranchProtectionRuleObservation) DeepCopy() *BranchProtectionRuleObservation {
	if in == nil {
		return nil
	}
	out := new(BranchProtectionRuleObservation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BranchProtectionRuleParameters) DeepCopyInto(out *BranchProtectionRuleParameters) {
	*out = *in
	if in.Repository != nil {
		in, out := &in.Repository, &out.Repository
		*out = new(string)
		**out = **in
	}
	if in.RepositoryRef != nil {
		in, out := &in.RepositoryRef, &out.RepositoryRef
		*out = new(v1.Reference)
		(*in).DeepCopyInto(*out)
	}
	if in.RepositorySelector != nil {
		in, out := &in.RepositorySelector, &out.RepositorySelector
		*out = new(v1.Selector)
		(*in).DeepCopyInto(*out)
	}
	if in.BypassForcePushAllowances != nil {
		in, out := &in.BypassForcePushAllowances, &out.BypassForcePushAllowances
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.BypassPullRequestAllowances != nil {
		in, out := &in.BypassPullRequestAllowances, &out.BypassPullRequestAllowances
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.DismissesStaleReviews != nil {
		in, out := &in.DismissesStaleReviews, &out.DismissesStaleReviews
		*out = new(bool)
		**out = **in
	}
	if in.IsAdminEnforced != nil {
		in, out := &in.IsAdminEnforced, &out.IsAdminEnforced
		*out = new(bool)
		**out = **in
	}
	if in.PushAllowances != nil {
		in, out := &in.PushAllowances, &out.PushAllowances
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.RequiredApprovingReviewCount != nil {
		in, out := &in.RequiredApprovingReviewCount, &out.RequiredApprovingReviewCount
		*out = new(int32)
		**out = **in
	}
	if in.RequiredStatusCheckContexts != nil {
		in, out := &in.RequiredStatusCheckContexts, &out.RequiredStatusCheckContexts
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.RequiresCodeOwnerReviews != nil {
		in, out := &in.RequiresCodeOwnerReviews, &out.RequiresCodeOwnerReviews
		*out = new(bool)
		**out = **in
	}
	if in.RequiresCommitSignatures != nil {
		in, out := &in.RequiresCommitSignatures, &out.RequiresCommitSignatures
		*out = new(bool)
		**out = **in
	}
	if in.RequiresConversationResolution != nil {
		in, out := &in.RequiresConversationResolution, &out.RequiresConversationResolution
		*out = new(bool)
		**out = **in
	}
	if in.RequiresLinearHistory != nil {
		in, out := &in.RequiresLinearHistory, &out.RequiresLinearHistory
		*out = new(bool)
		**out = **in
	}
	if in.RequiresStrictStatusChecks != nil {
		in, out := &in.RequiresStrictStatusChecks, &out.RequiresStrictStatusChecks
		*out = new(bool)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BranchProtectionRuleParameters.
func (in *BranchProtectionRuleParameters) DeepCopy() *BranchProtectionRuleParameters {
	if in == nil {
		return nil
	}
	out := new(BranchProtectionRuleParameters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BranchProtectionRuleSpec) DeepCopyInto(out *BranchProtectionRuleSpec) {
	*out = *in
	in.ResourceSpec.DeepCopyInto(&out.ResourceSpec)
	in.ForProvider.DeepCopyInto(&out.ForProvider)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BranchProtectionRuleSpec.
func (in *BranchProtectionRuleSpec) DeepCopy() *BranchProtectionRuleSpec {
	if in == nil {
		return nil
	}
	out := new(BranchProtectionRuleSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BranchProtectionRuleStatus) DeepCopyInto(out *BranchProtectionRuleStatus) {
	*out = *in
	in.ResourceStatus.DeepCopyInto(&out.ResourceStatus)
	out.AtProvider = in.AtProvider
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BranchProtectionRuleStatus.
func (in *BranchProtectionRuleStatus) DeepCopy() *BranchProtectionRuleStatus {
	if in == nil {
		return nil
	}
	out := new(BranchProtectionRuleStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Content) DeepCopyInto(out *Content) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Content.
func (in *Content) DeepCopy() *Content {
	if in == nil {
		return nil
	}
	out := new(Content)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Content) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ContentList) DeepCopyInto(out *ContentList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Content, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ContentList.
func (in *ContentList) DeepCopy() *ContentList {
	if in == nil {
		return nil
	}
	out := new(ContentList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ContentList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ContentObservation) DeepCopyInto(out *ContentObservation) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ContentObservation.
func (in *ContentObservation) DeepCopy() *ContentObservation {
	if in == nil {
		return nil
	}
	out := new(ContentObservation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ContentParameters) DeepCopyInto(out *ContentParameters) {
	*out = *in
	if in.Repository != nil {
		in, out := &in.Repository, &out.Repository
		*out = new(string)
		**out = **in
	}
	if in.RepositoryRef != nil {
		in, out := &in.RepositoryRef, &out.RepositoryRef
		*out = new(v1.Reference)
		(*in).DeepCopyInto(*out)
	}
	if in.RepositorySelector != nil {
		in, out := &in.RepositorySelector, &out.RepositorySelector
		*out = new(v1.Selector)
		(*in).DeepCopyInto(*out)
	}
	if in.Branch != nil {
		in, out := &in.Branch, &out.Branch
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ContentParameters.
func (in *ContentParameters) DeepCopy() *ContentParameters {
	if in == nil {
		return nil
	}
	out := new(ContentParameters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ContentSpec) DeepCopyInto(out *ContentSpec) {
	*out = *in
	in.ResourceSpec.DeepCopyInto(&out.ResourceSpec)
	if in.Reconcile != nil {
		in, out := &in.Reconcile, &out.Reconcile
		*out = new(string)
		**out = **in
	}
	in.ForProvider.DeepCopyInto(&out.ForProvider)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ContentSpec.
func (in *ContentSpec) DeepCopy() *ContentSpec {
	if in == nil {
		return nil
	}
	out := new(ContentSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ContentStatus) DeepCopyInto(out *ContentStatus) {
	*out = *in
	in.ResourceStatus.DeepCopyInto(&out.ResourceStatus)
	out.AtProvider = in.AtProvider
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ContentStatus.
func (in *ContentStatus) DeepCopy() *ContentStatus {
	if in == nil {
		return nil
	}
	out := new(ContentStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Repository) DeepCopyInto(out *Repository) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Repository.
func (in *Repository) DeepCopy() *Repository {
	if in == nil {
		return nil
	}
	out := new(Repository)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Repository) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RepositoryList) DeepCopyInto(out *RepositoryList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Repository, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RepositoryList.
func (in *RepositoryList) DeepCopy() *RepositoryList {
	if in == nil {
		return nil
	}
	out := new(RepositoryList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *RepositoryList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RepositoryObservation) DeepCopyInto(out *RepositoryObservation) {
	*out = *in
	if in.CreatedAt != nil {
		in, out := &in.CreatedAt, &out.CreatedAt
		*out = (*in).DeepCopy()
	}
	if in.PushedAt != nil {
		in, out := &in.PushedAt, &out.PushedAt
		*out = (*in).DeepCopy()
	}
	if in.UpdatedAt != nil {
		in, out := &in.UpdatedAt, &out.UpdatedAt
		*out = (*in).DeepCopy()
	}
	if in.Topics != nil {
		in, out := &in.Topics, &out.Topics
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Permissions != nil {
		in, out := &in.Permissions, &out.Permissions
		*out = make(map[string]bool, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RepositoryObservation.
func (in *RepositoryObservation) DeepCopy() *RepositoryObservation {
	if in == nil {
		return nil
	}
	out := new(RepositoryObservation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RepositoryParameters) DeepCopyInto(out *RepositoryParameters) {
	*out = *in
	if in.Organization != nil {
		in, out := &in.Organization, &out.Organization
		*out = new(string)
		**out = **in
	}
	if in.Description != nil {
		in, out := &in.Description, &out.Description
		*out = new(string)
		**out = **in
	}
	if in.Homepage != nil {
		in, out := &in.Homepage, &out.Homepage
		*out = new(string)
		**out = **in
	}
	if in.Private != nil {
		in, out := &in.Private, &out.Private
		*out = new(bool)
		**out = **in
	}
	if in.Visibility != nil {
		in, out := &in.Visibility, &out.Visibility
		*out = new(string)
		**out = **in
	}
	if in.HasIssues != nil {
		in, out := &in.HasIssues, &out.HasIssues
		*out = new(bool)
		**out = **in
	}
	if in.HasProjects != nil {
		in, out := &in.HasProjects, &out.HasProjects
		*out = new(bool)
		**out = **in
	}
	if in.HasWiki != nil {
		in, out := &in.HasWiki, &out.HasWiki
		*out = new(bool)
		**out = **in
	}
	if in.IsTemplate != nil {
		in, out := &in.IsTemplate, &out.IsTemplate
		*out = new(bool)
		**out = **in
	}
	if in.TeamID != nil {
		in, out := &in.TeamID, &out.TeamID
		*out = new(int64)
		**out = **in
	}
	if in.AutoInit != nil {
		in, out := &in.AutoInit, &out.AutoInit
		*out = new(bool)
		**out = **in
	}
	if in.GitignoreTemplate != nil {
		in, out := &in.GitignoreTemplate, &out.GitignoreTemplate
		*out = new(string)
		**out = **in
	}
	if in.LicenseTemplate != nil {
		in, out := &in.LicenseTemplate, &out.LicenseTemplate
		*out = new(string)
		**out = **in
	}
	if in.AllowSquashMerge != nil {
		in, out := &in.AllowSquashMerge, &out.AllowSquashMerge
		*out = new(bool)
		**out = **in
	}
	if in.AllowMergeCommit != nil {
		in, out := &in.AllowMergeCommit, &out.AllowMergeCommit
		*out = new(bool)
		**out = **in
	}
	if in.AllowRebaseMerge != nil {
		in, out := &in.AllowRebaseMerge, &out.AllowRebaseMerge
		*out = new(bool)
		**out = **in
	}
	if in.DeleteBranchOnMerge != nil {
		in, out := &in.DeleteBranchOnMerge, &out.DeleteBranchOnMerge
		*out = new(bool)
		**out = **in
	}
	if in.HasPages != nil {
		in, out := &in.HasPages, &out.HasPages
		*out = new(bool)
		**out = **in
	}
	if in.HasDownloads != nil {
		in, out := &in.HasDownloads, &out.HasDownloads
		*out = new(bool)
		**out = **in
	}
	if in.DefaultBranch != nil {
		in, out := &in.DefaultBranch, &out.DefaultBranch
		*out = new(string)
		**out = **in
	}
	if in.Archived != nil {
		in, out := &in.Archived, &out.Archived
		*out = new(bool)
		**out = **in
	}
	if in.Template != nil {
		in, out := &in.Template, &out.Template
		*out = new(v1.Reference)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RepositoryParameters.
func (in *RepositoryParameters) DeepCopy() *RepositoryParameters {
	if in == nil {
		return nil
	}
	out := new(RepositoryParameters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RepositorySpec) DeepCopyInto(out *RepositorySpec) {
	*out = *in
	in.ResourceSpec.DeepCopyInto(&out.ResourceSpec)
	in.ForProvider.DeepCopyInto(&out.ForProvider)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RepositorySpec.
func (in *RepositorySpec) DeepCopy() *RepositorySpec {
	if in == nil {
		return nil
	}
	out := new(RepositorySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RepositoryStatus) DeepCopyInto(out *RepositoryStatus) {
	*out = *in
	in.ResourceStatus.DeepCopyInto(&out.ResourceStatus)
	in.AtProvider.DeepCopyInto(&out.AtProvider)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RepositoryStatus.
func (in *RepositoryStatus) DeepCopy() *RepositoryStatus {
	if in == nil {
		return nil
	}
	out := new(RepositoryStatus)
	in.DeepCopyInto(out)
	return out
}
