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

package reconciler

import (
	"context"
	"github.com/OpenNMS/opennms-operator/api/v1alpha1"
	"github.com/OpenNMS/opennms-operator/internal/model/values"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestUpdateValues(t *testing.T) {
	testValues := values.TemplateValues{
		Values: values.Values{
			Host: "testingHost",
		},
	}

	testNamespace := "testNamespace"

	testRecon := OpenNMSReconciler{
		DefaultValues: testValues,
	}

	crd := v1alpha1.OpenNMS{
		ObjectMeta: v1.ObjectMeta{
			Namespace: testNamespace,
		},
	}

	res := testRecon.UpdateValues(context.Background(), crd)


	assert.Equal(t, testNamespace, res.Values.Namespace, "should have populated values from reconcile request")
	assert.Equal(t, "testingHost", res.Values.Host, "should have used values from the default values")

	_, ok := testRecon.ValuesMap[testNamespace]
	assert.True(t, ok, "should have saved the created values to the reconciler's values map")
}
