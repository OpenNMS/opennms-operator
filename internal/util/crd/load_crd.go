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

package crd

import (
	"context"
	"github.com/OpenNMS/opennms-operator/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetCRDFromCluster(ctx context.Context, k8sClient client.Client, req ctrl.Request) (v1alpha1.OpenNMS, error) {
	var instanceCRD v1alpha1.OpenNMS
	err := k8sClient.Get(ctx, req.NamespacedName, &instanceCRD)
	if err != nil {
		return v1alpha1.OpenNMS{}, err
	}
	return instanceCRD, nil
}