/*
Copyright 2026.

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
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	platformv1 "github.com/drabottini/cloud-native-platform/operator/api/v1"
)

// WebAppReconciler reconciles a WebApp object
type WebAppReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=platform.dylan.io,resources=webapps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=platform.dylan.io,resources=webapps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=platform.dylan.io,resources=webapps/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *WebAppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	webApp := &platformv1.WebApp{}
	err := r.Get(ctx, req.NamespacedName, webApp)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	logger.Info("Reconciling WebApp", "name", webApp.Name, "namespace", webApp.Namespace)

	if err := r.reconcileDeployment(ctx, webApp); err != nil {
		webApp.Status.Phase = "Error"
		webApp.Status.Message = fmt.Sprintf("Failed to reconcile Deployment: %v", err)
		_ = r.Status().Update(ctx, webApp)
		return ctrl.Result{}, err
	}

	/*
		if err := r.reconcileService(ctx, &webapp); err != nil {
			webapp.Status.Phase = "Error"
			webapp.Status.Message = fmt.Sprintf("Failed to reconcile Service: %v", err)
			_ = r.Status().Update(ctx, &webapp)
			return ctrl.Result{}, err
		}

		if err := r.updateStatus(ctx, &webapp); err != nil {
			return ctrl.Result{}, err
		} */

	return ctrl.Result{}, nil
}

func replicasOrDefault(webapp *platformv1.WebApp) int32 {
	if webapp.Spec.Replicas == nil {
		return 1
	}
	return *webapp.Spec.Replicas
}

func (r *WebAppReconciler) desiredDeployment(webapp *platformv1.WebApp) (*appsv1.Deployment, error) {
	replicas := replicasOrDefault(webapp)

	labels := map[string]string{
		"app.kubernetes.io/name":       webapp.Name,
		"app.kubernetes.io/managed-by": "webapp-operator",
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      webapp.Name,
			Namespace: webapp.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  webapp.Name,
							Image: webapp.Spec.Image,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: webapp.Spec.Port,
								},
							},
						},
					},
				},
			},
		},
	}

	if err := ctrl.SetControllerReference(webapp, deployment, r.Scheme); err != nil {
		return nil, err
	}

	return deployment, nil
}

func (r *WebAppReconciler) reconcileDeployment(ctx context.Context, webapp *platformv1.WebApp) error {
	desired, err := r.desiredDeployment(webapp)
	if err != nil {
		return err
	}

	existingDeploy := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: desired.Name, Namespace: desired.Namespace}, existingDeploy)

	if apierrors.IsNotFound(err) {
		return r.Create(ctx, desired)
	}

	if err != nil {
		return err
	}

	updated := existingDeploy.DeepCopy()
	updated.Labels = desired.Labels
	updated.Spec.Replicas = desired.Spec.Replicas
	updated.Spec.Selector = desired.Spec.Selector
	updated.Spec.Template = desired.Spec.Template

	if equality.Semantic.DeepEqual(existingDeploy.Spec, updated.Spec) &&
		equality.Semantic.DeepEqual(existingDeploy.Labels, updated.Labels) {
		return nil
	}

	return r.Update(ctx, updated)
}

// SetupWithManager sets up the controller with the Manager.
func (r *WebAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&platformv1.WebApp{}).
		Named("webapp").
		Complete(r)
}
