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
package organizations

import (
	"context"

	"github.com/google/go-github/v48/github"
	"github.com/pkg/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/controller"
	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/ratelimiter"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"

	"github.com/crossplane-contrib/provider-github/apis/organizations/v1alpha1"
	ghclient "github.com/crossplane-contrib/provider-github/pkg/clients"
)

const (
	errUnexpectedObject = "The managed resource is not a Membership resource"
)

// SetupMembership adds a controller that reconciles Memberships.
func SetupMembership(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(v1alpha1.MembershipGroupKind)

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.MembershipGroupVersionKind),
		managed.WithExternalConnecter(&connector{client: mgr.GetClient(), newClientFn: ghclient.NewV3Client}),
		managed.WithConnectionPublishers(),
		managed.WithReferenceResolver(managed.NewAPISimpleReferenceResolver(mgr.GetClient())),
		managed.WithInitializers(managed.NewDefaultProviderConfig(mgr.GetClient())),
		managed.WithPollInterval(o.PollInterval),
		managed.WithLogger(o.Logger.WithValues("controller", name)),
		managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name))),
	)

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(o.ForControllerRuntime()).
		For(&v1alpha1.Membership{}).
		Complete(ratelimiter.NewReconciler(name, r, o.GlobalRateLimiter))
}

type connector struct {
	client      client.Client
	newClientFn func(string) (*github.Client, error)
}

func (c *connector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	cr, ok := mg.(*v1alpha1.Membership)
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
		client: client,
		kube:   c.client,
	}, nil
}

type external struct {
	client *github.Client
	kube   client.Client
}

func (e *external) Observe(ctx context.Context, mgd resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mgd.(*v1alpha1.Membership)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errUnexpectedObject)
	}

	// TODO(hasheddan): handle errors correctly
	m, _, err := e.client.Organizations.GetOrgMembership(ctx, cr.Spec.ForProvider.User, cr.Spec.ForProvider.Organization)
	if err != nil { // nolint:nilerr
		return managed.ExternalObservation{
			ResourceExists: false,
		}, nil
	}

	if m.State != nil && *m.State == "active" {
		cr.SetConditions(xpv1.Available())
	} else {
		cr.SetConditions(xpv1.Creating())
	}

	return managed.ExternalObservation{
		ResourceUpToDate: true,
		ResourceExists:   true,
	}, nil
}

func (e *external) Create(ctx context.Context, mgd resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mgd.(*v1alpha1.Membership)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errUnexpectedObject)
	}

	inv := &github.CreateOrgInvitationOptions{
		InviteeID: cr.Spec.ForProvider.InviteeID,
		Email:     cr.Spec.ForProvider.Email,
		Role:      cr.Spec.ForProvider.Role,
		TeamID:    []int64{},
	}
	_, _, err := e.client.Organizations.CreateOrgInvitation(ctx, cr.Spec.ForProvider.Organization, inv)
	if err != nil {
		return managed.ExternalCreation{}, err
	}
	return managed.ExternalCreation{ExternalNameAssigned: true}, nil

}

func (e *external) Update(ctx context.Context, mgd resource.Managed) (managed.ExternalUpdate, error) { // nolint:gocyclo
	_, ok := mgd.(*v1alpha1.Membership)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errUnexpectedObject)
	}

	return managed.ExternalUpdate{}, nil
}

func (e *external) Delete(ctx context.Context, mgd resource.Managed) error {
	cr, ok := mgd.(*v1alpha1.Membership)
	if !ok {
		return errors.New(errUnexpectedObject)
	}

	_, err := e.client.Organizations.RemoveMember(ctx, cr.Spec.ForProvider.Organization, cr.Spec.ForProvider.User)

	return err
}
