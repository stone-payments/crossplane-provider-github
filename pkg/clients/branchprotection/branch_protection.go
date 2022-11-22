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
	"sort"
	"strings"

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

// ActorTypes represents the possible
// types of an actor
type ActorTypes struct {
	Actor struct {
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

	// Sort strings to be compared
	sort.Strings(current.BypassForcePushAllowances)
	sort.Strings(current.BypassPullRequestAllowances)
	sort.Strings(current.PushAllowances)
	sort.Strings(desired.BypassForcePushAllowances)
	sort.Strings(desired.BypassPullRequestAllowances)
	sort.Strings(desired.PushAllowances)

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

		if v.Actor.User.Login != "" && v.Actor.User.Type == "User" {
			actor := fmt.Sprintf("/%v", string(v.Actor.User.Login))
			list = append(list, actor)
		}
	}
	return list
}

// GenerateCreateBranchProtectionRuleInput generates a githubv4.CreateBranchProtectionRuleInput
// based on the v1alpha1.BranchProtectionRuleParameters passed as parameter
func GenerateCreateBranchProtectionRuleInput(params v1alpha1.BranchProtectionRuleParameters, bypassForcePushIds, bypassPullRequestIds, pushIds []string, repositoryID string) githubv4.CreateBranchProtectionRuleInput { // nolint:gocyclo
	input := githubv4.CreateBranchProtectionRuleInput{
		Pattern:      githubv4.String(params.Pattern),
		RepositoryID: githubv4.NewID(githubv4.ID(repositoryID)),
		// Setting RequiresStatusChecks without defining the properties below it (RequiredStatusCheckContexts or RequiresStrictStatusChecks)
		// has no effect in the branch protection behavior. We need to have it set to true because we can't modify the child properties
		// with it disabled -- when it is disabled, it causes an update loop in the managed resource.
		RequiresStatusChecks: githubv4.NewBoolean(true),
	}
	var restrictsPushes, requiresApprovingReviews bool

	input.BypassForcePushActorIDs = githubv4NewIDSlice(githubv4IDSliceEmpty(bypassForcePushIds))
	input.BypassPullRequestActorIDs = githubv4NewIDSlice(githubv4IDSliceEmpty(bypassPullRequestIds))
	input.PushActorIDs = githubv4NewIDSlice(githubv4IDSliceEmpty(pushIds))

	restrictsPushes = restrictsPushes || len(pushIds) > 0
	requiresApprovingReviews = requiresApprovingReviews || len(bypassPullRequestIds) > 0

	if params.IsAdminEnforced != nil {
		input.IsAdminEnforced = (*githubv4.Boolean)(params.IsAdminEnforced)
	}

	if params.DismissesStaleReviews != nil {
		input.DismissesStaleReviews = (*githubv4.Boolean)(params.DismissesStaleReviews)
		requiresApprovingReviews = requiresApprovingReviews || *params.DismissesStaleReviews
	}

	if params.RequiredApprovingReviewCount != nil {
		input.RequiredApprovingReviewCount = (*githubv4.Int)(params.RequiredApprovingReviewCount)
		requiresApprovingReviews = requiresApprovingReviews || *params.RequiredApprovingReviewCount != 0
	}

	if params.RequiredStatusCheckContexts != nil {
		input.RequiredStatusCheckContexts = githubv4NewStringSlice(githubv4StringSlice(params.RequiredStatusCheckContexts))
	}

	if params.RequiresCodeOwnerReviews != nil {
		input.RequiresCodeOwnerReviews = (*githubv4.Boolean)(params.RequiresCodeOwnerReviews)
		requiresApprovingReviews = requiresApprovingReviews || *params.RequiresCodeOwnerReviews
	}

	if params.RequiresCommitSignatures != nil {
		input.RequiresCommitSignatures = (*githubv4.Boolean)(params.RequiresCommitSignatures)
	}

	if params.RequiresConversationResolution != nil {
		input.RequiresConversationResolution = (*githubv4.Boolean)(params.RequiresConversationResolution)
	}

	if params.RequiresLinearHistory != nil {
		input.RequiresLinearHistory = (*githubv4.Boolean)(params.RequiresLinearHistory)
	}

	if params.RequiresStrictStatusChecks != nil {
		input.RequiresStrictStatusChecks = (*githubv4.Boolean)(params.RequiresStrictStatusChecks)
	}

	input.RequiresApprovingReviews = (*githubv4.Boolean)(&requiresApprovingReviews)
	input.RestrictsPushes = (*githubv4.Boolean)(&restrictsPushes)

	return input
}

