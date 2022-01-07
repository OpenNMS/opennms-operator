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
	"github.com/OpenNMS/opennms-operator/internal/model/values"
	ctrl "sigs.k8s.io/controller-runtime"
)

//UpdateValues - update values for an instance based on it's crd
func (r *OpenNMSReconciler) UpdateValues(req ctrl.Request) values.TemplateValues {
	if r.ValuesMap == nil {
		r.ValuesMap = map[string]values.TemplateValues{}
	}

	namespace := req.Name

	// get CRD
	// crd := getCrd(namespace)

	templateValues, ok := r.ValuesMap[namespace]
	if !ok {
		templateValues = r.DefaultValues
	}

	templateValues.Values.Namespace = namespace

	//other updating based on the CRD

	r.ValuesMap[namespace] = templateValues

	return templateValues
}