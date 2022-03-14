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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// OpenNMSSpec defines the desired state of OpenNMS
type OpenNMSSpec struct {
	// Version of OpenNMS.
	Version Version `json:"version,omitempty"`

	// Domain name used in ingress rule
	Host string `json:"host,omitempty"`

	// K8s namespace to use
	Namespace string `json:"namespace"`

	// Users allowed login via tokens
	AllowedUsers []string `json:"allowedUsers,omitempty"`

	// Deploy an instance in a nonoperative testing mode
	TestDeploy bool `json:"testDeploy,omitempty"`

	// Defines what plugin for timeseries to use
	Timeseries Timeseries `json:"timeseries,omitempty"`

	// Defines cpu,mem and disk size for core
	Core BaseServiceResources `json:"core,omitempty"`

	// Defines cpu,mem and disk size for postgres
	Postgres BaseServiceResources `json:"postgres,omitempty"`

	// Defines the logic of ONMS image update
	ImageUpdateConfig ImageUpdateConfig `json:"imageUpdate,omitempty"`
}

//Version - defines the version of the ONMS core to use
type Version struct {
	Distribution string `json:"distribution"`
	Tag          string `json:"tag"`
}

//Timeseries - defines the timeseries DB backend to use
type Timeseries struct {
	Mode   string `json:"mode,omitempty"`
	Host   string `json:"host,omitempty"`
	Port   string `json:"port,omitempty"`
	ApiKey string `json:"apiKey,omitempty"`
}

//BaseServiceResources - defines basic resource needs of a service
type BaseServiceResources struct {
	MEM  string `json:"mem,omitempty"`
	Disk string `json:"disk,omitempty"`
	CPU  string `json:"cpu,omitempty"`
}

// OpenNMSStatus - defines the observed state of OpenNMS
type OpenNMSStatus struct {
	Image     ImageStatus     `json:"image,omitempty"`
	Readiness ReadinessStatus `json:"readiness,omitempty"`
	Nodes     []string        `json:"nodes,omitempty"`
}

// ImageUpdateConfig - defines current status of used image for OpenNMS container
type ImageUpdateConfig struct {
	// can have values of now/autoupdate/none
	Update string `json:"update,omitempty"`
	// represents number of minutes for recurrent checks of a new image
	Frequency int `json:"frequency,omitempty"`
}

// ImageStatus - defines current status of used image for OpenNMS container
type ImageStatus struct {
	// true if latest image used, false otherwise
	IsLatest bool `json:"isLatest"`
	// timestamp of a last image check in DockerHub
	CheckedAt string `json:"checkedAt,omitempty"`
	// readable message about image status
	Message string `json:"message,omitempty"`
}

//ReadinessStatus - the ready status of the ONMS instance
type ReadinessStatus struct {
	// if the ONMS instance is ready
	Ready bool `json:"ready,omitempty"`
	// reason an ONMS instance isn't ready
	Reason string `json:"reason,omitempty"`
	// the time the `ready` flag was last updated
	Timestamp string `json:"timestamp,omitempty"`
	// list of readinesses of the constituent services
	Services []ServiceStatus `json:"services,omitempty"`
}

type ServiceStatus struct {
	// if the service is ready
	Ready bool `json:"ready,omitempty"`
	// reason a service isn't ready
	Reason string `json:"reason,omitempty"`
	// the time the `ready` flag was last updated
	Timestamp string `json:"timestamp,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// OpenNMS - is the Schema for the opennms API
type OpenNMS struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OpenNMSSpec   `json:"spec,omitempty"`
	Status OpenNMSStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// OpenNMSList - contains a list of OpenNMS
type OpenNMSList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OpenNMS `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OpenNMS{}, &OpenNMSList{})
}
