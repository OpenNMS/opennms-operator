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
	"fmt"
	"github.com/OpenNMS/opennms-operator/api/v1alpha1"
	"github.com/OpenNMS/opennms-operator/internal/model/values"
)


//ConvertCRDToValues - convert an ONMS crd into a set of template values
func ConvertCRDToValues(crd v1alpha1.OpenNMS, defaultValues values.TemplateValues) values.TemplateValues {
	templateValues := defaultValues

	v := templateValues.Values
	spec := crd.Spec

	v.Namespace = spec.Namespace
	v.Host = spec.Host

	//ONMS Core
	v.OpenNMS = getCoreValues(spec, v.OpenNMS)

	//Postgres
	v.Postgres = getPostgresValues(spec, v.Postgres)

	templateValues.Values = v

	return templateValues
}


//getCoreValues - get ONMS core values from the crd
func getCoreValues(spec v1alpha1.OpenNMSSpec, v values.OpenNMSValues) values.OpenNMSValues {
	v.Image = getImage(spec)
	if spec.Core.CPU != "" {
		v.Resources.Request.Cpu = spec.Core.CPU
	}
	if spec.Core.MEM != "" {
		v.Resources.Request.Memory = spec.Core.MEM
	}
	if spec.Core.Disk != "" {
		v.VolumeSize = spec.Core.Disk
	}
	v.Timeseries = getTimeseriesValues(spec, v.Timeseries)
	return v
}


//getTimeseriesValues - get TS DB values from the crd
func getTimeseriesValues(spec v1alpha1.OpenNMSSpec, v values.TimeseriesValues) values.TimeseriesValues {
	if spec.Timeseries.Mode != "" {
		v.Mode = spec.Timeseries.Mode
	}
	if spec.Timeseries.Host != "" {
		v.Host = spec.Timeseries.Host
	}
	if spec.Timeseries.Port != "" {
		v.Port = spec.Timeseries.Port
	}
	if spec.Timeseries.ApiKey != "" {
		v.ApiKey = spec.Timeseries.ApiKey
	}
	return v
}

//getImage - get an ONMS image from the crd
func getImage(spec v1alpha1.OpenNMSSpec) string {
	distro := spec.Version.Distribution
	if distro == "" {
		distro = "horizon"
	}
	imageTag := spec.Version.Tag
	if imageTag == "" {
		imageTag = "bleeding"
	}
	return fmt.Sprintf("opennms/%s:%s", distro, imageTag)
}


//getPostgresValues - get Postgres DB values from the CRD
func getPostgresValues(spec v1alpha1.OpenNMSSpec, v values.PostgresValues) values.PostgresValues {
	if spec.Postgres.Disk != "" {
		v.VolumeSize = spec.Postgres.Disk
	}
	return v
}