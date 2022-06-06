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

package handlers

import (
	"github.com/OpenNMS/opennms-operator/internal/model/values"
	"github.com/OpenNMS/opennms-operator/internal/util/yaml"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type OpenNMSHandler struct {
	ServiceHandlerObject
}

func (h *OpenNMSHandler) ProvideConfig(values values.TemplateValues) []client.Object {
	var configMap corev1.ConfigMap
	var coreService corev1.Service
	var coreDeployment appsv1.Deployment
	var apiService corev1.Service
	var apiDeployment appsv1.Deployment
	var uiDeployment appsv1.Deployment
	var uiService corev1.Service

	yaml.LoadYaml(filepath("opennms/opennms-configmap.yaml"), values, &configMap)
	yaml.LoadYaml(filepath("opennms/opennms-core-service.yaml"), values, &coreService)
	yaml.LoadYaml(filepath("opennms/opennms-core-deployment.yaml"), values, &coreDeployment)
	yaml.LoadYaml(filepath("opennms/opennms-api-service.yaml"), values, &apiService)
	yaml.LoadYaml(filepath("opennms/opennms-api-deployment.yaml"), values, &apiDeployment)
	yaml.LoadYaml(filepath("opennms/opennms-ui-deployment.yaml"), values, &uiDeployment)
	yaml.LoadYaml(filepath("opennms/opennms-ui-service.yaml"), values, &uiService)

	h.Config = []client.Object{
		&configMap,
		&coreService,
		&coreDeployment,
		&apiService,
		&apiDeployment,
		&uiDeployment,
		&uiService,
	}

	return h.Config
}
