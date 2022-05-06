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

package image

import (
	"context"
	"github.com/OpenNMS/opennms-operator/api/v1alpha1"
	"github.com/OpenNMS/opennms-operator/internal/scheme"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"
)

func TestImageChecker_MarkInstanceServiceForUpdate(t *testing.T) {
	instanceName := "1234"
	serviceName1 := "sn1"
	serviceName2 := "sn2"
	onmsInstance := v1alpha1.OpenNMS{
		ObjectMeta: v1.ObjectMeta{
			Name: instanceName,
		},
	}

	mockClient := fake.NewClientBuilder().WithScheme(scheme.GetScheme()).WithObjects(&onmsInstance).Build()

	ic := NewImageChecker(mockClient, 60)

	ic.markInstanceServiceForUpdate(context.Background(), instanceName, serviceName1, "1", "1")
	mockClient.Get(context.Background(), types.NamespacedName{Name: instanceName}, &onmsInstance)
	assert.True(t, onmsInstance.Status.Image.IsLatest, "should mark instance as on latest")

	ic.markInstanceServiceForUpdate(context.Background(), instanceName, serviceName1, "1", "2")
	mockClient.Get(context.Background(), types.NamespacedName{Name: instanceName}, &onmsInstance)
	assert.False(t, onmsInstance.Status.Image.IsLatest, "should mark service as needing update")
	assert.Equal(t, onmsInstance.Status.Image.ServicesToUpdate, serviceName1, "should insert service name into list of services to update")

	ic.markInstanceServiceForUpdate(context.Background(), instanceName, serviceName2, "1", "2")
	mockClient.Get(context.Background(), types.NamespacedName{Name: instanceName}, &onmsInstance)
	assert.Equal(t, onmsInstance.Status.Image.ServicesToUpdate, "sn1,sn2", "should append extra service to update to list")
}
