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
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/shurcooL/githubv4"

	"github.com/crossplane-contrib/provider-github/apis/repositories/v1alpha1"
	"github.com/crossplane-contrib/provider-github/pkg/clients/branchprotection"
	"github.com/crossplane-contrib/provider-github/pkg/clients/branchprotection/fake"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/crossplane/crossplane-runtime/pkg/test"
	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	unexpectedObject   resource.Managed
	errBoom            = errors.New("boom")
	fakeID             = "id"
	fakeRepository     = "fake"
	fakePattern        = "pattern"
	fakeFalse          = false
	fakeTrue           = true
	errUnexpectedQuery = "unexpected query"
	fakeActors         = []string{"/fake1", "/fake/fake2"}
)

type protectionOption func(*v1alpha1.BranchProtectionRule)

func newBranchProtectionRule(opts ...protectionOption) *v1alpha1.BranchProtectionRule {
	bpr := &v1alpha1.BranchProtectionRule{}

	for _, f := range opts {
		f(bpr)
	}
	return bpr
}

func withActors(actors []string) protectionOption {
	return func(bpr *v1alpha1.BranchProtectionRule) {
		bpr.Spec.ForProvider.PushAllowances = actors
		bpr.Spec.ForProvider.BypassForcePushAllowances = actors
		bpr.Spec.ForProvider.BypassPullRequestAllowances = actors
	}
}

func withPattern(pattern string) protectionOption {
	return func(bpr *v1alpha1.BranchProtectionRule) { bpr.Spec.ForProvider.Pattern = pattern }
}

func withRepository(repository string) protectionOption {
	return func(bpr *v1alpha1.BranchProtectionRule) { bpr.Spec.ForProvider.Repository = &repository }
}

func withCommitSignatures(value bool) protectionOption {
	return func(bpr *v1alpha1.BranchProtectionRule) { bpr.Spec.ForProvider.RequiresCommitSignatures = &value }
}

func newExternalBranchProtectionRule() branchprotection.ExternalBranchProtectionRule {
	return branchprotection.ExternalBranchProtectionRule{
		Pattern: (*githubv4.String)(&fakePattern),
		Repository: struct{ Name *githubv4.String }{
			Name: (*githubv4.String)(&fakeRepository),
		},
		RequiresCommitSignatures: (*githubv4.Boolean)(&fakeFalse),
	}
}

type args struct {
	kube   client.Client
	mg     resource.Managed
	github branchprotection.Service
}

