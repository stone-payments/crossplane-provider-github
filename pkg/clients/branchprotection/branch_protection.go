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
package branchprotection

import (
	"context"
	"fmt"

	"github.com/crossplane-contrib/provider-github/apis/repositories/v1alpha1"
	ghclient "github.com/crossplane-contrib/provider-github/pkg/clients"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/mitchellh/copystructure"
	"github.com/pkg/errors"
	"github.com/shurcooL/githubv4"
)

const (
	errCheckUpToDate = "unable to determine if external resource is up to date"
)

// Service defines the GraphQL operations
type Service interface {
	Query(ctx context.Context, q interface{}, variables map[string]interface{}) error
	Mutate(ctx context.Context, m interface{}, input githubv4.Input, variables map[string]interface{}) error
}

// Team represents a GitHub Team
type Team struct {
	Slug githubv4.String `graphql:"slug"`
	Type githubv4.String `graphql:"__typename"`
}

// User represents a GitHub User
type User struct {
	Login githubv4.String `graphql:"login"`
	Type  githubv4.String `graphql:"__typename"`
}

// App represents a GitHub App
type App struct {
	Slug githubv4.String `graphql:"slug"`
	Type githubv4.String `graphql:"__typename"`
}

// ActorTypes represents the possible
// types of an actor
type ActorTypes struct {
	Actor struct {
		App  App  `graphql:"... on App"`
		Team Team `graphql:"... on Team"`
		User User `graphql:"... on User"`
	}
}

// ExternalBranchProtectionRule represents the GraphQL
// schema of the BranchProtection
type ExternalBranchProtectionRule struct {
	Repository struct {
		Name *githubv4.String
	}
	PushAllowances struct {
		Nodes []ActorTypes
	} `graphql:"pushAllowances(first: 100)"`
	ReviewDismissalAllowances struct {
		Nodes []ActorTypes
	} `graphql:"reviewDismissalAllowances(first: 100)"`
	BypassPullRequestAllowances struct {
		Nodes []ActorTypes
	} `graphql:"bypassPullRequestAllowances(first: 100)"`
	BypassForcePushAllowances struct {
		Nodes []ActorTypes
	} `graphql:"bypassForcePushAllowances(first: 100)"`
	AllowsDeletions                *githubv4.Boolean
	AllowsForcePushes              *githubv4.Boolean
	DismissesStaleReviews          *githubv4.Boolean
	ID                             *githubv4.ID
	IsAdminEnforced                *githubv4.Boolean
	Pattern                        *githubv4.String
	RequiredApprovingReviewCount   *githubv4.Int
	RequiredStatusCheckContexts    []githubv4.String
	RequiresApprovingReviews       *githubv4.Boolean
	RequiresCodeOwnerReviews       *githubv4.Boolean
	RequiresCommitSignatures       *githubv4.Boolean
	RequiresLinearHistory          *githubv4.Boolean
	RequiresConversationResolution *githubv4.Boolean
	RequiresStatusChecks           *githubv4.Boolean
	RequiresStrictStatusChecks     *githubv4.Boolean
	RestrictsPushes                *githubv4.Boolean
	RestrictsReviewDismissals      *githubv4.Boolean
}

