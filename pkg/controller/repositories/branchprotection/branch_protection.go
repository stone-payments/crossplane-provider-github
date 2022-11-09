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
package branchprotection

import (
	"context"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"github.com/shurcooL/githubv4"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/ratelimiter"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"

	"github.com/crossplane-contrib/provider-github/apis/repositories/v1alpha1"
	ghclient "github.com/crossplane-contrib/provider-github/pkg/clients"
	"github.com/crossplane-contrib/provider-github/pkg/clients/branchprotection"
)

const (
	errUnexpectedObject          = "The managed resource is not a BranchProtectionRule resource"
	errCheckBranchProtectionRule = "Cannot check if GitHub BranchProtectionRule exists"
	errGetBranchProtectionRule   = "Cannot get GitHub BranchProtectionRule"
	errCheckUpToDate             = "unable to determine if external resource is up to date"
)

// SetupBranchProtectionRule adds a controller that reconciles BranchProtectionRule.
func SetupBranchProtectionRule(mgr ctrl.Manager, l logging.Logger, rl workqueue.RateLimiter) error {
	name := managed.ControllerName(v1alpha1.BranchProtectionRuleGroupKind)

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(controller.Options{
			RateLimiter: ratelimiter.NewDefaultManagedRateLimiter(rl),
		}).
		For(&v1alpha1.BranchProtectionRule{}).
		Complete(managed.NewReconciler(mgr,
			resource.ManagedKind(v1alpha1.BranchProtectionRuleGroupVersionKind),
			managed.WithExternalConnecter(
				&connector{
					client:      mgr.GetClient(),
					newClientFn: branchprotection.NewClient,
				},
			),
			managed.WithConnectionPublishers(),
			managed.WithReferenceResolver(managed.NewAPISimpleReferenceResolver(mgr.GetClient())),
			managed.WithInitializers(
				managed.NewDefaultProviderConfig(mgr.GetClient()),
				managed.NewNameAsExternalName(mgr.GetClient()),
			),
			managed.WithLogger(l.WithValues("controller", name)),
			managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name)))))
}

type connector struct {
	client      client.Client
	newClientFn func(string) (branchprotection.Service, error)
}

func (c *connector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	cr, ok := mg.(*v1alpha1.BranchProtectionRule)
	if !ok {
		return nil, errors.New(errUnexpectedObject)
	}
	cfg, err := ghclient.GetConfig(ctx, c.client, cr)
	if err != nil {
		return nil, err
	}

	client, err := c.newClientFn(string(cfg))
	if err != nil {
		return nil, err
	}

	return &external{
		gh:     client,
		client: c.client,
	}, nil
}

type external struct {
	gh     branchprotection.Service
	client client.Client
}

func (e *external) Observe(ctx context.Context, mgd resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mgd.(*v1alpha1.BranchProtectionRule)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errUnexpectedObject)
	}

	isCreated, err := e.CheckBranchProtectionRuleExistance(ctx, cr)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errCheckBranchProtectionRule)
	}

	if !isCreated {
		return managed.ExternalObservation{
			ResourceExists: false,
		}, nil
	}

	external, err := e.GetBranchProtectionRule(ctx, cr.Status.AtProvider.ID)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errGetBranchProtectionRule)
	}

	currentSpec := cr.Spec.ForProvider.DeepCopy()
	branchprotection.LateInitialize(&cr.Spec.ForProvider, external)
	lateInitialized := !cmp.Equal(currentSpec, &cr.Spec.ForProvider)

	cr.Status.SetConditions(xpv1.Available())

	diff, err := branchprotection.IsUpToDate(cr.Spec.ForProvider, external)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errCheckUpToDate)
	}

	return managed.ExternalObservation{
		ResourceExists:          true,
		ResourceUpToDate:        diff == "",
		ResourceLateInitialized: lateInitialized,
		Diff:                    diff,
	}, nil
}

func (e *external) Create(ctx context.Context, mgd resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mgd.(*v1alpha1.BranchProtectionRule)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errUnexpectedObject)
	}
	cr.SetConditions(xpv1.Creating())
	fmt.Println("CREATE BRANCHPROTECTIONRULE")

	return managed.ExternalCreation{}, nil
}

func (e *external) Update(ctx context.Context, mgd resource.Managed) (managed.ExternalUpdate, error) {
	_, ok := mgd.(*v1alpha1.BranchProtectionRule)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errUnexpectedObject)
	}

	fmt.Println("UPDATE BRANCHPROTECTIONRULE")

	return managed.ExternalUpdate{}, nil
}

func (e *external) Delete(ctx context.Context, mgd resource.Managed) error {
	_, ok := mgd.(*v1alpha1.BranchProtectionRule)
	if !ok {
		return errors.New(errUnexpectedObject)
	}

	fmt.Println("DELETE BRANCHPROTECTIONRULE")

	return nil
}

// CheckBranchProtectionRuleExistance checks if a BranchProtectionRule pattern
// exists in the desired repository.
func (e *external) CheckBranchProtectionRuleExistance(ctx context.Context, cr *v1alpha1.BranchProtectionRule) (bool, error) {
	var query struct {
		Repository struct {
			ID                    githubv4.String `graphql:"id"`
			BranchProtectionRules struct {
				Nodes []struct {
					Pattern githubv4.String `graphql:"pattern"`
					ID      githubv4.String `graphql:"id"`
				} `graphql:"nodes"`
			} `graphql:"branchProtectionRules(first: 100)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	if cr.Spec.ForProvider.Repository == nil {
		return false, errors.New("required spec.forProvider.repository field is empty")
	}

	variables := map[string]interface{}{
		"owner": githubv4.String(cr.Spec.ForProvider.Owner),
		"name":  githubv4.String(*cr.Spec.ForProvider.Repository),
	}

	if err := e.gh.Query(ctx, &query, variables); err != nil {
		return false, err
	}

	cr.Spec.ForProvider.RepositoryID = (*string)(&query.Repository.ID)

	for _, node := range query.Repository.BranchProtectionRules.Nodes {
		if node.Pattern == githubv4.String(cr.Spec.ForProvider.Pattern) {
			cr.Status.AtProvider.ID = string(node.ID)
			return true, nil
		}
	}

	return false, nil
}

// GetBranchProtectionRule fetches the state of the desired
// BranchProtectionRule in GitHub
func (e *external) GetBranchProtectionRule(ctx context.Context, id string) (branchprotection.ExternalBranchProtectionRule, error) {
	var query struct {
		Node struct {
			Node branchprotection.ExternalBranchProtectionRule `graphql:"... on BranchProtectionRule"`
		} `graphql:"node(id: $id)"`
	}

	variables := map[string]interface{}{
		"id": githubv4.ID(id),
	}

	if err := e.gh.Query(ctx, &query, variables); err != nil {
		return branchprotection.ExternalBranchProtectionRule{}, err
	}

	return query.Node.Node, nil
}