func TestObserve(t *testing.T) {
	type want struct {
		eo  managed.ExternalObservation
		err error
	}

	cases := map[string]struct {
		reason string
		args   args
		want   want
	}{
		"ResourceIsNotBranchProtectionRule": {
			reason: "Must return an error resource is not BranchProtectionRule",
			args: args{
				mg: unexpectedObject,
			},
			want: want{
				eo:  managed.ExternalObservation{},
				err: errors.New(errUnexpectedObject),
			},
		},
		"CannotGetBranchProtectionRule": {
			reason: "Must return an error if BranchProtectionRule query fails",
			args: args{
				mg: newBranchProtectionRule(
					withRepository(fakeRepository),
				),
				github: &fake.MockServiceBranchProtection{
					MockQuery: func(ctx context.Context, q interface{}, variables map[string]interface{}) error {
						return errBoom
					},
				},
			},
			want: want{
				eo:  managed.ExternalObservation{},
				err: errors.Wrap(errBoom, errCheckBranchProtectionRule),
			},
		},
		"CannotFindBranchProtectionRule": {
			reason: "Must return resource doesn't exists if the BranchProtectionRule cannot be found",
			args: args{
				mg: newBranchProtectionRule(
					withPattern("fake"),
					withRepository(fakeRepository),
				),
				github: &fake.MockServiceBranchProtection{
					MockQuery: func(ctx context.Context, q interface{}, variables map[string]interface{}) error {
						query, ok := q.(*struct {
							Repository struct {
								ID                    githubv4.String "graphql:\"id\""
								BranchProtectionRules struct {
									Nodes []struct {
										Pattern githubv4.String "graphql:\"pattern\""
										ID      githubv4.String "graphql:\"id\""
									} "graphql:\"nodes\""
								} "graphql:\"branchProtectionRules(first: 100)\""
							} "graphql:\"repository(owner: $owner, name: $name)\""
						})
						if !ok {
							return errors.New(errUnexpectedQuery)
						}

						query.Repository.ID = githubv4.String(fakeID)
						query.Repository.BranchProtectionRules.Nodes = []struct {
							Pattern githubv4.String "graphql:\"pattern\""
							ID      githubv4.String "graphql:\"id\""
						}{
							{
								Pattern: "test",
							},
						}
						return nil
					},
				},
			},
			want: want{
				eo:  managed.ExternalObservation{ResourceExists: false},
				err: nil,
			},
		},
		"BranchProtectionRuleIsNotUpToDate": {
			reason: "Must return ResourceUpToDate as false if BranchProtectionRule is outdated",
			args: args{
				mg: newBranchProtectionRule(
					withRepository(fakeRepository),
					withPattern(fakePattern),
					withCommitSignatures(fakeTrue),
				),
				github: &fake.MockServiceBranchProtection{
					MockQuery: func(ctx context.Context, q interface{}, variables map[string]interface{}) error {
						checkQuery, ok := q.(*struct {
							Repository struct {
								ID                    githubv4.String "graphql:\"id\""
								BranchProtectionRules struct {
									Nodes []struct {
										Pattern githubv4.String "graphql:\"pattern\""
										ID      githubv4.String "graphql:\"id\""
									} "graphql:\"nodes\""
								} "graphql:\"branchProtectionRules(first: 100)\""
							} "graphql:\"repository(owner: $owner, name: $name)\""
						})
						if ok {
							checkQuery.Repository.ID = githubv4.String(fakeID)
							checkQuery.Repository.BranchProtectionRules.Nodes = []struct {
								Pattern githubv4.String "graphql:\"pattern\""
								ID      githubv4.String "graphql:\"id\""
							}{
								{
									Pattern: githubv4.String(fakePattern),
								},
							}

							return nil
						}

						getQuery, ok := q.(*struct {
							Node struct {
								Node branchprotection.ExternalBranchProtectionRule "graphql:\"... on BranchProtectionRule\""
							} "graphql:\"node(id: $id)\""
						})
						if !ok {
							return errors.New(errUnexpectedQuery)
						}

						getQuery.Node.Node = newExternalBranchProtectionRule()

						return nil
					},
				},
			},
			want: want{
				eo: managed.ExternalObservation{
					ResourceExists:   true,
					ResourceUpToDate: false,
				},
				err: nil,
			},
		},
		"BranchProtectionRuleIsUpToDate": {
			reason: "Must return ResourceUpToDate as true if BranchProtectionRule is up to date",
			args: args{
				mg: newBranchProtectionRule(
					withRepository(fakeRepository),
					withPattern(fakePattern),
					withCommitSignatures(fakeFalse),
				),
				github: &fake.MockServiceBranchProtection{
					MockQuery: func(ctx context.Context, q interface{}, variables map[string]interface{}) error {
						checkQuery, ok := q.(*struct {
							Repository struct {
								ID                    githubv4.String "graphql:\"id\""
								BranchProtectionRules struct {
									Nodes []struct {
										Pattern githubv4.String "graphql:\"pattern\""
										ID      githubv4.String "graphql:\"id\""
									} "graphql:\"nodes\""
								} "graphql:\"branchProtectionRules(first: 100)\""
							} "graphql:\"repository(owner: $owner, name: $name)\""
						})
						if ok {
							checkQuery.Repository.ID = githubv4.String(fakeID)
							checkQuery.Repository.BranchProtectionRules.Nodes = []struct {
								Pattern githubv4.String "graphql:\"pattern\""
								ID      githubv4.String "graphql:\"id\""
							}{
								{
									Pattern: githubv4.String(fakePattern),
								},
							}

							return nil
						}

						getQuery, ok := q.(*struct {
							Node struct {
								Node branchprotection.ExternalBranchProtectionRule "graphql:\"... on BranchProtectionRule\""
							} "graphql:\"node(id: $id)\""
						})
						if !ok {
							return errors.New(errUnexpectedQuery)
						}

						getQuery.Node.Node = newExternalBranchProtectionRule()

						return nil
					},
				},
			},
			want: want{
				eo: managed.ExternalObservation{
					ResourceExists:   true,
					ResourceUpToDate: true,
				},
				err: nil,
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := external{
				client: tc.args.kube,
				gh:     tc.args.github,
			}
			got, err := e.Observe(context.Background(), tc.args.mg)
			if diff := cmp.Diff(tc.want.eo, got, cmpopts.IgnoreFields(managed.ExternalObservation{}, "Diff")); diff != "" {
				t.Errorf("Observe(...): -want, +got:\n%s", diff)
			}
			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("Observe(...): -want error, +got error:\n%s", diff)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	type want struct {
		eo  managed.ExternalCreation
		err error
	}

	cases := map[string]struct {
		reason string
		args   args
		want   want
	}{
		"ResourceIsNotBranchProtectionRule": {
			reason: "Must return an error resource is not BranchProtectionRule",
			args: args{
				mg: unexpectedObject,
			},
			want: want{
				eo:  managed.ExternalCreation{},
				err: errors.New(errUnexpectedObject),
			},
		},
		"MutateFails": {
			reason: "Must return an error if the Mutation fails",
			args: args{
				github: &fake.MockServiceBranchProtection{
					MockMutate: func(ctx context.Context, m interface{}, input githubv4.Input, variables map[string]interface{}) error {
						return errBoom
					},
				},
				mg: newBranchProtectionRule(),
			},
			want: want{
				eo:  managed.ExternalCreation{},
				err: errors.Wrap(errBoom, errCreateBranchProtectionRule),
			},
		},
		"QueryActorsFails": {
			reason: "Must return an error if Actors IDs can't be queried in the API",
			args: args{
				github: &fake.MockServiceBranchProtection{
					MockQuery: func(ctx context.Context, q interface{}, variables map[string]interface{}) error {
						return errBoom
					},
				},
				mg: newBranchProtectionRule(
					withActors(fakeActors),
				),
			},
			want: want{
				eo:  managed.ExternalCreation{},
				err: errors.Wrap(errBoom, errCreateBranchProtectionRule),
			},
		},
		"Success": {
			reason: "Must not return an error if everything goes well",
			args: args{
				github: &fake.MockServiceBranchProtection{
					MockQuery: func(ctx context.Context, q interface{}, variables map[string]interface{}) error {
						queryTeam, ok := q.(*struct {
							Organization struct {
								Team struct{ ID string } "graphql:\"team(slug: $slug)\""
							} "graphql:\"organization(login: $organization)\""
						})
						if ok {
							queryTeam.Organization.Team.ID = fakeID

							return nil
						}

						queryUser, ok := q.(*struct {
							User struct{ ID string } "graphql:\"user(login: $user)\""
						})
						if !ok {
							return errors.New(errUnexpectedQuery)
						}

						queryUser.User.ID = fakeID

						return nil
					},
					MockMutate: func(ctx context.Context, m interface{}, input githubv4.Input, variables map[string]interface{}) error {
						return nil
					},
				},
				mg: newBranchProtectionRule(
					withActors(fakeActors),
				),
			},
			want: want{
				eo:  managed.ExternalCreation{},
				err: nil,
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := external{
				gh: tc.args.github,
			}
			got, err := e.Create(context.Background(), tc.args.mg)
			if diff := cmp.Diff(tc.want.eo, got); diff != "" {
				t.Errorf("Create(...): -want, +got:\n%s", diff)
			}
			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("Create(...): -want error, +got error:\n%s", diff)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	type want struct {
		eo  managed.ExternalUpdate
		err error
	}

	cases := map[string]struct {
		reason string
		args   args
		want   want
	}{
		"ResourceIsNotBranchProtectionRule": {
			reason: "Must return an error resource is not BranchProtectionRule",
			args: args{
				mg: unexpectedObject,
			},
			want: want{
				eo:  managed.ExternalUpdate{},
				err: errors.New(errUnexpectedObject),
			},
		},
		"MutateFails": {
			reason: "Must return an error if the Mutation fails",
			args: args{
				github: &fake.MockServiceBranchProtection{
					MockMutate: func(ctx context.Context, m interface{}, input githubv4.Input, variables map[string]interface{}) error {
						return errBoom
					},
				},
				mg: newBranchProtectionRule(),
			},
			want: want{
				eo:  managed.ExternalUpdate{},
				err: errors.Wrap(errBoom, errUpdateBranchProtectionRule),
			},
		},
		"QueryActorsFails": {
			reason: "Must return an error if Actors IDs can't be queried in the API",
			args: args{
				github: &fake.MockServiceBranchProtection{
					MockQuery: func(ctx context.Context, q interface{}, variables map[string]interface{}) error {
						return errBoom
					},
				},
				mg: newBranchProtectionRule(
					withActors(fakeActors),
				),
			},
			want: want{
				eo:  managed.ExternalUpdate{},
				err: errors.Wrap(errBoom, errUpdateBranchProtectionRule),
			},
		},
		"Success": {
			reason: "Must not return an error if everything goes well",
			args: args{
				github: &fake.MockServiceBranchProtection{
					MockQuery: func(ctx context.Context, q interface{}, variables map[string]interface{}) error {
						queryTeam, ok := q.(*struct {
							Organization struct {
								Team struct{ ID string } "graphql:\"team(slug: $slug)\""
							} "graphql:\"organization(login: $organization)\""
						})
						if ok {
							queryTeam.Organization.Team.ID = fakeID

							return nil
						}

						queryUser, ok := q.(*struct {
							User struct{ ID string } "graphql:\"user(login: $user)\""
						})
						if !ok {
							return errors.New(errUnexpectedQuery)
						}

						queryUser.User.ID = fakeID

						return nil
					},
					MockMutate: func(ctx context.Context, m interface{}, input githubv4.Input, variables map[string]interface{}) error {
						return nil
					},
				},
				mg: newBranchProtectionRule(
					withActors(fakeActors),
				),
			},
			want: want{
				eo:  managed.ExternalUpdate{},
				err: nil,
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := external{
				gh: tc.args.github,
			}
			got, err := e.Update(context.Background(), tc.args.mg)
			if diff := cmp.Diff(tc.want.eo, got); diff != "" {
				t.Errorf("Update(...): -want, +got:\n%s", diff)
			}
			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("Update(...): -want error, +got error:\n%s", diff)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	type want struct {
		err error
	}

	cases := map[string]struct {
		reason string
		args   args
		want   want
	}{
		"ResourceIsNotBranchProtectionRule": {
			reason: "Must return an error resource is not BranchProtectionRule",
			args: args{
				mg: unexpectedObject,
			},
			want: want{
				err: errors.New(errUnexpectedObject),
			},
		},
		"DeleteFailed": {
			reason: "Must return error if DeleteBranchProtectionRule fails",
			args: args{
				mg: newBranchProtectionRule(),
				github: &fake.MockServiceBranchProtection{
					MockMutate: func(ctx context.Context, m interface{}, input githubv4.Input, variables map[string]interface{}) error {
						return errBoom
					},
				},
			},
			want: want{
				err: errors.Wrap(errBoom, errDeleteBranchProtectionRule),
			},
		},
		"Success": {
			reason: "Must not fail if all calls succeed",
			args: args{
				mg: newBranchProtectionRule(),
				github: &fake.MockServiceBranchProtection{
					MockMutate: func(ctx context.Context, m interface{}, input githubv4.Input, variables map[string]interface{}) error {
						return nil
					},
				},
			},
			want: want{
				err: nil,
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := external{
				gh: tc.args.github,
			}
			err := e.Delete(context.Background(), tc.args.mg)
			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("Delete(...): -want error, +got error:\n%s", diff)
			}
		})
	}
}