// NewClient creates a new *githubv4.Client
func NewClient(token string) (Service, error) {
	c, err := ghclient.NewV4Client(token)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// LateInitialize fills the empty values in the parameters struct
// if they are defined in the external struct
func LateInitialize(params *v1alpha1.BranchProtectionRuleParameters, external ExternalBranchProtectionRule) { // nolint:gocyclo
	if params.BypassForcePushAllowances == nil && external.BypassForcePushAllowances.Nodes != nil {
		params.BypassForcePushAllowances = transformActorTypesToSlice(
			external.BypassForcePushAllowances.Nodes,
			params.Owner,
		)
	}

	if params.BypassPullRequestAllowances == nil && external.BypassPullRequestAllowances.Nodes != nil {
		params.BypassPullRequestAllowances = transformActorTypesToSlice(
			external.BypassPullRequestAllowances.Nodes,
			params.Owner,
		)
	}

	if params.PushAllowances == nil && external.PushAllowances.Nodes != nil {
		params.PushAllowances = transformActorTypesToSlice(
			external.PushAllowances.Nodes,
			params.Owner,
		)
	}

	if params.RequiredStatusCheckContexts == nil && external.RequiredStatusCheckContexts != nil {
		for _, v := range external.RequiredStatusCheckContexts {
			params.RequiredStatusCheckContexts = append(params.RequiredStatusCheckContexts, string(v))
		}
	}

	if params.DismissesStaleReviews == nil && external.DismissesStaleReviews != nil {
		params.DismissesStaleReviews = (*bool)(external.DismissesStaleReviews)
	}

	if params.IsAdminEnforced == nil && external.IsAdminEnforced != nil {
		params.IsAdminEnforced = (*bool)(external.IsAdminEnforced)
	}

	if params.RequiredApprovingReviewCount == nil && external.RequiredApprovingReviewCount != nil {
		params.RequiredApprovingReviewCount = (*int32)(external.RequiredApprovingReviewCount)
	}

	if params.RequiresCodeOwnerReviews == nil && external.RequiresCodeOwnerReviews != nil {
		params.RequiresCodeOwnerReviews = (*bool)(external.RequiresCodeOwnerReviews)
	}

	if params.RequiresCommitSignatures == nil && external.RequiresCommitSignatures != nil {
		params.RequiresCommitSignatures = (*bool)(external.RequiresCommitSignatures)
	}

	if params.RequiresConversationResolution == nil && external.RequiresConversationResolution != nil {
		params.RequiresConversationResolution = (*bool)(external.RequiresConversationResolution)
	}

	if params.RequiresLinearHistory == nil && external.RequiresLinearHistory != nil {
		params.RequiresLinearHistory = (*bool)(external.RequiresLinearHistory)
	}

	if params.RequiresStrictStatusChecks == nil && external.RequiresStrictStatusChecks != nil {
		params.RequiresStrictStatusChecks = (*bool)(external.RequiresStrictStatusChecks)
	}
}

// IsUpToDate checks whether the desired state is the current state
func IsUpToDate(desired v1alpha1.BranchProtectionRuleParameters, external ExternalBranchProtectionRule) (string, error) {
	copy, err := copystructure.Copy(desired)
	if err != nil {
		return "", errors.Wrap(err, errCheckUpToDate)
	}
	clone, ok := copy.(v1alpha1.BranchProtectionRuleParameters)
	if !ok {
		return "", errors.New(errCheckUpToDate)
	}

	current := GenerateParametersFromExternal(clone, external)

	return cmp.Diff(
		desired,
		current,
		cmpopts.EquateEmpty(),
	), nil
}

// GenerateParametersFromExternal overrides the parameters field values based
// on the external ones.
func GenerateParametersFromExternal(params v1alpha1.BranchProtectionRuleParameters, external ExternalBranchProtectionRule) v1alpha1.BranchProtectionRuleParameters { // nolint:gocyclo
	if external.BypassForcePushAllowances.Nodes != nil {
		params.BypassForcePushAllowances = transformActorTypesToSlice(
			external.BypassForcePushAllowances.Nodes,
			params.Owner,
		)
	}

	if external.BypassPullRequestAllowances.Nodes != nil {
		params.BypassPullRequestAllowances = transformActorTypesToSlice(
			external.BypassPullRequestAllowances.Nodes,
			params.Owner,
		)
	}

	if external.DismissesStaleReviews != nil {
		params.DismissesStaleReviews = (*bool)(external.DismissesStaleReviews)
	}

	if external.IsAdminEnforced != nil {
		params.IsAdminEnforced = (*bool)(external.IsAdminEnforced)
	}

	if external.PushAllowances.Nodes != nil {
		params.PushAllowances = transformActorTypesToSlice(
			external.PushAllowances.Nodes,
			params.Owner,
		)
	}

	if external.RequiredApprovingReviewCount != nil {
		params.RequiredApprovingReviewCount = (*int32)(external.RequiredApprovingReviewCount)
	}

	if external.RequiredStatusCheckContexts != nil {
		statusCheckContexts := make([]string, 0)
		for _, v := range external.RequiredStatusCheckContexts {
			statusCheckContexts = append(statusCheckContexts, (string)(v))
		}
		params.RequiredStatusCheckContexts = statusCheckContexts
	}

	if external.RequiresCodeOwnerReviews != nil {
		params.RequiresCodeOwnerReviews = (*bool)(external.RequiresCodeOwnerReviews)
	}

	if external.RequiresCommitSignatures != nil {
		params.RequiresCommitSignatures = (*bool)(external.RequiresCommitSignatures)
	}

	if external.RequiresConversationResolution != nil {
		params.RequiresConversationResolution = (*bool)(external.RequiresConversationResolution)
	}

	if external.RequiresLinearHistory != nil {
		params.RequiresLinearHistory = (*bool)(external.RequiresLinearHistory)
	}

	if external.RequiresStrictStatusChecks != nil {
		params.RequiresStrictStatusChecks = (*bool)(external.RequiresStrictStatusChecks)
	}

	return params
}

func transformActorTypesToSlice(actors []ActorTypes, org string) []string {
	list := make([]string, 0)
	for _, v := range actors {
		if v.Actor.Team.Slug != "" && v.Actor.Team.Type == "Team" {
			actor := fmt.Sprintf("/%v/%v", org, string(v.Actor.Team.Slug))
			list = append(list, actor)
		}

		if v.Actor.User.Login != "" && v.Actor.Team.Type == "User" {
			actor := fmt.Sprintf("/%v", string(v.Actor.User.Login))
			list = append(list, actor)
		}

		if v.Actor.App.Slug != "" && v.Actor.Team.Type == "App" {
			actor := fmt.Sprintf("/app/%v", string(v.Actor.App.Slug))
			list = append(list, actor)
		}
	}
	return list
}
