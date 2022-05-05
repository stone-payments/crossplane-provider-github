/*
Copyright 2022 The Crossplane Authors.
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

package v1alpha1

import (
	"context"

	"github.com/crossplane-contrib/provider-github/apis/repositories/v1alpha1"
	"github.com/crossplane/crossplane-runtime/pkg/reference"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// RepositoryName extracts the Name of a Repository.
func RepositoryName() reference.ExtractValueFn {
	return func(mg resource.Managed) string {
		r, ok := mg.(*v1alpha1.Repository)
		if !ok {
			return ""
		}
		return r.Status.AtProvider.Name
	}
}

// ResolveReferences of this OrphanRef
func (mg *OrphanRef) ResolveReferences(ctx context.Context, c client.Reader) error {
	r := reference.NewAPIResolver(c, mg)

	// Resolve spec.forProvider.repository
	rsp, err := r.Resolve(ctx, reference.ResolutionRequest{
		CurrentValue: reference.FromPtrValue(mg.Spec.ForProvider.Repository),
		Reference:    mg.Spec.ForProvider.RepositoryRef,
		Selector:     mg.Spec.ForProvider.RepositorySelector,
		To:           reference.To{Managed: &v1alpha1.Repository{}, List: &v1alpha1.RepositoryList{}},
		Extract:      RepositoryName(),
	})
	if err != nil {
		return errors.Wrap(err, "spec.forProvider.repository")
	}
	mg.Spec.ForProvider.Repository = reference.ToPtrValue(rsp.ResolvedValue)
	mg.Spec.ForProvider.RepositoryRef = rsp.ResolvedReference

	return nil
}
