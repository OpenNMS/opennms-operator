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

	assert.Equal(t, spec.Core.CPU, v.Resources.Request.Cpu, "should pull the correct value")
	assert.Equal(t, spec.Core.MEM, v.Resources.Request.Memory, "should pull the correct value")
	assert.Equal(t, spec.Core.Disk, v.VolumeSize, "should pull the correct value")

	v = getCoreValues(v1alpha1.OpenNMSSpec{}, v)

	assert.Equal(t, "testcpu", v.Resources.Request.Cpu, "value should remain unchanged when spec is unset")
	assert.Equal(t, "testmem", v.Resources.Request.Memory, "value should remain unchanged when spec is unset")
	assert.Equal(t, "testdisk", v.VolumeSize, "value should remain unchanged when spec is unset")

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

func TestGetImage(t *testing.T) {
	spec := v1alpha1.OpenNMSSpec{
		Version: v1alpha1.Version{
			Distribution: "testdist",
			Tag:          "testtag",
		},
	}

	res := getImage(spec)

	assert.Equal(t, "opennms/testdist:testtag", res, "should create the correct image name")

	res = getImage(v1alpha1.OpenNMSSpec{})

	assert.Equal(t, "opennms/horizon:bleeding", res, "should create default image when version is unset")
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