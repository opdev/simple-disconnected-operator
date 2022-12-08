package controllers

import (
	"context"

	toolsv1alpha1 "github.com/opdev/simple-disconnected-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// StatusReconciler reconciles the DisconnectedFriendlyApp Status
type StatusReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=tools.opdev.io,resources=disconnectedfriendlyapps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=tools.opdev.io,resources=disconnectedfriendlyapps/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=tools.opdev.io,resources=disconnectedfriendlyapps/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments/finalizers,verbs=update

// Reconcile will ensure that the status for DisconnectedFriendlyApp.
func (r *StatusReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx).WithValues("reconciler", "status")
	l.Info("status reconciliation initiated.")
	defer l.Info("status reconciliation complete.")
	instanceKey := req.NamespacedName

	// Get the DisconnectedFriendlyApp instance to make sure it still exists.
	var instance toolsv1alpha1.DisconnectedFriendlyApp
	err := r.Client.Get(ctx, instanceKey, &instance)

	if apierrors.IsNotFound(err) {
		return ctrl.Result{Requeue: false}, nil
	}

	if err != nil {
		return ctrl.Result{Requeue: true}, err
	}

	patchDiff := client.MergeFrom(instance.DeepCopy())
	instance.PopulateStatus()

	if err = r.Status().Patch(ctx, &instance, patchDiff); err != nil {
		return ctrl.Result{Requeue: true}, err
	}

	return ctrl.Result{Requeue: false}, nil // success
}

// SetupWithManager sets up the controller with the Manager.
func (r *StatusReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&toolsv1alpha1.DisconnectedFriendlyApp{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
