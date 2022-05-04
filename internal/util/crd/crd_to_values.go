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

	if spec.TestDeploy {
		v.TestDeploy = spec.TestDeploy
		v = overrideImages(spec, v)
	}

	templateValues.Values = v

	return templateValues
}

//getCoreValues - get ONMS core values from the crd
func getCoreValues(spec v1alpha1.OpenNMSSpec, v values.OpenNMSValues) values.OpenNMSValues {
	v.Image = spec.Version
	if spec.Core.CPU != "" {
		v.Resources.Request.Cpu = spec.Core.CPU
		v.Resources.Limits.Cpu = spec.Core.CPU
	}
	if spec.Core.MEM != "" {
		v.Resources.Request.Memory = spec.Core.MEM
		v.Resources.Limits.Memory = spec.Core.MEM
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

//overrideImages - overrides images with noop images for deployment testing purposes
func overrideImages(spec v1alpha1.OpenNMSSpec, v values.Values) values.Values {
	noopServiceImage := "lipanski/docker-static-website:latest"
	noopJobImage := "alpine:latest"

	v.OpenNMS.Image = noopServiceImage
	v.Postgres.Image = noopServiceImage
	v.Grafana.Image = noopServiceImage
	v.Auth.Image = noopServiceImage
	v.Stunnel.Image = noopServiceImage

	v.Ingress.ControllerImage = noopServiceImage
	v.Ingress.SecretJobImage = noopJobImage
	v.Ingress.WebhookPatchJobImage = noopJobImage

	return v
}

//getPostgresValues - get Postgres DB values from the CRD
func getPostgresValues(spec v1alpha1.OpenNMSSpec, v values.PostgresValues) values.PostgresValues {
	if spec.Postgres.Disk != "" {
		v.VolumeSize = spec.Postgres.Disk
	}
	return v
}
