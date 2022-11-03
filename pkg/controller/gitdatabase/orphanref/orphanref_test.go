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

package orphanref

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/crossplane-contrib/provider-github/apis/gitdatabase/v1alpha1"
	"github.com/crossplane-contrib/provider-github/pkg/clients/orphanref"
	"github.com/crossplane-contrib/provider-github/pkg/clients/orphanref/fake"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/crossplane/crossplane-runtime/pkg/test"
	"github.com/google/go-github/v48/github"
	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	unexpectedObject resource.Managed
	errBoom          = errors.New("boom")
	notFound         = 404
	internalError    = 500
	ok               = 200
	dummyText        = "example"
)

type orphanrefOption func(*v1alpha1.OrphanRef)

func newOrphanRef(opts ...orphanrefOption) *v1alpha1.OrphanRef {
	r := &v1alpha1.OrphanRef{}

	for _, f := range opts {
		f(r)
	}
	return r
}

func withRepository(repository string) orphanrefOption {
	return func(c *v1alpha1.OrphanRef) { c.Spec.ForProvider.Repository = &repository }
}

func withRef(b string) orphanrefOption {
	return func(c *v1alpha1.OrphanRef) {
		name := fmt.Sprintf("refs/heads/%s", b)
		meta.SetExternalName(c, name)
	}
}

