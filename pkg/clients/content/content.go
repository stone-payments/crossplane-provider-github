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

	"github.com/crossplane-contrib/provider-github/apis/repositories/v1alpha1"
	ghclient "github.com/crossplane-contrib/provider-github/pkg/clients"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/v48/github"
	"github.com/mitchellh/copystructure"
)

// Service defines the Content operations
type Service interface {
	CreateFile(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentFileOptions) (*github.RepositoryContentResponse, *github.Response, error)
	GetContents(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentGetOptions) (fileContent *github.RepositoryContent, directoryContent []*github.RepositoryContent, resp *github.Response, err error)
	UpdateFile(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentFileOptions) (*github.RepositoryContentResponse, *github.Response, error)
	DeleteFile(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentFileOptions) (*github.RepositoryContentResponse, *github.Response, error)
}

// NewService creates a new Service based on the *github.Client
// returned by the NewClient SDK method.
func NewService(token string) *Service {
	c := ghclient.NewClient(token)
	r := Service(c.Repositories)
	return &r
}

// GenerateObservation generates a v1alpha1.ContentObservation
func GenerateObservation(content github.RepositoryContent) v1alpha1.ContentObservation {
	return v1alpha1.ContentObservation{
		HTMLURL: ghclient.StringValue(content.HTMLURL),
		URL:     ghclient.StringValue(content.URL),
		SHA:     ghclient.StringValue(content.SHA),
	}
}

// IsUpToDate checks if the spec is up to date with the external provider API
func IsUpToDate(params *v1alpha1.ContentParameters, content *github.RepositoryContent) (bool, error) {
	generated, err := copystructure.Copy(content)
	if err != nil {
		return true, err
	}
	clone, ok := generated.(*github.RepositoryContent)
	if !ok {
		return true, err
	}

	desired := OverrideParameters(*params, *clone)

	decodedContent, err := content.GetContent()
	if err != nil {
		return true, err
	}
	content.Content = &decodedContent

	return cmp.Equal(
		desired,
		*content,
	), nil
}

// OverrideParameters overrides the fields in github.RepositoryContent that were defined in
// v1alpha1.ContentParameters
func OverrideParameters(params v1alpha1.ContentParameters, content github.RepositoryContent) github.RepositoryContent {
	content.Content = &params.Content
	return content
}
