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

package fake

import (
	"context"

	"github.com/google/go-github/v48/github"

	"github.com/crossplane-contrib/provider-github/pkg/clients/orphanref"
)

// This ensures that the mock implements the Service interface
var _ orphanref.Service = (*MockService)(nil)

// MockService is a mock implementation of the Service
type MockService struct {
	MockCreateTree   func(ctx context.Context, owner string, repo string, baseTree string, entries []*github.TreeEntry) (*github.Tree, *github.Response, error)
	MockCreateCommit func(ctx context.Context, owner string, repo string, commit *github.Commit) (*github.Commit, *github.Response, error)
	MockCreateRef    func(ctx context.Context, owner string, repo string, ref *github.Reference) (*github.Reference, *github.Response, error)
	MockGetRef       func(ctx context.Context, owner string, repo string, ref string) (*github.Reference, *github.Response, error)
	MockDeleteRef    func(ctx context.Context, owner string, repo string, ref string) (*github.Response, error)
}

// CreateTree is a mock implementation that redirects to MockCreateTree func field
func (m *MockService) CreateTree(ctx context.Context, owner string, repo string, baseTree string, entries []*github.TreeEntry) (*github.Tree, *github.Response, error) {
	return m.MockCreateTree(ctx, owner, repo, baseTree, entries)
}

// CreateCommit is a mock implementation that redirects to MockCreateCommit func field
func (m *MockService) CreateCommit(ctx context.Context, owner string, repo string, commit *github.Commit) (*github.Commit, *github.Response, error) {
	return m.MockCreateCommit(ctx, owner, repo, commit)
}

// CreateRef is a mock implementation that redirects to MockCreateRef func field
func (m *MockService) CreateRef(ctx context.Context, owner string, repo string, ref *github.Reference) (*github.Reference, *github.Response, error) {
	return m.MockCreateRef(ctx, owner, repo, ref)
}

// GetRef is a mock implementation that redirects to MockGetRef func field
func (m *MockService) GetRef(ctx context.Context, owner string, repo string, ref string) (*github.Reference, *github.Response, error) {
	return m.MockGetRef(ctx, owner, repo, ref)
}

// DeleteRef is a mock implementation that redirects to MockDeleteRef func field
func (m *MockService) DeleteRef(ctx context.Context, owner string, repo string, ref string) (*github.Response, error) {
	return m.MockDeleteRef(ctx, owner, repo, ref)
}
