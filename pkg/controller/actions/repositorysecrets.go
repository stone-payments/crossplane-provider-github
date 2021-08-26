package repositorysecret

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

import (
	"context"

	"github.com/pkg/errors"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/crossplane/crossplane-runtime/pkg/ratelimiter"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"

	"github.com/crossplane-contrib/provider-github/apis/actions/v1alpha1"
	ghclient "github.com/crossplane-contrib/provider-github/pkg/clients"
	repositorysecret "github.com/crossplane-contrib/provider-github/pkg/clients/actions"
)

const (
	errUnexpectedObject        = "The managed resource is not a Repository Secrets resource"
	errCreateRepositorySecrets = "cannot create Repository Secrets"
	errUpdateRepositorySecrets = "cannot update Repository Secrets"
	errDeleteRepositorySecrets = "cannot delete Repository Secrets"
)

// SetupRepositorySecret adds a controller that reconciles secrets.
func SetupRepositorySecret(mgr ctrl.Manager, l logging.Logger, rl workqueue.RateLimiter) error {
	name := managed.ControllerName(v1alpha1.RepositorySecretGroupKind)

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(controller.Options{
			RateLimiter: ratelimiter.NewDefaultManagedRateLimiter(rl),
		}).
		For(&v1alpha1.RepositorySecret{}).
		Complete(managed.NewReconciler(mgr,
			resource.ManagedKind(v1alpha1.RepositorySecretGroupVersionKind),
			managed.WithExternalConnecter(
				&connector{
					client:      mgr.GetClient(),
					newClientFn: repositorysecret.NewService,
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
	newClientFn func(string) *repositorysecret.Service
}

func (c *connector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	cr, ok := mg.(*v1alpha1.RepositorySecret)
	if !ok {
		return nil, errors.New(errUnexpectedObject)
	}
	cfg, err := ghclient.GetConfig(ctx, c.client, cr)
	if err != nil {
		return nil, err
	}
	return &external{*c.newClientFn(string(cfg)), c.client}, nil
}

type external struct {
	gh     repositorysecret.Service
	client client.Client
}

func (e *external) Observe(ctx context.Context, mgd resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mgd.(*v1alpha1.RepositorySecret)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errUnexpectedObject)
	}

	cr.Status.SetConditions(xpv1.Available())
	if len(cr.Status.AtProvider.EncryptValue) == 0 {
		return managed.ExternalObservation{}, nil
	}

	upToDate, err := repositorysecret.IsUpToDate(ctx, e.client, &cr.Spec.ForProvider, &cr.Status.AtProvider, meta.GetExternalName(cr), e.gh)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, "Error to verify if is up to date")
	}

	return managed.ExternalObservation{
		ResourceUpToDate: upToDate,
		ResourceExists:   true,
	}, nil
}

func (e *external) Create(ctx context.Context, mgd resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mgd.(*v1alpha1.RepositorySecret)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errUnexpectedObject)
	}

	hash, time, err := repositorysecret.CreateOrUpdateSec(ctx, &cr.Spec.ForProvider, meta.GetExternalName(cr), e.client, e.gh)
	if err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, errCreateRepositorySecrets)
	}

	cr.Status.AtProvider.LastUpdate = time
	cr.Status.AtProvider.EncryptValue = hash
	cr.SetConditions(xpv1.Creating())
	return managed.ExternalCreation{}, nil
}

func (e *external) Update(ctx context.Context, mgd resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mgd.(*v1alpha1.RepositorySecret)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errUnexpectedObject)
	}

	hash, time, err := repositorysecret.CreateOrUpdateSec(ctx, &cr.Spec.ForProvider, meta.GetExternalName(cr), e.client, e.gh)
	if err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, errUpdateRepositorySecrets)
	}

	cr.Status.AtProvider.LastUpdate = time
	cr.Status.AtProvider.EncryptValue = hash
	return managed.ExternalUpdate{}, nil
}

func (e *external) Delete(ctx context.Context, mgd resource.Managed) error {
	cr, ok := mgd.(*v1alpha1.RepositorySecret)
	if !ok {
		return errors.New(errUnexpectedObject)
	}

	_, err := e.gh.DeleteRepoSecret(ctx, cr.Spec.ForProvider.Owner, cr.Spec.ForProvider.Repository, meta.GetExternalName(cr))
	if err != nil {
		return errors.Wrap(err, errDeleteRepositorySecrets)
	}

	cr.Status.AtProvider.LastUpdate = ""
	cr.Status.AtProvider.EncryptValue = ""
	return nil
}
