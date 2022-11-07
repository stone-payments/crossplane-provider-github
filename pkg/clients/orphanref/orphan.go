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
package orphanref

import (
	"context"

	"github.com/crossplane-contrib/provider-github/apis/gitdatabase/v1alpha1"
	ghclient "github.com/crossplane-contrib/provider-github/pkg/clients"
	"github.com/google/go-github/v48/github"
)

// Service defines the OrphanRef operations
type Service interface {
	CreateTree(ctx context.Context, owner string, repo string, baseTree string, entries []*github.TreeEntry) (*github.Tree, *github.Response, error)
	CreateCommit(ctx context.Context, owner string, repo string, commit *github.Commit) (*github.Commit, *github.Response, error)
	CreateRef(ctx context.Context, owner string, repo string, ref *github.Reference) (*github.Reference, *github.Response, error)
	GetRef(ctx context.Context, owner string, repo string, ref string) (*github.Reference, *github.Response, error)
	DeleteRef(ctx context.Context, owner string, repo string, ref string) (*github.Response, error)
}

// NewService creates a new Service based on the *github.Client
// returned by the NewClient SDK method.
func NewService(token string) (*Service, error) {
	c, err := ghclient.NewV3Client(token)
	if err != nil {
		return nil, err
	}

	r := Service(c.Git)
	return &r, nil
}

// GenerateObservation generates a v1alpha1.OrphanRefObservation
func GenerateObservation(ref *github.Reference) v1alpha1.OrphanRefObservation {
	return v1alpha1.OrphanRefObservation{
		URL: ghclient.StringValue(ref.URL),
		Ref: ghclient.StringValue(ref.Ref),
	}
}
