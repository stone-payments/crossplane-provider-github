/*
Copyright 2021 The Crossplane Authors.
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
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/shurcooL/githubv4"

	"github.com/crossplane-contrib/provider-github/apis/repositories/v1alpha1"
)

var (
	fakeRepository          = "sample"
	fakePattern             = "fak*"
	fakeOwner               = "owner"
	fakeUser                = "user"
	fakeTeamSlug            = "team"
	fakeFormattedUser       = "/" + fakeUser
	fakeFormattedTeam       = "/" + fakeOwner + "/" + fakeTeamSlug
	fakeTrue                = true
	fakeFalse               = false
	fakeCount         int32 = 2
	fakeCheck               = "check"
	fakeIDs                 = []string{"fake", "fake2"}
	fakeID                  = "fake3"
)

func params() *v1alpha1.BranchProtectionRuleParameters {
	return &v1alpha1.BranchProtectionRuleParameters{
		Owner:      fakeOwner,
		Pattern:    fakePattern,
		Repository: &fakeRepository,
		BypassForcePushAllowances: []string{
			fakeFormattedTeam,
			fakeFormattedUser,
		},
		BypassPullRequestAllowances: []string{
			fakeFormattedTeam,
			fakeFormattedUser,
		},
		DismissesStaleReviews: &fakeTrue,
		IsAdminEnforced:       &fakeTrue,
		PushAllowances: []string{
			fakeFormattedTeam,
			fakeFormattedUser,
		},
		RequiredApprovingReviewCount: &fakeCount,
		RequiresCodeOwnerReviews:     &fakeTrue,
		RequiredStatusCheckContexts: []string{
			fakeCheck,
		},
		RequiresCommitSignatures:       &fakeTrue,
		RequiresConversationResolution: &fakeTrue,
		RequiresLinearHistory:          &fakeTrue,
		RequiresStrictStatusChecks:     &fakeTrue,
	}
}

func createInput() githubv4.CreateBranchProtectionRuleInput {
	return githubv4.CreateBranchProtectionRuleInput{
		Pattern:                      githubv4.String(fakePattern),
		RepositoryID:                 githubv4.NewID(githubv4.ID(fakeID)),
		BypassForcePushActorIDs:      githubv4NewIDSlice(githubv4IDSliceEmpty(fakeIDs)),
		BypassPullRequestActorIDs:    githubv4NewIDSlice(githubv4IDSliceEmpty(fakeIDs)),
		DismissesStaleReviews:        (*githubv4.Boolean)(&fakeTrue),
		IsAdminEnforced:              (*githubv4.Boolean)(&fakeTrue),
		PushActorIDs:                 githubv4NewIDSlice(githubv4IDSliceEmpty(fakeIDs)),
		RequiredApprovingReviewCount: (*githubv4.Int)(&fakeCount),
		RequiresCodeOwnerReviews:     (*githubv4.Boolean)(&fakeTrue),
		RequiredStatusCheckContexts: &[]githubv4.String{
			(githubv4.String)(fakeCheck),
		},
		RequiresCommitSignatures:       (*githubv4.Boolean)(&fakeTrue),
		RequiresConversationResolution: (*githubv4.Boolean)(&fakeTrue),
		RequiresLinearHistory:          (*githubv4.Boolean)(&fakeTrue),
		RequiresStrictStatusChecks:     (*githubv4.Boolean)(&fakeTrue),
		RequiresApprovingReviews:       (*githubv4.Boolean)(&fakeTrue),
		RequiresStatusChecks:           (*githubv4.Boolean)(&fakeTrue),
		RestrictsPushes:                (*githubv4.Boolean)(&fakeTrue),
	}
}

func updateInput() githubv4.UpdateBranchProtectionRuleInput {
	return githubv4.UpdateBranchProtectionRuleInput{
		Pattern:                      githubv4.NewString(githubv4.String(fakePattern)),
		BranchProtectionRuleID:       fakeID,
		BypassForcePushActorIDs:      githubv4NewIDSlice(githubv4IDSliceEmpty(fakeIDs)),
		BypassPullRequestActorIDs:    githubv4NewIDSlice(githubv4IDSliceEmpty(fakeIDs)),
		DismissesStaleReviews:        (*githubv4.Boolean)(&fakeTrue),
		IsAdminEnforced:              (*githubv4.Boolean)(&fakeTrue),
		PushActorIDs:                 githubv4NewIDSlice(githubv4IDSliceEmpty(fakeIDs)),
		RequiredApprovingReviewCount: (*githubv4.Int)(&fakeCount),
		RequiresCodeOwnerReviews:     (*githubv4.Boolean)(&fakeTrue),
		RequiredStatusCheckContexts: &[]githubv4.String{
			(githubv4.String)(fakeCheck),
		},
		RequiresCommitSignatures:       (*githubv4.Boolean)(&fakeTrue),
		RequiresConversationResolution: (*githubv4.Boolean)(&fakeTrue),
		RequiresLinearHistory:          (*githubv4.Boolean)(&fakeTrue),
		RequiresStrictStatusChecks:     (*githubv4.Boolean)(&fakeTrue),
		RequiresApprovingReviews:       (*githubv4.Boolean)(&fakeTrue),
		RequiresStatusChecks:           (*githubv4.Boolean)(&fakeTrue),
		RestrictsPushes:                (*githubv4.Boolean)(&fakeTrue),
	}
}

func unsyncedBranchProtection() ExternalBranchProtectionRule {
	return ExternalBranchProtectionRule{
		RequiresCommitSignatures:       (*githubv4.Boolean)(&fakeFalse),
		RequiresConversationResolution: (*githubv4.Boolean)(&fakeFalse),
		RequiresLinearHistory:          (*githubv4.Boolean)(&fakeFalse),
		RequiresStrictStatusChecks:     (*githubv4.Boolean)(&fakeFalse),
	}
}

func syncedBranchProtection() ExternalBranchProtectionRule {
	return ExternalBranchProtectionRule{
		Pattern: (*githubv4.String)(&fakePattern),
		Repository: struct{ Name *githubv4.String }{
			Name: (*githubv4.String)(&fakeRepository),
		},
		BypassForcePushAllowances: struct{ Nodes []ActorTypes }{
			Nodes: []ActorTypes{
				{
					Actor: struct {
						Team Team "graphql:\"... on Team\""
						User User "graphql:\"... on User\""
					}{
						Team: Team{
							Slug: githubv4.String(fakeTeamSlug),
							Type: "Team",
						},
						User: User{},
					},
				},
				{
					Actor: struct {
						Team Team "graphql:\"... on Team\""
						User User "graphql:\"... on User\""
					}{
						User: User{
							Login: githubv4.String(fakeUser),
							Type:  "User",
						},
						Team: Team{},
					},
				},
			},
		},
		BypassPullRequestAllowances: struct{ Nodes []ActorTypes }{
			Nodes: []ActorTypes{
				{
					Actor: struct {
						Team Team "graphql:\"... on Team\""
						User User "graphql:\"... on User\""
					}{
						Team: Team{
							Slug: githubv4.String(fakeTeamSlug),
							Type: "Team",
						},
						User: User{},
					},
				},
				{
					Actor: struct {
						Team Team "graphql:\"... on Team\""
						User User "graphql:\"... on User\""
					}{
						User: User{
							Login: githubv4.String(fakeUser),
							Type:  "User",
						},
						Team: Team{},
					},
				},
			},
		},
		DismissesStaleReviews: (*githubv4.Boolean)(&fakeTrue),
		IsAdminEnforced:       (*githubv4.Boolean)(&fakeTrue),
		PushAllowances: struct{ Nodes []ActorTypes }{
			Nodes: []ActorTypes{
				{
					Actor: struct {
						Team Team "graphql:\"... on Team\""
						User User "graphql:\"... on User\""
					}{
						Team: Team{
							Slug: githubv4.String(fakeTeamSlug),
							Type: "Team",
						},
						User: User{},
					},
				},
				{
					Actor: struct {
						Team Team "graphql:\"... on Team\""
						User User "graphql:\"... on User\""
					}{
						User: User{
							Login: githubv4.String(fakeUser),
							Type:  "User",
						},
						Team: Team{},
					},
				},
			},
		},
		RequiredApprovingReviewCount: (*githubv4.Int)(&fakeCount),
		RequiresCodeOwnerReviews:     (*githubv4.Boolean)(&fakeTrue),
		RequiredStatusCheckContexts: []githubv4.String{
			githubv4.String(fakeCheck),
		},
		RequiresCommitSignatures:       (*githubv4.Boolean)(&fakeTrue),
		RequiresConversationResolution: (*githubv4.Boolean)(&fakeTrue),
		RequiresLinearHistory:          (*githubv4.Boolean)(&fakeTrue),
		RequiresStrictStatusChecks:     (*githubv4.Boolean)(&fakeTrue),
	}
}

func TestLateInitialize(t *testing.T) {
	type args struct {
		external ExternalBranchProtectionRule
		params   *v1alpha1.BranchProtectionRuleParameters
	}
	cases := map[string]struct {
		args
		out *v1alpha1.BranchProtectionRuleParameters
	}{
		"Must initialize empty Parameters fields if they are given in External": {
			args: args{
				params: &v1alpha1.BranchProtectionRuleParameters{
					Repository: &fakeRepository,
					Pattern:    fakePattern,
					Owner:      fakeOwner,
				},
				external: syncedBranchProtection(),
			},
			out: params(),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			LateInitialize(tc.args.params, tc.args.external)
			if diff := cmp.Diff(tc.out, tc.args.params); diff != "" {
				t.Errorf("LateInitialize(...): -want, +got:\n%s", diff)
			}
		})
	}
}

func TestIsUpToDate(t *testing.T) {
	type args struct {
		external ExternalBranchProtectionRule
		params   *v1alpha1.BranchProtectionRuleParameters
	}
	cases := map[string]struct {
		args
		upToDate bool
	}{
		"NotUpToDate": {
			args: args{
				external: unsyncedBranchProtection(),
				params:   params(),
			},
			upToDate: false,
		},
		"UpToDate": {
			args: args{
				external: syncedBranchProtection(),
				params:   params(),
			},
			upToDate: true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got, _ := IsUpToDate(*tc.args.params, tc.args.external)
			if diff := cmp.Diff(tc.upToDate, got == ""); diff != "" {
				t.Errorf("IsUpToDate(...): -want, +got:\n%s", diff)
			}
		})
	}
}

func TestGenerateParametersFromExternal(t *testing.T) {
	type args struct {
		external ExternalBranchProtectionRule
		params   *v1alpha1.BranchProtectionRuleParameters
	}
	cases := map[string]struct {
		args
		want v1alpha1.BranchProtectionRuleParameters
	}{
		"OverrideAllFields": {
			args: args{
				external: syncedBranchProtection(),
				params: &v1alpha1.BranchProtectionRuleParameters{
					Owner:      fakeOwner,
					Pattern:    fakePattern,
					Repository: &fakeRepository,
				},
			},
			want: *params(),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := GenerateParametersFromExternal(*tc.args.params, tc.args.external)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("GenerateParametersFromExternal(...): -want, +got:\n%s", diff)
			}
		})
	}
}

func TestGenerateCreateBranchProtectionRuleInput(t *testing.T) {
	type args struct {
		repositoryID string
		ids          []string
		params       *v1alpha1.BranchProtectionRuleParameters
	}
	cases := map[string]struct {
		args
		want githubv4.CreateBranchProtectionRuleInput
	}{
		"GenerateInput": {
			args: args{
				params:       params(),
				repositoryID: fakeID,
				ids:          fakeIDs,
			},
			want: createInput(),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := GenerateCreateBranchProtectionRuleInput(*tc.args.params, tc.args.ids, tc.args.ids, tc.args.ids, tc.args.repositoryID)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("GenerateCreateBranchProtectionRuleInput(...): -want, +got:\n%s", diff)
			}
		})
	}
}

func TestGenerateUpdateBranchProtectionRuleInput(t *testing.T) {
	type args struct {
		repositoryID string
		ids          []string
		params       *v1alpha1.BranchProtectionRuleParameters
	}
	cases := map[string]struct {
		args
		want githubv4.UpdateBranchProtectionRuleInput
	}{
		"GenerateInput": {
			args: args{
				params:       params(),
				repositoryID: fakeID,
				ids:          fakeIDs,
			},
			want: updateInput(),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := GenerateUpdateBranchProtectionRuleInput(*tc.args.params, tc.args.ids, tc.args.ids, tc.args.ids, tc.args.repositoryID)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("GenerateUpdateBranchProtectionRuleInput(...): -want, +got:\n%s", diff)
			}
		})
	}
}

func TestIsTeamActor(t *testing.T) {
	type args struct {
		slug string
	}
	cases := map[string]struct {
		args
		want bool
	}{
		"IsNotTeam": {
			args: args{
				slug: "/user",
			},
			want: false,
		},
		"IsTeam": {
			args: args{
				slug: "/org/slug",
			},
			want: true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := IsTeamActor(tc.args.slug)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("IsTeamActor(...): -want, +got:\n%s", diff)
			}
		})
	}
}
