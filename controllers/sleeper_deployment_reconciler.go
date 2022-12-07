package controllers

import (
	"context"

	"github.com/imdario/mergo"
	toolsv1alpha1 "github.com/opdev/simple-disconnected-operator/api/v1alpha1"
	"github.com/opdev/simple-disconnected-operator/image"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// SleeperDeploymentReconciler reconciles the deployment resource.
type SleeperDeploymentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=tools.opdev.io,resources=disconnectedfriendlyapps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=tools.opdev.io,resources=disconnectedfriendlyapps/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=tools.opdev.io,resources=disconnectedfriendlyapps/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;update;patch
//+kubebuilder:rbac:groups=apps,resources=deployments/finalizers,verbs=update

// Reconcile will ensure that the Sleeper Deployment for DisconnectedFriendlyApp
// reaches the desired state.
func (r *SleeperDeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx).WithValues("reconciler", "sleeper")
	l.Info("deployment reconciliation initiated.")
	defer l.Info("deployment reconciliation complete.")
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

	l.Info("Resolving Deployment Image", "image", image.Sleeper())
	new := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sleeper-" + instanceKey.Name,
			Namespace: instanceKey.Namespace,
			Labels: map[string]string{
				"app": "sleeper",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: instance.SleeperReplicas(),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "sleeper",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "sleeper-",
					Namespace:    instanceKey.Namespace,
					Labels: map[string]string{
						"app": "sleeper",
					},
					Annotations: map[string]string{},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "sleeper",
							Image:           image.Sleeper(),
							Command:         []string{"/bin/sleep", "infinity"},
							ImagePullPolicy: corev1.PullAlways,
						},
					},
				},
			},
		},
	}

	err = ctrl.SetControllerReference(&instance, &new, r.Scheme)
	if err != nil {
		return ctrl.Result{Requeue: true}, err
	}

	// If the deployment exists, get it and patch it
	var existing appsv1.Deployment
	err = r.Client.Get(ctx, client.ObjectKeyFromObject(&new), &existing)

	if apierrors.IsNotFound(err) {
		// create the resource because it does not exist.
		l.Info("creating resource", new.Kind, new.Name)
		if err := r.Client.Create(ctx, &new); err != nil {
			return ctrl.Result{Requeue: true}, err
		}
	}

	if err != nil {
		return ctrl.Result{Requeue: true}, err
	}

	l.Info("updating resources if necessary", existing.Kind, existing.GetName())
	patchDiff := client.MergeFrom(existing.DeepCopy())
	if err = mergo.Merge(&existing, new, mergo.WithOverride); err != nil {
		return ctrl.Result{Requeue: true}, err
	}

	if err = r.Patch(ctx, &existing, patchDiff); err != nil {
		return ctrl.Result{Requeue: true}, err
	}

	return ctrl.Result{Requeue: false}, nil // success
}

// SetupWithManager sets up the controller with the Manager.
func (r *SleeperDeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&toolsv1alpha1.DisconnectedFriendlyApp{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
