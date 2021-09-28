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

	ghclient "github.com/crossplane-contrib/provider-github/pkg/clients"
	"github.com/google/go-github/v33/github"
)

// Service defines the Content operations
type Service interface {
	CreateFile(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentFileOptions) (*github.RepositoryContentResponse, *github.Response, error)
	GetContents(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentGetOptions) (fileContent *github.RepositoryContent, directoryContent []*github.RepositoryContent, resp *github.Response, err error)
}

// NewService creates a new Service based on the *github.Client
// returned by the NewClient SDK method.
func NewService(token string) *Service {
	c := ghclient.NewClient(token)
	r := Service(c.Repositories)
	return &r
}
