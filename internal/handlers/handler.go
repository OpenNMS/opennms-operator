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
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var ConfigFilePath = "./charts/opennms/templates/"
var OperatorFilePath = "./charts/opennms-operator/templates/"

type ServiceHandler interface {
	//ProvideConfig - provide k8s config for the service
	ProvideConfig(values values.TemplateValues) []client.Object

	//HandleUpdate - handle the update process for the service ???
	//HandleUpdate(values values.TemplateValues) ???
}

type ServiceHandlerObject struct {
	Config []client.Object
}

func filepath(filename string) string {
	return ConfigFilePath + filename
}

func opfilepath(filename string) string {
	return OperatorFilePath + filename
}
