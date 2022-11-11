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
package fake

import (
	"context"

	"github.com/crossplane-contrib/provider-github/pkg/clients/branchprotection"
	"github.com/shurcooL/githubv4"
)

// This ensures that the mock implements the Service interface
var _ branchprotection.Service = (*MockServiceBranchProtection)(nil)

// MockServiceBranchProtection is a mock implementation of the Service
type MockServiceBranchProtection struct {
	MockQuery  func(ctx context.Context, q interface{}, variables map[string]interface{}) error
	MockMutate func(ctx context.Context, m interface{}, input githubv4.Input, variables map[string]interface{}) error
}

// Query is a fake Query SDK method
func (m *MockServiceBranchProtection) Query(ctx context.Context, q interface{}, variables map[string]interface{}) error {
	return m.MockQuery(ctx, q, variables)
}

// Mutate is a fake Mutate SDK method
func (m *MockServiceBranchProtection) Mutate(ctx context.Context, mutate interface{}, input githubv4.Input, variables map[string]interface{}) error {
	return m.MockMutate(ctx, mutate, input, variables)
}
