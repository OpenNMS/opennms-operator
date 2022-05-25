//go:build unit
// +build unit

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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvertCRDToValues(t *testing.T) {
	crd := v1alpha1.OpenNMS{
		Spec: v1alpha1.OpenNMSSpec{
			Namespace: "testns",
			Host:      "testhost",
		},
	}

	res := ConvertCRDToValues(crd, values.TemplateValues{})

	assert.Equal(t, crd.Spec.Namespace, res.Values.Namespace, "should pull the correct value")
	assert.Equal(t, crd.Spec.Host, res.Values.Host, "should pull the correct value")
}

func TestGetCoreValues(t *testing.T) {
	spec := v1alpha1.OpenNMSSpec{
		Core: v1alpha1.BaseServiceResources{
			CPU:  "testcpu",
			MEM:  "testmem",
			Disk: "testdisk",
		},
	}

	v := values.OpenNMSValues{}

	v = getCoreValues(spec, v)

	assert.Equal(t, spec.Core.CPU, v.Core.Resources.Request.Cpu, "should pull the correct value")
	assert.Equal(t, spec.Core.MEM, v.Core.Resources.Request.Memory, "should pull the correct value")
	assert.Equal(t, spec.Core.Disk, v.Core.VolumeSize, "should pull the correct value")

	v = getCoreValues(v1alpha1.OpenNMSSpec{}, v)

	assert.Equal(t, "testcpu", v.Core.Resources.Request.Cpu, "value should remain unchanged when spec is unset")
	assert.Equal(t, "testmem", v.Core.Resources.Request.Memory, "value should remain unchanged when spec is unset")
	assert.Equal(t, "testdisk", v.Core.VolumeSize, "value should remain unchanged when spec is unset")
}

func TestGetAPIValues(t *testing.T) {
	spec := v1alpha1.OpenNMSSpec{
		API: v1alpha1.BaseServiceResources{
			CPU:  "testcpu",
			MEM:  "testmem",
			Disk: "testdisk",
		},
	}

	v := values.OpenNMSValues{}

	v = getAPIValues(spec, v)

	assert.Equal(t, spec.API.CPU, v.API.Resources.Request.Cpu, "should pull the correct value")
	assert.Equal(t, spec.API.MEM, v.API.Resources.Request.Memory, "should pull the correct value")
	assert.Equal(t, spec.API.Disk, v.API.VolumeSize, "should pull the correct value")

	v = getAPIValues(v1alpha1.OpenNMSSpec{}, v)

	assert.Equal(t, "testcpu", v.API.Resources.Request.Cpu, "value should remain unchanged when spec is unset")
	assert.Equal(t, "testmem", v.API.Resources.Request.Memory, "value should remain unchanged when spec is unset")
	assert.Equal(t, "testdisk", v.API.VolumeSize, "value should remain unchanged when spec is unset")
}

func TestGetUIValues(t *testing.T) {
	spec := v1alpha1.OpenNMSSpec{
		UI: v1alpha1.BaseServiceResources{
			CPU:  "testcpu",
			MEM:  "testmem",
			Disk: "testdisk",
		},
	}

	v := values.OpenNMSValues{}

	v = getUIValues(spec, v)

	assert.Equal(t, spec.UI.CPU, v.UI.Resources.Request.Cpu, "should pull the correct value")
	assert.Equal(t, spec.UI.MEM, v.UI.Resources.Request.Memory, "should pull the correct value")
	assert.Equal(t, spec.UI.Disk, v.UI.VolumeSize, "should pull the correct value")

	v = getUIValues(v1alpha1.OpenNMSSpec{}, v)

	assert.Equal(t, "testcpu", v.UI.Resources.Request.Cpu, "value should remain unchanged when spec is unset")
	assert.Equal(t, "testmem", v.UI.Resources.Request.Memory, "value should remain unchanged when spec is unset")
	assert.Equal(t, "testdisk", v.UI.VolumeSize, "value should remain unchanged when spec is unset")
}

func TestGetTimeseriesValues(t *testing.T) {
	spec := v1alpha1.OpenNMSSpec{
		Timeseries: v1alpha1.Timeseries{
			Mode:   "testmode",
			Host:   "testhost",
			Port:   "testport",
			ApiKey: "testkey",
		},
	}

	v := values.TimeseriesValues{}

	v = getTimeseriesValues(spec, v)

	assert.Equal(t, spec.Timeseries.Mode, v.Mode, "should pull the correct value")
	assert.Equal(t, spec.Timeseries.Host, v.Host, "should pull the correct value")
	assert.Equal(t, spec.Timeseries.Port, v.Port, "should pull the correct value")
	assert.Equal(t, spec.Timeseries.ApiKey, v.ApiKey, "should pull the correct value")

	v = getTimeseriesValues(v1alpha1.OpenNMSSpec{}, v)

	assert.Equal(t, "testmode", v.Mode, "value should remain unchanged when spec is unset")
	assert.Equal(t, "testhost", v.Host, "value should remain unchanged when spec is unset")
	assert.Equal(t, "testport", v.Port, "value should remain unchanged when spec is unset")
	assert.Equal(t, "testkey", v.ApiKey, "value should remain unchanged when spec is unset")
}

func TestGetPostgresValues(t *testing.T) {
	spec := v1alpha1.OpenNMSSpec{
		Postgres: v1alpha1.BaseServiceResources{
			Disk: "testing",
		},
	}

	v := values.PostgresValues{}

	v = getPostgresValues(spec, v)

	assert.Equal(t, spec.Postgres.Disk, v.VolumeSize, "should pull the correct value")

	v = getPostgresValues(v1alpha1.OpenNMSSpec{}, v)

	assert.Equal(t, "testing", v.VolumeSize, "value should remain unchanged when spec is unset")
}
