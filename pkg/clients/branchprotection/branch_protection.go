/*
Copyright 2022 The Crossplane Authors.
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

	ghclient "github.com/crossplane-contrib/provider-github/pkg/clients"
	"github.com/shurcooL/githubv4"
)

// Service defines the GraphQL operations
type Service interface {
	Query(ctx context.Context, q interface{}, variables map[string]interface{}) error
	Mutate(ctx context.Context, m interface{}, input githubv4.Input, variables map[string]interface{}) error
}

// NewClient creates a new *githubv4.Client
func NewClient(token string) (Service, error) {
	c, err := ghclient.NewV4Client(token)
	if err != nil {
		return nil, err
	}

	return c, nil
}
