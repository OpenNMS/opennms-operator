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
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

func (ic *ImageChecker) getInstance(ctx context.Context, instance string) v1alpha1.OpenNMS {
	ns := types.NamespacedName{Name: instance}
	ic.Log.Info("Searching for OpenNMS instance", "name", ns.Name)

	var onmsInstance v1alpha1.OpenNMS
	err := ic.Client.Get(ctx, ns, &onmsInstance)
	if err != nil {
		if errors.IsNotFound(err) {
			ic.Log.Info("OpenNMS resource not found")
			return onmsInstance
		}
		// Error reading the object - requeue the request.
		ic.Log.Error(err, "Failed to get OpenNMS")
		return onmsInstance
	}
	return onmsInstance
}
