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

type OpenNMSValues struct {
	Core ServiceValues `yaml:"Core"`
	API  ServiceValues `yaml:"API"`
	UI   ServiceValues `yaml:"UI"`
}

type ServiceValues struct {
	ServiceName        string         `yaml:"ServiceName"`
	Image              string         `yaml:"Image"`
	VolumeSize         string         `yaml:"VolumeSize"`
	InitContainerImage string         `yaml:"InitContainerImage"`
	Resources          ResourceValues `yaml:"Resources"`
}

type ResourceValues struct {
	Limits   ResourceDefinition `yaml:"Limits"`
	Requests ResourceDefinition `yaml:"Requests"`
}

type ResourceDefinition struct {
	Cpu    string `yaml:"Cpu"`
	Memory string `yaml:"Memory"`
}

type TimeseriesValues struct {
	Mode   string `yaml:"Mode"`
	Host   string `yaml:"Host"`
	Port   string `yaml:"Port"`
	ApiKey string `yaml:"ApiKey"`
}
