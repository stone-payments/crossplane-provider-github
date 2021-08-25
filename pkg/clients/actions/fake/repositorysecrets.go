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
)

// MockServiceRepositorysecret is a mock implementation of the Service
type MockServiceRepositorysecret struct {
	MockGetRepoSecret            func(ctx context.Context, owner, repo, name string) (*github.Secret, *github.Response, error)
	MockGetRepoPublicKey         func(ctx context.Context, owner, repo string) (*github.PublicKey, *github.Response, error)
	MockCreateOrUpdateRepoSecret func(ctx context.Context, owner, repo string, eSecret *github.EncryptedSecret) (*github.Response, error)
	MockDeleteRepoSecret         func(ctx context.Context, owner, repo, name string) (*github.Response, error)
}

// GetRepoSecret is a fake GetRepoSecret SDK method
func (m *MockServiceRepositorysecret) GetRepoSecret(ctx context.Context, owner, repo, name string) (*github.Secret, *github.Response, error) {
	return m.MockGetRepoSecret(ctx, owner, repo, name)
}

// GetRepoPublicKey is a fake GetRepoPublicKey SDK method
func (m *MockServiceRepositorysecret) GetRepoPublicKey(ctx context.Context, owner, repo string) (*github.PublicKey, *github.Response, error) {
	return m.MockGetRepoPublicKey(ctx, owner, repo)
}

// CreateOrUpdateRepoSecret is a fake CreateOrUpdateRepoSecret SDK method
func (m *MockServiceRepositorysecret) CreateOrUpdateRepoSecret(ctx context.Context, owner, repo string, eSecret *github.EncryptedSecret) (*github.Response, error) {
	return m.MockCreateOrUpdateRepoSecret(ctx, owner, repo, eSecret)
}

// DeleteRepoSecret is a fake DeleteRepoSecret SDK method
func (m *MockServiceRepositorysecret) DeleteRepoSecret(ctx context.Context, owner, repo, name string) (*github.Response, error) {
	return m.MockDeleteRepoSecret(ctx, owner, repo, name)
}
