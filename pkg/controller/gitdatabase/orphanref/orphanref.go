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

package orphanref

import (
	"context"
	"fmt"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/crossplane/crossplane-runtime/pkg/ratelimiter"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/google/go-github/v48/github"
	"github.com/pkg/errors"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	"github.com/crossplane-contrib/provider-github/apis/gitdatabase/v1alpha1"
	ghclient "github.com/crossplane-contrib/provider-github/pkg/clients"
	"github.com/crossplane-contrib/provider-github/pkg/clients/orphanref"
)

const (
	errGetOrphanRef     = "cannot get GitHub repository orphanRef"
	errCreateOrphanRef  = "cannot create OrphanRef"
	errDeleteOrphanRef  = "cannot delete OrphanRef"
	errUnexpectedObject = "The managed resource is not a OrphanRef resource"
	errRepositoryEmpty  = "the repository value cannot be empty"
)

// SetupOrphanRef adds a controller that reconciles OrphanRefs.
func SetupOrphanRef(mgr ctrl.Manager, l logging.Logger, rl workqueue.RateLimiter) error {
	name := managed.ControllerName(v1alpha1.OrphanRefGroupKind)

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(controller.Options{
			RateLimiter: ratelimiter.NewDefaultManagedRateLimiter(rl),
		}).
		For(&v1alpha1.OrphanRef{}).
		Complete(managed.NewReconciler(mgr,
			resource.ManagedKind(v1alpha1.OrphanRefGroupVersionKind),
			managed.WithExternalConnecter(
				&orphanRefConnector{
					client:      mgr.GetClient(),
					newClientFn: orphanref.NewService,
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

type orphanRefConnector struct {
	client      client.Client
	newClientFn func(string) *orphanref.Service
}

func (c *orphanRefConnector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	cr, ok := mg.(*v1alpha1.OrphanRef)
	if !ok {
		return nil, errors.New(errUnexpectedObject)
	}
	cfg, err := ghclient.GetConfig(ctx, c.client, cr)
	if err != nil {
		return nil, err
	}

	return &orphanRefExternal{*c.newClientFn(string(cfg)), c.client}, nil
}

type orphanRefExternal struct {
	gh     orphanref.Service
	client client.Client
}

func (e *orphanRefExternal) Observe(ctx context.Context, mgd resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mgd.(*v1alpha1.OrphanRef)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errUnexpectedObject)
	}

	// Repository field is required
	if cr.Spec.ForProvider.Repository == nil {
		return managed.ExternalObservation{}, errors.New(errRepositoryEmpty)
	}

	refName := getRefName(cr)
	ref, res, err := e.gh.GetRef(ctx,
		cr.Spec.ForProvider.Owner,
		*cr.Spec.ForProvider.Repository,
		refName,
	)

	if err != nil {
		if res.StatusCode == 404 {
			return managed.ExternalObservation{ResourceExists: false}, nil
		}
		return managed.ExternalObservation{}, errors.Wrap(err, errGetOrphanRef)
	}

	cr.Status.AtProvider = orphanref.GenerateObservation(ref)
	cr.SetConditions(xpv1.Available())

	return managed.ExternalObservation{
		ResourceUpToDate: true,
		ResourceExists:   true,
	}, nil
}

func (e *orphanRefExternal) Create(ctx context.Context, mgd resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mgd.(*v1alpha1.OrphanRef)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errUnexpectedObject)
	}
	// Repository field is required
	if cr.Spec.ForProvider.Repository == nil {
		return managed.ExternalCreation{}, errors.New(errRepositoryEmpty)
	}

	tree, _, err := e.gh.CreateTree(ctx,
		cr.Spec.ForProvider.Owner,
		*cr.Spec.ForProvider.Repository,
		"",
		[]*github.TreeEntry{
			{
				Path:    &cr.Spec.ForProvider.Path,
				Mode:    ghclient.StringPtr("100644"),
				Type:    ghclient.StringPtr("blob"),
				Content: ghclient.StringPtr(""),
			},
		},
	)
	if err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, errCreateOrphanRef)
	}

	commit, _, err := e.gh.CreateCommit(ctx,
		cr.Spec.ForProvider.Owner,
		*cr.Spec.ForProvider.Repository,
		&github.Commit{
			Tree:    tree,
			Message: cr.Spec.ForProvider.Message,
		},
	)
	if err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, errCreateOrphanRef)
	}

	ref := getRefName(cr)
	_, _, err = e.gh.CreateRef(ctx,
		cr.Spec.ForProvider.Owner,
		*cr.Spec.ForProvider.Repository,
		&github.Reference{
			Ref: &ref,
			Object: &github.GitObject{
				SHA: commit.SHA,
			},
		},
	)

	if err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, errCreateOrphanRef)
	}

	return managed.ExternalCreation{}, err
}

func (e *orphanRefExternal) Update(ctx context.Context, mgd resource.Managed) (managed.ExternalUpdate, error) {
	return managed.ExternalUpdate{}, nil
}

func (e *orphanRefExternal) Delete(ctx context.Context, mgd resource.Managed) error {
	cr, ok := mgd.(*v1alpha1.OrphanRef)
	if !ok {
		return errors.New(errUnexpectedObject)
	}

	ref := getRefName(cr)
	_, err := e.gh.DeleteRef(ctx,
		cr.Spec.ForProvider.Owner,
		*cr.Spec.ForProvider.Repository,
		ref)

	return errors.Wrap(err, errDeleteOrphanRef)
}

func getRefName(cr *v1alpha1.OrphanRef) string {
	name := meta.GetExternalName(cr)
	return fmt.Sprintf("refs/heads/%s", name)
}
