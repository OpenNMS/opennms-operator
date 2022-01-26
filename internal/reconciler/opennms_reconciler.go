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
	"fmt"
	"github.com/OpenNMS/opennms-operator/api/v1alpha1"
	"github.com/OpenNMS/opennms-operator/config"
	"github.com/OpenNMS/opennms-operator/internal/handlers"
	"github.com/OpenNMS/opennms-operator/internal/model/values"
	"github.com/OpenNMS/opennms-operator/internal/util/crd"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"
)

// OpenNMSReconciler reconciles a OpenNMS object
type OpenNMSReconciler struct {
	client.Client
	Log           logr.Logger
	Scheme        *runtime.Scheme
	CodecFactory  serializer.CodecFactory
	Config        config.OperatorConfig
	DefaultValues values.TemplateValues
	Handlers      []handlers.ServiceHandler
	ValuesMap     map[string]values.TemplateValues
}

func (r *OpenNMSReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	instance, err := crd.GetCRDFromCluster(ctx, r.Client, req)
	if err != nil && errors.IsNotFound(err) {
		r.Log.Info("queued instance does not exist " + req.Name)
		return ctrl.Result{}, err
	} else if err != nil {
		r.Log.Error(err, "failed to get crd for instance "+req.Name)
		return ctrl.Result{}, err
	}
	valuesForInstance := r.UpdateValues(ctx, instance)

	for _, handler := range r.Handlers {
		for _, resource := range handler.ProvideConfig(valuesForInstance) {
			kind := reflect.ValueOf(resource).Elem().Type().String()
			deployedResource, exists := r.getResourceFromCluster(ctx, resource)
			if !exists {
				r.updateStatus(ctx, &instance, false, "instance starting")
				r.Log.Info("creating resource", "name", resource.GetName(), "namespace", resource.GetNamespace(), "kind", kind)
				err := r.Create(ctx, resource)
				if err != nil {
					r.Log.Error(err, "error creating resource", resource.GetName(), "namespace", resource.GetNamespace(), "kind", kind, "error", err)
					return ctrl.Result{}, err
				}
				if kind == "v1.Deployment" || kind == "v1.Job" || kind == "v1.StatefulSet" {
					return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
				}
			} else {
				r.Log.Info("updating resource", "name", resource.GetName(), "namespace", resource.GetNamespace(), "kind", kind)
				r.updateStatus(ctx, &instance, false, "updating instance resource")
				var res *reconcile.Result
				var err error
				switch kind {
				case "v1.Deployment":
					res, err = r.updateDeployment(ctx, resource, deployedResource)
				case "v1.StatefulSet":
					res, err = r.updateStatefulSet(ctx, resource, deployedResource)
				case "v1.Job":
					res, err = r.updateJob(ctx, resource, deployedResource)
				case "v1.Secret":
					res, err = r.updateSecret(ctx, resource, deployedResource)
				case "v1.ConfigMap":
					res, err = r.updateConfigMap(ctx, resource, deployedResource)
				}
				if err != nil {
					r.Log.Info("error updating resource", "name", resource.GetName(), "namespace", resource.GetNamespace(), "kind", kind, "error", err)
					r.updateStatus(ctx, &instance, false, fmt.Sprintf("Error: failed to update resource: %s %s %s", resource.GetNamespace(), kind, resource.GetName()))
					return reconcile.Result{}, err
				}
				if res != nil {
					return *res, nil
				}
			}
		}
	}
	return ctrl.Result{}, nil
}

func (r *OpenNMSReconciler) getResourceFromCluster(ctx context.Context, resource client.Object) (client.Object, bool) {
	r.Log.Info("Attempting get for", "name", resource.GetName(), "namespace", resource.GetNamespace())
	deployedResource := resource.DeepCopyObject().(client.Object)
	err := r.Get(ctx, types.NamespacedName{Name: resource.GetName(), Namespace: resource.GetNamespace()}, deployedResource)
	return deployedResource, !(err != nil && errors.IsNotFound(err))
}

func (r *OpenNMSReconciler) updateStatus(ctx context.Context, instance *v1alpha1.OpenNMS, ready bool, reason string) {
	if instance.Status.Ready != ready || instance.Status.Reason != reason {
		instance.Status.Ready = ready
		instance.Status.Reason = reason
		_ = r.Status().Update(ctx, instance)
	}
}
