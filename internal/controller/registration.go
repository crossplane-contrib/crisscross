/*
Copyright 2020 The Crossplane Authors.

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

package controller

import (
	"context"
	"strings"
	"time"

	"github.com/pkg/errors"
	kunstructured "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"
	kcontroller "sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/crossplane/crossplane-runtime/pkg/controller"
	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/crossplane/crossplane-runtime/pkg/resource/unstructured"
	rmanaged "github.com/crossplane/crossplane-runtime/pkg/resource/unstructured/managed"

	"github.com/crossplane/crossplane/pkg/controller/apiextensions/composite"
	"github.com/hasheddan/crisscross/apis/v1alpha1"
)

const (
	tinyWait  = 3 * time.Second
	shortWait = 30 * time.Second

	timeout        = 2 * time.Minute
	maxConcurrency = 5
	finalizer      = "defined.apiextensions.crossplane.io"

	errGetRegistration = "cannot get Registration"
	errStartController = "cannot start managed resource controller"
	errUpdateStatus    = "cannot update Registration status"
	errAddFinalizer    = "cannot add finalizer to Registration"
)

// Event reasons.
const (
	reasonStartController     event.Reason = "EstablishComposite"
	reasonTerminateController event.Reason = "TerminateManagedController"
)

// A controllerEngine can start and stop Kubernetes controllers on demand.
type controllerEngine interface {
	IsRunning(name string) bool
	Start(name string, o kcontroller.Options, w ...controller.Watch) error
	Stop(name string)
	Err(name string) error
}

// Setup adds a controller that reconciles Registrations.
func Setup(mgr ctrl.Manager, log logging.Logger) error {
	name := "crisscross/" + strings.ToLower(v1alpha1.RegistrationGroupKind)

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		For(&v1alpha1.Registration{}).
		WithOptions(kcontroller.Options{MaxConcurrentReconciles: maxConcurrency}).
		Complete(NewReconciler(mgr,
			WithLogger(log.WithValues("controller", name)),
			WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name)))))
}

// ReconcilerOption is used to configure the Reconciler.
type ReconcilerOption func(*Reconciler)

// WithLogger specifies how the Reconciler should log messages.
func WithLogger(log logging.Logger) ReconcilerOption {
	return func(r *Reconciler) {
		r.log = log
	}
}

// WithRecorder specifies how the Reconciler should record Kubernetes events.
func WithRecorder(er event.Recorder) ReconcilerOption {
	return func(r *Reconciler) {
		r.record = er
	}
}

// WithFinalizer specifies how the Reconciler should finalize
// CompositeResourceDefinitions.
func WithFinalizer(f resource.Finalizer) ReconcilerOption {
	return func(r *Reconciler) {
		r.finalizer = f
	}
}

// WithControllerEngine specifies how the Reconciler should manage the
// lifecycles of composite controllers.
func WithControllerEngine(c controllerEngine) ReconcilerOption {
	return func(r *Reconciler) {
		r.controller = c
	}
}

// WithClientApplicator specifies how the Reconciler should interact with the
// Kubernetes API.
func WithClientApplicator(ca resource.ClientApplicator) ReconcilerOption {
	return func(r *Reconciler) {
		r.client = ca
	}
}

// NewReconciler returns a Reconciler of Registrations.
func NewReconciler(mgr manager.Manager, opts ...ReconcilerOption) *Reconciler {
	kube := unstructured.NewClient(mgr.GetClient())

	r := &Reconciler{
		mgr: mgr,
		client: resource.ClientApplicator{
			Client:     kube,
			Applicator: resource.NewAPIUpdatingApplicator(kube),
		},
		controller: controller.NewEngine(mgr),
		finalizer:  resource.NewAPIFinalizer(kube, finalizer),
		log:        logging.NewNopLogger(),
		record:     event.NewNopRecorder(),
	}

	for _, f := range opts {
		f(r)
	}
	return r
}

// A Reconciler reconciles Registrations.
type Reconciler struct {
	client     resource.ClientApplicator
	mgr        manager.Manager
	controller controllerEngine
	finalizer  resource.Finalizer
	log        logging.Logger
	record     event.Recorder
}

// Reconcile a Registration by starting a controller that watches for its
// referenced type and calls out to its endpoint.
func (r *Reconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {

	log := r.log.WithValues("request", req)
	log.Debug("Reconciling")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	d := &v1alpha1.Registration{}
	if err := r.client.Get(ctx, req.NamespacedName, d); err != nil {
		log.Debug(errGetRegistration, "error", err)
		return reconcile.Result{}, errors.Wrap(resource.IgnoreNotFound(err), errGetRegistration)
	}

	log = log.WithValues(
		"uid", d.GetUID(),
		"version", d.GetResourceVersion(),
		"name", d.GetName(),
	)

	if meta.WasDeleted(d) {
		r.controller.Stop("managed/" + d.GetName())
		log.Debug("Stopped managed resource controller")
		r.record.Event(d, event.Normal(reasonTerminateController, "Stopped managed resource controller"))
		return reconcile.Result{RequeueAfter: tinyWait}, nil
	}

	if err := r.finalizer.AddFinalizer(ctx, d); err != nil {
		log.Debug(errAddFinalizer, "error", err)
		r.record.Event(d, event.Warning(reasonStartController, errors.Wrap(err, errAddFinalizer)))
		return reconcile.Result{RequeueAfter: shortWait}, nil
	}

	if err := r.controller.Err("managed/" + d.GetName()); err != nil {
		log.Debug("Managed resource controller encountered an error", "error", err)
	}

	// Build managed reconciler
	recorder := r.record.WithAnnotations("controller", composite.ControllerName(d.GetName()))
	managedGVK := schema.FromAPIVersionAndKind(d.Spec.TypeRef.APIVersion, d.Spec.TypeRef.Kind)
	uClient := unstructured.NewClient(r.mgr.GetClient())
	o := kcontroller.Options{Reconciler: managed.NewReconciler(r.mgr,
		resource.ManagedKind(managedGVK),
		func() resource.Managed {
			return rmanaged.New(rmanaged.WithGroupVersionKind(schema.GroupVersionKind(managedGVK)))
		},
		managed.WithClient(uClient),
		// TODO(hasheddan): handle connection publishing
		managed.WithFinalizer(resource.NewAPIFinalizer(uClient, "finalizer.managedresource.crisscross.crossplane.io")),
		managed.WithInitializers(
			managed.NewDefaultProviderConfig(uClient),
			managed.NewNameAsExternalName(uClient),
		),
		managed.WithExternalConnecter(&connector{endpoint: d.Spec.Endpoint}),
		managed.WithReferenceResolver(managed.NewAPISimpleReferenceResolver(uClient)),
		managed.WithLogger(log.WithValues("controller", "managed/"+d.GetName())),
		managed.WithRecorder(recorder))}

	u := &kunstructured.Unstructured{}
	u.SetGroupVersionKind(managedGVK)

	if err := r.controller.Start("managed/"+d.GetName(), o, controller.For(u, &handler.EnqueueRequestForObject{})); err != nil {
		log.Debug(errStartController, "error", err)
		r.record.Event(d, event.Warning(reasonStartController, errors.Wrap(err, errStartController)))
		return reconcile.Result{RequeueAfter: shortWait}, nil
	}

	r.record.Event(d, event.Normal(reasonStartController, "(Re)started managed resource controller"))
	return reconcile.Result{Requeue: false}, errors.Wrap(r.client.Status().Update(ctx, d), errUpdateStatus)
}
