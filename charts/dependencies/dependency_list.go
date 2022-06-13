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

package dependencies

import (
	"github.com/OpenNMS/opennms-operator/internal/handlers"
	helmclient "github.com/mittwald/go-helm-client"
	"helm.sh/helm/v3/pkg/repo"
)

var repositories = []repo.Entry{
	{
		Name: "elastic",
		URL:  "https://helm.elastic.co",
	},
	{
		Name: "jetstack",
		URL:  "https://charts.jetstack.io",
	},
	{
		Name: "mittwald",
		URL:  "https://helm.mittwald.de",
	},
	{
		Name: "opennms",
		URL:  "https://opennms.github.io/horizon-stream/charts/packaged",
	},
}

var charts = []helmclient.ChartSpec{
	{
		ChartName:   "elastic/eck-operator",
		ReleaseName: "elastic-system",
		Namespace:   "elastic-system",

		CreateNamespace: true,
	},
	{
		ChartName:   "jetstack/cert-manager",
		ReleaseName: "cert-manager",
		Namespace:   "cert-manager",

		//--set installCRDs=true
		ValuesYaml: "installCRDs: true",

		Version:         "v1.7.0",
		CreateNamespace: true,
	},
	{
		ChartName:   "mittwald/kubernetes-replicator",
		ReleaseName: "kubernetes-replicator",
		Namespace:   "kubernetes-replicator",
	},
	{
		ChartName:   "opennms/onms-kafka",
		ReleaseName: "kafka",
		Namespace:   "kafka",
	},
	{
		ChartName:   "opennms/onms-keycloak",
		ReleaseName: "keycloak",
		Namespace:   "keycloak",
	},
}

var handlerslist = []handlers.ServiceHandler{
	&handlers.CertHandler{},
}
