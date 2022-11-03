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

package content

import (
	"context"
	"net/http"
	"testing"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/google/go-cmp/cmp"

	"github.com/crossplane-contrib/provider-github/apis/repositories/v1alpha1"
	"github.com/crossplane-contrib/provider-github/pkg/clients/content"
	"github.com/crossplane-contrib/provider-github/pkg/clients/content/fake"
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

type contentOption func(*v1alpha1.Content)

func newContent(opts ...contentOption) *v1alpha1.Content {
	r := &v1alpha1.Content{}

	for _, f := range opts {
		f(r)
	}
	return r
}

func withRepository(repository string) contentOption {
	return func(c *v1alpha1.Content) { c.Spec.ForProvider.Repository = &repository }
}

func withReconcile(reconcile string) contentOption {
	return func(c *v1alpha1.Content) { c.Spec.Reconcile = &reconcile }
}

func withBranch(b string) contentOption {
	return func(c *v1alpha1.Content) { c.Spec.ForProvider.Branch = &b }
}

func withConditions(condition xpv1.Condition) contentOption {
	return func(c *v1alpha1.Content) { c.SetConditions(condition) }
}

type args struct {
	kube   client.Client
	mg     resource.Managed
	github content.Service
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
		"ResourceIsNotContent": {
			reason: "Must return an error resource is not Repository",
			args: args{
				mg: unexpectedObject,
			},
			want: want{
				eo:  managed.ExternalObservation{},
				err: errors.New(errUnexpectedObject),
			},
		},
		"MustNotReturnError404": {
			reason: "Must not return an error if GET content returns 404 status code",
			args: args{
				mg: newContent(
					withBranch("main"),
					withRepository(dummyText),
				),
				github: &fake.MockService{
					MockGetContents: func(tx context.Context, owner, repo, path string, opts *github.RepositoryContentGetOptions) (fileContent *github.RepositoryContent, directoryContent []*github.RepositoryContent, resp *github.Response, err error) {
						return &github.RepositoryContent{},
							nil,
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
			reason: "Should return an error if GET content returns a status code different than 404",
			args: args{
				mg: newContent(
					withRepository(dummyText),
				),
				github: &fake.MockService{
					MockGetContents: func(tx context.Context, owner, repo, path string, opts *github.RepositoryContentGetOptions) (fileContent *github.RepositoryContent, directoryContent []*github.RepositoryContent, resp *github.Response, err error) {
						return &github.RepositoryContent{},
							nil,
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
				err: errors.Wrap(errBoom, errGetContent),
			},
		},
		"SkipReconcile": {
			reason: "Should skip reconcile logic if reconcile option is Disabled",
			args: args{
				mg: newContent(
					withReconcile("Disabled"),
					withConditions(xpv1.Available()),
					withRepository(dummyText),
				),
				github: &fake.MockService{
					MockGetContents: func(tx context.Context, owner, repo, path string, opts *github.RepositoryContentGetOptions) (fileContent *github.RepositoryContent, directoryContent []*github.RepositoryContent, resp *github.Response, err error) {
						return &github.RepositoryContent{},
							nil,
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
		"Success": {
			reason: "Should not return errors if everything goes well",
			args: args{
				mg: newContent(
					withReconcile("Disabled"),
					withRepository(dummyText),
				),
				github: &fake.MockService{
					MockGetContents: func(tx context.Context, owner, repo, path string, opts *github.RepositoryContentGetOptions) (fileContent *github.RepositoryContent, directoryContent []*github.RepositoryContent, resp *github.Response, err error) {
						return &github.RepositoryContent{
								HTMLURL: &dummyText,
								URL:     &dummyText,
								SHA:     &dummyText,
							},
							nil,
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
			e := contentExternal{
				client: tc.args.kube,
				gh:     tc.args.github,
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
		"ResourceIsNotContent": {
			reason: "Must return an error resource is not Repository",
			args: args{
				mg: unexpectedObject,
			},
			want: want{
				eo:  managed.ExternalCreation{},
				err: errors.New(errUnexpectedObject),
			},
		},
		"CreationFailed": {
			reason: "Must return an error if the repository creation fails",
			args: args{
				github: &fake.MockService{
					MockCreateFile: func(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentFileOptions) (*github.RepositoryContentResponse, *github.Response, error) {
						return nil,
							nil,
							errBoom
					},
				},
				mg: newContent(
					withRepository(dummyText),
				),
			},
			want: want{
				eo:  managed.ExternalCreation{},
				err: errors.Wrap(errBoom, errCreateContent),
			},
		},
		"Success": {
			reason: "Must not return an error if everything goes well",
			args: args{
				github: &fake.MockService{
					MockCreateFile: func(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentFileOptions) (*github.RepositoryContentResponse, *github.Response, error) {
						return &github.RepositoryContentResponse{},
							&github.Response{},
							nil
					},
				},
				mg: newContent(
					withRepository(dummyText),
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
			e := contentExternal{
				client: tc.args.kube,
				gh:     tc.args.github,
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
		"ResourceIsNotContent": {
			reason: "Must return an error if resource is not Content",
			args: args{
				mg: unexpectedObject,
			},
			want: want{
				eo:  managed.ExternalUpdate{},
				err: errors.New(errUnexpectedObject),
			},
		},
		"UpdateSuccessful": {
			reason: "Must not return an error if everything succeeds",
			args: args{
				mg: newContent(
					withBranch("main"),
					withRepository(dummyText),
				),
				github: &fake.MockService{
					MockUpdateFile: func(ctx context.Context, owner string, repo string, path string, opts *github.RepositoryContentFileOptions) (*github.RepositoryContentResponse, *github.Response, error) {
						return nil,
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
				eo:  managed.ExternalUpdate{},
				err: nil,
			},
		},
		"UpdateFailed": {
			reason: "Must return an error if the update request fails",
			args: args{
				mg: newContent(
					withBranch("main"),
					withRepository(dummyText),
				),
				github: &fake.MockService{
					MockUpdateFile: func(ctx context.Context, owner string, repo string, path string, opts *github.RepositoryContentFileOptions) (*github.RepositoryContentResponse, *github.Response, error) {
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
				eo:  managed.ExternalUpdate{},
				err: errors.Wrap(errBoom, errUpdateContent),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := contentExternal{
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
		"ResourceIsNotContent": {
			reason: "Must return an error if resource is not Content",
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
				mg: newContent(
					withBranch("main"),
					withRepository(dummyText),
				),
				github: &fake.MockService{
					MockDeleteFile: func(ctx context.Context, owner string, repo string, path string, opts *github.RepositoryContentFileOptions) (*github.RepositoryContentResponse, *github.Response, error) {
						return nil,
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
				err: nil,
			},
		},
		"DeleteFailed": {
			reason: "Must return an error if the delete request fails",
			args: args{
				mg: newContent(
					withBranch("main"),
					withRepository(dummyText),
				),
				github: &fake.MockService{
					MockDeleteFile: func(ctx context.Context, owner string, repo string, path string, opts *github.RepositoryContentFileOptions) (*github.RepositoryContentResponse, *github.Response, error) {
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
				err: errors.Wrap(errBoom, errDeleteContent),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := contentExternal{
				gh: tc.args.github,
			}
			err := e.Delete(context.Background(), tc.args.mg)
			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("Delete(...): -want error, +got error:\n%s", diff)
			}
		})
	}
}
