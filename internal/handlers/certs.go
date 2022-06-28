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
	v1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type CertHandler struct {
	ServiceHandlerObject
}

func (h *CertHandler) ProvideConfig(values values.TemplateValues) []client.Object {
	var ci v1.ClusterIssuer
	var cert v1.Certificate

	yaml.LoadYaml(opfilepath("cert-manager/cert-cluster-issuer.yaml"), values, &ci)
	yaml.LoadYaml(opfilepath("cert-manager/cert-primary.yaml"), values, &cert)

	h.Config = []client.Object{
		&ci,
		&cert,
	}

	return h.Config
}