// GenerateUpdateBranchProtectionRuleInput generates a githubv4.UpdateBranchProtectionRuleInput
// based on the v1alpha1.BranchProtectionRuleParameters passed as parameter
func GenerateUpdateBranchProtectionRuleInput(params v1alpha1.BranchProtectionRuleParameters, bypassForcePushIds, bypassPullRequestIds, pushIds []string, id string) githubv4.UpdateBranchProtectionRuleInput { // nolint:gocyclo
	input := githubv4.UpdateBranchProtectionRuleInput{
		Pattern:                githubv4.NewString(githubv4.String(params.Pattern)),
		BranchProtectionRuleID: id,
		// Setting RequiresStatusChecks without defining the properties below it (RequiredStatusCheckContexts or RequiresStrictStatusChecks)
		// has no effect in the branch protection behavior. We need to have it set to true because we can't modify the child properties
		// with it disabled -- when it is disabled, it causes an update loop in the managed resource.
		RequiresStatusChecks: githubv4.NewBoolean(true),
	}

	var restrictsPushes, requiresApprovingReviews bool

	input.BypassForcePushActorIDs = githubv4NewIDSlice(githubv4IDSliceEmpty(bypassForcePushIds))
	input.BypassPullRequestActorIDs = githubv4NewIDSlice(githubv4IDSliceEmpty(bypassPullRequestIds))
	input.PushActorIDs = githubv4NewIDSlice(githubv4IDSliceEmpty(pushIds))

	restrictsPushes = restrictsPushes || len(pushIds) > 0
	requiresApprovingReviews = requiresApprovingReviews || len(bypassPullRequestIds) > 0

	if params.IsAdminEnforced != nil {
		input.IsAdminEnforced = (*githubv4.Boolean)(params.IsAdminEnforced)
	}

	if params.DismissesStaleReviews != nil {
		input.DismissesStaleReviews = (*githubv4.Boolean)(params.DismissesStaleReviews)
		requiresApprovingReviews = requiresApprovingReviews || *params.DismissesStaleReviews
	}

	if params.RequiredApprovingReviewCount != nil {
		input.RequiredApprovingReviewCount = (*githubv4.Int)(params.RequiredApprovingReviewCount)
		requiresApprovingReviews = requiresApprovingReviews || *params.RequiredApprovingReviewCount != 0
	}

	if params.RequiredStatusCheckContexts != nil {
		input.RequiredStatusCheckContexts = githubv4NewStringSlice(githubv4StringSlice(params.RequiredStatusCheckContexts))
	}

	if params.RequiresCodeOwnerReviews != nil {
		input.RequiresCodeOwnerReviews = (*githubv4.Boolean)(params.RequiresCodeOwnerReviews)
		requiresApprovingReviews = requiresApprovingReviews || *params.RequiresCodeOwnerReviews
	}

	if params.RequiresCommitSignatures != nil {
		input.RequiresCommitSignatures = (*githubv4.Boolean)(params.RequiresCommitSignatures)
	}

	if params.RequiresConversationResolution != nil {
		input.RequiresConversationResolution = (*githubv4.Boolean)(params.RequiresConversationResolution)
	}

	if params.RequiresLinearHistory != nil {
		input.RequiresLinearHistory = (*githubv4.Boolean)(params.RequiresLinearHistory)
	}

	if params.RequiresStrictStatusChecks != nil {
		input.RequiresStrictStatusChecks = (*githubv4.Boolean)(params.RequiresStrictStatusChecks)
	}

	input.RequiresApprovingReviews = (*githubv4.Boolean)(&requiresApprovingReviews)
	input.RestrictsPushes = (*githubv4.Boolean)(&restrictsPushes)

	return input
}

// IsTeamActor returns if the slug passed as parameter
// is from a Team actor
func IsTeamActor(slug string) bool {
	return len(strings.Split(slug, "/")) == 3
}

func githubv4IDSliceEmpty(slice []string) []githubv4.ID {
	ids := make([]githubv4.ID, 0)
	for _, s := range slice {
		ids = append(ids, githubv4.ID(s))
	}
	return ids
}

func githubv4StringSlice(slice []string) []githubv4.String {
	ss := make([]githubv4.String, 0)
	for _, s := range slice {
		ss = append(ss, githubv4.String(s))
	}
	return ss
}

func githubv4NewIDSlice(v []githubv4.ID) *[]githubv4.ID { return &v }

func githubv4NewStringSlice(v []githubv4.String) *[]githubv4.String { return &v }
