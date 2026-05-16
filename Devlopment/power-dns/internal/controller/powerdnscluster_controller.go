package controller

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	platformv1alpha1 "github.com/example/powerdns-platform/api/v1alpha1"
)

type PowerDNSClusterReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

func NewPowerDNSClusterReconciler(cli client.Client, scheme *runtime.Scheme, recorder record.EventRecorder) *PowerDNSClusterReconciler {
	return &PowerDNSClusterReconciler{Client: cli, Scheme: scheme, Recorder: recorder}
}

func (r *PowerDNSClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	instance := &platformv1alpha1.PowerDNSCluster{}
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if err := r.reconcileWorkloads(ctx, instance); err != nil {
		logger.Error(err, "unable to reconcile workloads")
		return ctrl.Result{RequeueAfter: 20 * time.Second}, err
	}

	instance.Status.ObservedGeneration = instance.Generation
	instance.Status.Phase = platformv1alpha1.PowerDNSClusterPhaseReconciling
	instance.Status.ReadyReplicas = r.readyReplicas(ctx, instance)
	instance.Status.Endpoints = platformv1alpha1.PowerDNSEndpointsStatus{
		API:        instance.Name + "-powerdns." + instance.Namespace + ".svc.cluster.local",
		DNS:        instance.Name + "-powerdns." + instance.Namespace + ".svc.cluster.local",
		PostgreSQL: instance.Name + "-postgresql." + instance.Namespace + ".svc.cluster.local",
	}
	if instance.Status.ReadyReplicas >= r.desiredReadyReplicas(instance) {
		instance.Status.Phase = platformv1alpha1.PowerDNSClusterPhaseReady
	}

	if err := r.Status().Update(ctx, instance); err != nil {
		logger.Error(err, "unable to update status")
		return ctrl.Result{RequeueAfter: 20 * time.Second}, err
	}

	return ctrl.Result{RequeueAfter: 60 * time.Second}, nil
}

func (r *PowerDNSClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&platformv1alpha1.PowerDNSCluster{}).
		WithEventFilter(predicate.Funcs{
			CreateFunc: func(e event.CreateEvent) bool { return true },
			UpdateFunc: func(e event.UpdateEvent) bool { return true },
			DeleteFunc: func(e event.DeleteEvent) bool { return true },
			GenericFunc: func(e event.GenericEvent) bool { return true },
		}).
		Complete(r)
}