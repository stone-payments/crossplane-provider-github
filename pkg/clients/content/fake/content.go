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

	"github.com/google/go-github/v33/github"

	"github.com/crossplane-contrib/provider-github/pkg/clients/content"
)

// This ensures that the mock implements the Service interface
var _ content.Service = (*MockService)(nil)

// MockService is a mock implementation of the Service
type MockService struct {
	MockGetContents func(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentGetOptions) (fileContent *github.RepositoryContent, directoryContent []*github.RepositoryContent, resp *github.Response, err error)
	MockCreateFile  func(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentFileOptions) (*github.RepositoryContentResponse, *github.Response, error)
	MockUpdateFile  func(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentFileOptions) (*github.RepositoryContentResponse, *github.Response, error)
	MockDeleteFile  func(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentFileOptions) (*github.RepositoryContentResponse, *github.Response, error)
}

// CreateFile is a fake Create SDK method
func (m *MockService) CreateFile(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentFileOptions) (*github.RepositoryContentResponse, *github.Response, error) {
	return m.MockCreateFile(ctx, owner, repo, path, opts)
}

// GetContents is a fake GetContents SDK method
func (m *MockService) GetContents(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentGetOptions) (fileContent *github.RepositoryContent, directoryContent []*github.RepositoryContent, resp *github.Response, err error) {
	return m.MockGetContents(ctx, owner, repo, path, opts)
}

// UpdateFile is a fake UpdateFile SDK method
func (m *MockService) UpdateFile(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentFileOptions) (*github.RepositoryContentResponse, *github.Response, error) {
	return m.MockUpdateFile(ctx, owner, repo, path, opts)
}

// DeleteFile is a fake DeleteFile SDK method
func (m *MockService) DeleteFile(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentFileOptions) (*github.RepositoryContentResponse, *github.Response, error) {
	return m.MockDeleteFile(ctx, owner, repo, path, opts)
}
