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

package values

type TemplateValues struct {
	Values Values
}

type Values struct {
	Namespace  string         `yaml:"Namespace"`
	Host       string         `yaml:"Host"`
	TestDeploy bool           `yaml:"TestDeploy"`
	OpenNMS    OpenNMSValues  `yaml:"OpenNMS"`
	TLS        TLSValues      `yaml:"TLS"`
	Auth       AuthValues     `yaml:"Auth"`
	Postgres   PostgresValues `yaml:"Postgres"`
	Grafana    GrafanaValues  `yaml:"Grafana"`
	Ingress    IngressValues  `yaml:"Ingress"`
	Stunnel    StunnelValues  `yaml:"Stunnel"`
}
