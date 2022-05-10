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

package image

import (
	"context"
	"github.com/OpenNMS/opennms-operator/api/v1alpha1"
	"github.com/OpenNMS/opennms-operator/internal/util/crd"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"strings"
	"time"
)

const (
	None = "none"
	Now  = "now"
)

//markInstanceServiceForUpdate - mark a given service in given instance as having an update available
func (iu *ImageUpdater) markInstanceServiceForUpdate(ctx context.Context, instanceName, serviceName, oldDigest, newDigest string) {
	instance, err := crd.GetInstance(ctx, iu.Client, types.NamespacedName{Name: instanceName})
	if err != nil {
		if errors.IsNotFound(err) {
			iu.Log.Info("OpenNMS resource not found", "name", instanceName)
			return
		}
		// Error reading the object - requeue the request.
		iu.Log.Error(err, "Failed to get OpenNMS", "name", instanceName)
		return
	}
	now := time.Now().Format(time.RFC3339)

	isLatest := oldDigest == newDigest
	instance.Status.Image.IsLatest = isLatest
	instance.Status.Image.CheckedAt = now
	if !isLatest {
		if instance.Status.Image.ServicesToUpdate == "" {
			instance.Status.Image.ServicesToUpdate = serviceName
		} else {
			instance.Status.Image.ServicesToUpdate = instance.Status.Image.ServicesToUpdate + "," + serviceName
		}
	}
	if err := iu.Status().Update(ctx, &instance); err != nil {
		iu.Log.Error(err, "Failed to update OpenNMS status", "name", instanceName)
	}
	return
}

//UpdateServices - update the services in a given instance
func (iu *ImageUpdater) UpdateServices(instance v1alpha1.OpenNMS) {
	if instance.Spec.ImageUpdateConfig.Update == Now {
		ctx := context.Background()
		iu.forceUpdateServices(ctx, instance.GetName(), instance.Status.Image.ServicesToUpdate)
		instance.Spec.ImageUpdateConfig.Update = None
		if err := iu.Client.Update(ctx, &instance); err != nil {
			iu.Log.Error(err, "Failed to update OpenNMS image.update")
		}
		instance.Status.Image.IsLatest = true
		instance.Status.Image.ServicesToUpdate = ""
		if err := iu.Client.Status().Update(ctx, &instance); err != nil {
			iu.Log.Error(err, "Failed to update status after image update")
		}
	}
}

//forceUpdateServices - force update a list of services for an instance
func (iu *ImageUpdater) forceUpdateServices(ctx context.Context, instanceName, servicesToUpdate string) {
	serviceList := strings.Split(servicesToUpdate, ",")
	for _, service := range serviceList {
		iu.forceUpdateService(ctx, instanceName, service)
	}
}

//forceUpdateService - force update a given service for an instance
func (iu *ImageUpdater) forceUpdateService(ctx context.Context, instanceName, serviceName string) {
	if success := iu.forceUpdateStatefulSet(ctx, instanceName, serviceName); success {
		iu.Log.Info("StatefulSet in instance update", "instance", instanceName, "service", serviceName)
		return
	} else if success := iu.forceUpdateDeployment(ctx, instanceName, serviceName); success {
		iu.Log.Info("Deployment in instance update", "instance", instanceName, "service", serviceName)
		return
	} else {
		iu.Log.Info("Service could not be found in instance", "instance", instanceName, "service", serviceName)
	}
}

//forceUpdateStatefulSet - force update a statefulset in a given instance
func (iu *ImageUpdater) forceUpdateStatefulSet(ctx context.Context, instanceName, serviceName string) bool {
	var statefulSet v1.StatefulSet
	err := iu.Client.Get(ctx, types.NamespacedName{Namespace: instanceName, Name: serviceName}, &statefulSet)
	if err != nil {
		if errors.IsNotFound(err) {
			iu.Log.Info("Deployment resource not found", "instance", instanceName, "service", serviceName)
		} else {
			iu.Log.Error(err, "Error getting deployment resource", "instance", instanceName, "service", serviceName)
		}
		return false
	}
	an := statefulSet.Spec.Template.Annotations
	if an == nil {
		an = map[string]string{}
	}
	an["lastUpdate"] = time.Now().Format(time.RFC3339)
	statefulSet.Spec.Template.Annotations = an
	err = iu.Client.Update(ctx, &statefulSet)
	if err != nil {
		return false
	}
	return true
}

//forceUpdateDeployment - force update a deployment in a given instance
func (iu *ImageUpdater) forceUpdateDeployment(ctx context.Context, instanceName, serviceName string) bool {
	var deployment v1.Deployment
	err := iu.Client.Get(ctx, types.NamespacedName{Namespace: instanceName, Name: serviceName}, &deployment)
	if err != nil {
		if errors.IsNotFound(err) {
			iu.Log.Info("Deployment resource not found", "instance", instanceName, "service", serviceName)
		} else {
			iu.Log.Error(err, "Error getting deployment resource", "instance", instanceName, "service", serviceName)
		}
		return false
	}
	an := deployment.Spec.Template.Annotations
	if an == nil {
		an = map[string]string{}
	}
	an["lastUpdate"] = time.Now().Format(time.RFC3339)
	deployment.Spec.Template.Annotations = an
	err = iu.Client.Update(ctx, &deployment)
	if err != nil {
		return false
	}
	return true
}
