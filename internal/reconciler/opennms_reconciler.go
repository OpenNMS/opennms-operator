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
	"github.com/OpenNMS/opennms-operator/config"
	"github.com/OpenNMS/opennms-operator/internal/handlers"
	"github.com/OpenNMS/opennms-operator/internal/model/values"
	"github.com/OpenNMS/opennms-operator/internal/util/crd"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
	ValuesMap	  map[string]values.TemplateValues
}

func (r *OpenNMSReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	instanceCRD, err := crd.GetCRDFromCluster(ctx, r.Client, req)
	if err != nil {
		r.Log.Error(err, "failed to get crd for instance " + req.Name)
		return ctrl.Result{}, err
	}
	valuesForInstance := r.UpdateValues(ctx, instanceCRD)

	for _, handler := range r.Handlers {
		for _, resource := range handler.ProvideConfig(valuesForInstance) {
			err := r.Create(ctx, resource)
			if err != nil {
				r.Log.Error(err, "error creating resource")
				return ctrl.Result{}, err
			}
		}
	}

	return ctrl.Result{}, nil
}


