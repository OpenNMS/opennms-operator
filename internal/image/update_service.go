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
	"github.com/OpenNMS/opennms-operator/internal/util/crd"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"time"
)

const (
	None = "none"
	Now  = "now"
)

//markInstanceServiceForUpdate - mark a given service in given instance as having an update available
func (ic *ImageChecker) markInstanceServiceForUpdate(ctx context.Context, instanceName, serviceName, oldDigest, newDigest string) {
	instance, err := crd.GetInstance(ctx, ic.Client, types.NamespacedName{Name: instanceName})
	if err != nil {
		if errors.IsNotFound(err) {
			ic.Log.Info("OpenNMS resource not found", "name", instanceName)
			return
		}
		// Error reading the object - requeue the request.
		ic.Log.Error(err, "Failed to get OpenNMS", "name", instanceName)
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
	if err := ic.Status().Update(ctx, &instance); err != nil {
		ic.Log.Error(err, "Failed to update OpenNMS status", "name", instanceName)
	}
	return
}
