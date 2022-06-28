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
	esv1 "github.com/elastic/cloud-on-k8s/pkg/apis/elasticsearch/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ElasticsearchHandler struct {
	ServiceHandlerObject
}

func (h *ElasticsearchHandler) ProvideConfig(values values.TemplateValues) []client.Object {
	var es esv1.Elasticsearch

	yaml.LoadYaml(filepath("elasticsearch/elasticsearch-deployment.yaml"), values, &es)

	h.Config = []client.Object{
		&es,
	}

	return h.Config
}
