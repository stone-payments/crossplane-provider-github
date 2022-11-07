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

	"github.com/pkg/errors"
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
	errUnexpectedObject = "The managed resource is not a BranchProtectionRule resource"
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
	_, ok := mgd.(*v1alpha1.BranchProtectionRule)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errUnexpectedObject)
	}

	return managed.ExternalObservation{}, nil
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
