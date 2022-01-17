/*
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

package reconciler

import (
	"context"
	"github.com/OpenNMS/opennms-operator/internal/util/subsets"
	v1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"
)

func (r *OpenNMSReconciler) updateDeployment(ctx context.Context, resource client.Object, deployedResource client.Object) (reconcile.Result, error) {
	if !subsets.SubsetEqual(resource, deployedResource) {
		if err := r.Update(ctx, resource); err != nil {
			return reconcile.Result{}, err
		}
		return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
	} else {
		// Determine if the resources are fully created, otherwise wait longer
		deployment := deployedResource.(*v1.Deployment)
		if deployment.Status.ReadyReplicas != deployment.Status.Replicas {
			return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
		}
	}
	return reconcile.Result{}, nil
}

func (r *OpenNMSReconciler) updateStatefulSet(ctx context.Context, resource client.Object, deployedResource client.Object) (reconcile.Result, error) {
	if !subsets.SubsetEqual(resource, deployedResource) {
		if err := r.Update(ctx, resource); err != nil {
			return reconcile.Result{}, err
		}
		return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
	} else {
		// Determine if the resources are fully created, otherwise wait longer
		statefulset := deployedResource.(*v1.StatefulSet)
		if statefulset.Status.ReadyReplicas != statefulset.Status.Replicas {
			return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
		}
	}
	return reconcile.Result{}, nil
}

func (r *OpenNMSReconciler) updateJob(ctx context.Context, resource client.Object, deployedResource client.Object) (reconcile.Result, error) {
	job := deployedResource.(*batchv1.Job)
	if job.Status.Succeeded < 1 {
		return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
	}
	return reconcile.Result{}, nil
}

func (r *OpenNMSReconciler) updateSecret(ctx context.Context, resource client.Object, deployedResource client.Object) (reconcile.Result, error) {
	if resource.GetName() == "onms-allowed-users" && !subsets.SubsetEqual(resource, deployedResource) {
		err := r.updateAllowedUsersSecret(ctx, resource)
		if err != nil {
			return reconcile.Result{}, err
		}
	} else if !subsets.SubsetEqual(resource, deployedResource) {
		if err := r.Update(ctx, resource); err != nil {
			return reconcile.Result{}, err
		}
	}
	return reconcile.Result{}, nil
}

func (r *OpenNMSReconciler) updateConfigMap(ctx context.Context, resource client.Object, deployedResource client.Object) (reconcile.Result, error) {
	if !subsets.SubsetEqual(resource, deployedResource) {
		if err := r.Update(ctx, resource); err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
	}
	return reconcile.Result{}, nil
}

func (r *OpenNMSReconciler) updateAllowedUsersSecret(ctx context.Context, resource client.Object) error {
	// Need to update and restart auth container
	if err := r.Update(ctx, resource); err != nil {
		return err
	}
	// Update deployment to restart the auth container
	deployment := &v1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{Name: "onms-auth", Namespace: resource.GetNamespace()}, deployment)
	if err != nil {
		return err
	} else {
		// Update deployment to cause a restart. Needs to be something in the spec.template section for a restart
		if deployment.Spec.Template.Annotations == nil {
			deployment.Spec.Template.Annotations = make(map[string]string)
		}
		deployment.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().String()
		if err := r.Update(ctx, deployment); err != nil {
			return err
		}
	}
	return nil
}