type args struct {
	kube client.Client
	mg   resource.Managed
	gh   orphanref.Service
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
		"ResourceIsNotOrphanRef": {
			reason: "Must return an error resource is not OrphanRef",
			args: args{
				mg: unexpectedObject,
			},
			want: want{
				eo:  managed.ExternalObservation{},
				err: errors.New(errUnexpectedObject),
			},
		},
		"RepositoryFieldCannotBeNil": {
			reason: "Must return an error if repository field is nil",
			args: args{
				mg: &v1alpha1.OrphanRef{},
			},
			want: want{
				eo:  managed.ExternalObservation{},
				err: errors.New(errRepositoryEmpty),
			},
		},
		"MustNotReturnError404": {
			reason: "Must not return an error if GET ref returns 404 status code",
			args: args{
				mg: newOrphanRef(
					withRef("main"),
					withRepository(dummyText),
				),
				gh: &fake.MockService{
					MockGetRef: func(ctx context.Context, owner, repo, ref string) (*github.Reference, *github.Response, error) {
						return nil,
							&github.Response{
								Response: &http.Response{
									StatusCode: notFound,
								},
							},
							errBoom
					},
				},
			},
			want: want{
				eo:  managed.ExternalObservation{},
				err: nil,
			},
		},
		"InternalError": {
			reason: "Should return an error if GET ref returns a status code different than 404",
			args: args{
				mg: newOrphanRef(
					withRepository(dummyText),
				),
				gh: &fake.MockService{
					MockGetRef: func(ctx context.Context, owner, repo, ref string) (*github.Reference, *github.Response, error) {
						return nil,
							&github.Response{
								Response: &http.Response{
									StatusCode: internalError,
								},
							},
							errBoom
					},
				},
			},
			want: want{
				eo:  managed.ExternalObservation{},
				err: errors.Wrap(errBoom, errGetOrphanRef),
			},
		},
		"Success": {
			reason: "Should not return errors if resource is found",
			args: args{
				mg: newOrphanRef(
					withRef("main"),
					withRepository(dummyText),
				),
				gh: &fake.MockService{
					MockGetRef: func(ctx context.Context, owner, repo, ref string) (*github.Reference, *github.Response, error) {
						return &github.Reference{
								Ref: github.String("refs/heads/main"),
								URL: github.String("fakeurl.com"),
							},
							&github.Response{
								Response: &http.Response{
									StatusCode: ok,
								},
							},
							nil
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
			e := orphanRefExternal{
				client: tc.args.kube,
				gh:     tc.args.gh,
			}
			got, err := e.Observe(context.Background(), tc.args.mg)
			if diff := cmp.Diff(tc.want.eo, got); diff != "" {
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
		"ResourceIsNotOrphanRef": {
			reason: "Must return an error resource is not Repository",
			args: args{
				mg: unexpectedObject,
			},
			want: want{
				eo:  managed.ExternalCreation{},
				err: errors.New(errUnexpectedObject),
			},
		},
		"RepositoryFieldCannotBeNil": {
			reason: "Must return an error if repository field is nil",
			args: args{
				mg: &v1alpha1.OrphanRef{},
			},
			want: want{
				eo:  managed.ExternalCreation{},
				err: errors.New(errRepositoryEmpty),
			},
		},
		"TreeCreationFailed": {
			reason: "Must return an error if the git tree creation fails",
			args: args{
				mg: newOrphanRef(
					withRepository(dummyText),
					withRef("main"),
				),
				gh: &fake.MockService{
					MockCreateTree: func(ctx context.Context, owner, repo, baseTree string, entries []*github.TreeEntry) (*github.Tree, *github.Response, error) {
						return nil, nil, errBoom
					},
				},
			},
			want: want{
				eo:  managed.ExternalCreation{},
				err: errors.Wrap(errBoom, errCreateOrphanRef),
			},
		},
		"CommitCreationFailed": {
			reason: "Must return an error if the initial commit creation fails",
			args: args{
				mg: newOrphanRef(
					withRepository(dummyText),
					withRef("main"),
				),
				gh: &fake.MockService{
					MockCreateTree: func(ctx context.Context, owner, repo, baseTree string, entries []*github.TreeEntry) (*github.Tree, *github.Response, error) {
						return &github.Tree{}, nil, nil
					},
					MockCreateCommit: func(ctx context.Context, owner, repo string, commit *github.Commit) (*github.Commit, *github.Response, error) {
						return nil, nil, errBoom
					},
				},
			},
			want: want{
				eo:  managed.ExternalCreation{},
				err: errors.Wrap(errBoom, errCreateOrphanRef),
			},
		},
		"RefCreationFailed": {
			reason: "Must return an error if the initial commit creation fails",
			args: args{
				mg: newOrphanRef(
					withRepository(dummyText),
					withRef("main"),
				),
				gh: &fake.MockService{
					MockCreateTree: func(ctx context.Context, owner, repo, baseTree string, entries []*github.TreeEntry) (*github.Tree, *github.Response, error) {
						return &github.Tree{}, nil, nil
					},
					MockCreateCommit: func(ctx context.Context, owner, repo string, commit *github.Commit) (*github.Commit, *github.Response, error) {
						return &github.Commit{
							SHA: github.String("abc"),
						}, nil, nil
					},
					MockCreateRef: func(ctx context.Context, owner, repo string, ref *github.Reference) (*github.Reference, *github.Response, error) {
						return nil, nil, errBoom
					},
				},
			},
			want: want{
				eo:  managed.ExternalCreation{},
				err: errors.Wrap(errBoom, errCreateOrphanRef),
			},
		},
		"Success": {
			reason: "Must not return an error if everything goes well",
			args: args{
				mg: newOrphanRef(
					withRepository(dummyText),
					withRef("main"),
				),
				gh: &fake.MockService{
					MockCreateTree: func(ctx context.Context, owner, repo, baseTree string, entries []*github.TreeEntry) (*github.Tree, *github.Response, error) {
						return &github.Tree{}, nil, nil
					},
					MockCreateCommit: func(ctx context.Context, owner, repo string, commit *github.Commit) (*github.Commit, *github.Response, error) {
						return &github.Commit{
							SHA: github.String("abc"),
						}, nil, nil
					},
					MockCreateRef: func(ctx context.Context, owner, repo string, ref *github.Reference) (*github.Reference, *github.Response, error) {
						return &github.Reference{}, nil, nil
					},
				},
			},
			want: want{
				eo:  managed.ExternalCreation{},
				err: nil,
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := orphanRefExternal{
				client: tc.args.kube,
				gh:     tc.args.gh,
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

func TestDelete(t *testing.T) {
	type want struct {
		err error
	}

	cases := map[string]struct {
		reason string
		args   args
		want   want
	}{
		"ResourceIsNotOrphanRef": {
			reason: "Must return an error if resource is not OrphanRef",
			args: args{
				mg: unexpectedObject,
			},
			want: want{
				err: errors.New(errUnexpectedObject),
			},
		},
		"DeleteSuccessful": {
			reason: "Must not return an error if everything succeeds",
			args: args{
				mg: newOrphanRef(
					withRef("main"),
					withRepository(dummyText),
				),
				gh: &fake.MockService{
					MockDeleteRef: func(ctx context.Context, owner, repo, ref string) (*github.Response, error) {
						return &github.Response{
								Response: &http.Response{
									StatusCode: ok,
								},
							},
							nil
					},
				},
			},
			want: want{
				err: nil,
			},
		},
		"DeleteFailed": {
			reason: "Must return an error if the delete request fails",
			args: args{
				mg: newOrphanRef(
					withRef("main"),
					withRepository(dummyText),
				),
				gh: &fake.MockService{
					MockDeleteRef: func(ctx context.Context, owner, repo, ref string) (*github.Response, error) {
						return &github.Response{
								Response: &http.Response{
									StatusCode: internalError,
								},
							},
							errBoom
					},
				},
			},
			want: want{
				err: errors.Wrap(errBoom, errDeleteOrphanRef),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := orphanRefExternal{
				gh: tc.args.gh,
			}
			err := e.Delete(context.Background(), tc.args.mg)
			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("Delete(...): -want error, +got error:\n%s", diff)
			}
		})
	}
}
