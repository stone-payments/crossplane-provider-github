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
package content

import (
	"context"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/ratelimiter"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/google/go-github/v48/github"
	"github.com/pkg/errors"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	"github.com/crossplane-contrib/provider-github/apis/repositories/v1alpha1"
	ghclient "github.com/crossplane-contrib/provider-github/pkg/clients"
	"github.com/crossplane-contrib/provider-github/pkg/clients/content"
)

const (
	disabledReconcile   = "Disabled"
	errGetContent       = "cannot get GitHub repository content"
	errCreateContent    = "cannot create Content"
	errUpdateContent    = "cannot update Content"
	errDeleteContent    = "cannot delete Content"
	errUnexpectedObject = "The managed resource is not a Content resource"
	errRepositoryEmpty  = "the repository value cannot be empty"
	errCheckUpToDate    = "unable to determine if external resource is up to date"
)

// SetupContent adds a controller that reconciles Repositories.
func SetupContent(mgr ctrl.Manager, l logging.Logger, rl workqueue.RateLimiter) error {
	name := managed.ControllerName(v1alpha1.ContentGroupKind)

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(controller.Options{
			RateLimiter: ratelimiter.NewDefaultManagedRateLimiter(rl),
		}).
		For(&v1alpha1.Content{}).
		Complete(managed.NewReconciler(mgr,
			resource.ManagedKind(v1alpha1.ContentGroupVersionKind),
			managed.WithExternalConnecter(
				&contentConnector{
					client:      mgr.GetClient(),
					newClientFn: content.NewService,
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

type contentConnector struct {
	client      client.Client
	newClientFn func(string) (*content.Service, error)
}

func (c *contentConnector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	cr, ok := mg.(*v1alpha1.Content)
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

	return &contentExternal{
		gh:     *client,
		client: c.client,
	}, nil
}

type contentExternal struct {
	gh     content.Service
	client client.Client
}

func (e *contentExternal) Observe(ctx context.Context, mgd resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mgd.(*v1alpha1.Content)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errUnexpectedObject)
	}

	if skipReconcile(cr) {
		return reconciledObservation(), nil
	}

	var opts *github.RepositoryContentGetOptions
	if cr.Spec.ForProvider.Branch != nil {
		opts = &github.RepositoryContentGetOptions{
			Ref: *cr.Spec.ForProvider.Branch,
		}
	}

	if cr.Spec.ForProvider.Repository == nil {
		return managed.ExternalObservation{}, errors.New(errRepositoryEmpty)
	}

	fc, _, res, err := e.gh.GetContents(ctx,
		cr.Spec.ForProvider.Owner,
		*cr.Spec.ForProvider.Repository,
		cr.Spec.ForProvider.Path,
		opts,
	)

	if err != nil {
		if res.StatusCode == 404 {
			return managed.ExternalObservation{}, nil
		}
		return managed.ExternalObservation{}, errors.Wrap(err, errGetContent)
	}

	cr.Status.AtProvider = content.GenerateObservation(*fc)
	cr.SetConditions(xpv1.Available())

	upToDate, err := content.IsUpToDate(&cr.Spec.ForProvider, fc)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errCheckUpToDate)
	}

	return managed.ExternalObservation{
		ResourceUpToDate: upToDate,
		ResourceExists:   true,
	}, nil
}

func (e *contentExternal) Create(ctx context.Context, mgd resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mgd.(*v1alpha1.Content)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errUnexpectedObject)
	}

	if cr.Spec.ForProvider.Repository == nil {
		return managed.ExternalCreation{}, errors.New(errRepositoryEmpty)
	}

	_, _, err := e.gh.CreateFile(ctx,
		cr.Spec.ForProvider.Owner,
		*cr.Spec.ForProvider.Repository,
		cr.Spec.ForProvider.Path,
		&github.RepositoryContentFileOptions{
			Message: &cr.Spec.ForProvider.Message,
			Branch:  cr.Spec.ForProvider.Branch,
			Content: []byte(cr.Spec.ForProvider.Content),
		},
	)

	if err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, errCreateContent)
	}

	return managed.ExternalCreation{}, err
}

func (e *contentExternal) Update(ctx context.Context, mgd resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mgd.(*v1alpha1.Content)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errUnexpectedObject)
	}

	_, _, err := e.gh.UpdateFile(ctx,
		cr.Spec.ForProvider.Owner,
		*cr.Spec.ForProvider.Repository,
		cr.Spec.ForProvider.Path,
		&github.RepositoryContentFileOptions{
			Message: &cr.Spec.ForProvider.Message,
			Branch:  cr.Spec.ForProvider.Branch,
			Content: []byte(cr.Spec.ForProvider.Content),
			SHA:     &cr.Status.AtProvider.SHA,
		},
	)
	return managed.ExternalUpdate{}, errors.Wrap(err, errUpdateContent)
}

func (e *contentExternal) Delete(ctx context.Context, mgd resource.Managed) error {
	cr, ok := mgd.(*v1alpha1.Content)
	if !ok {
		return errors.New(errUnexpectedObject)
	}

	_, _, err := e.gh.DeleteFile(ctx,
		cr.Spec.ForProvider.Owner,
		*cr.Spec.ForProvider.Repository,
		cr.Spec.ForProvider.Path,
		&github.RepositoryContentFileOptions{
			Message: &cr.Spec.ForProvider.Message,
			Branch:  cr.Spec.ForProvider.Branch,
			Content: []byte(cr.Spec.ForProvider.Content),
			SHA:     &cr.Status.AtProvider.SHA,
		},
	)
	return errors.Wrap(err, errDeleteContent)
}

func skipReconcile(cr *v1alpha1.Content) bool {
	return cr.GetCondition(xpv1.TypeReady).Reason == xpv1.Available().Reason && *cr.Spec.Reconcile == disabledReconcile
}

func reconciledObservation() managed.ExternalObservation {
	return managed.ExternalObservation{
		ResourceExists:   true,
		ResourceUpToDate: true,
	}
}
