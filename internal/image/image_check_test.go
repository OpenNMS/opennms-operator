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
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"
)

func TestImageChecker_InitScheduler(t *testing.T) {
	var ic ImageChecker
	ic.InitScheduler()
	assert.NotNil(t, ic.Scheduler, "scheduler should have been init'd")
	ic.Scheduler.Stop()
}

func TestImageChecker_StartImageCheckerForInstance(t *testing.T) {
	instance := v1alpha1.OpenNMS{}
	instance.SetName("test")
	service := appsv1.Deployment{}
	service.SetName("service")
	services := []client.Object{
		&service,
	}
	ic := NewImageChecker(nil, 60)
	ic.StartImageCheckerForInstance(instance, services)

	res := ic.RunningInstances["test"]

	assert.Equal(t, services, res, "should have started the image checker for the instance")
	ic.Scheduler.Stop()
}

func TestImageChecker_StopImageCheckerForInstance(t *testing.T) {
	instance := v1alpha1.OpenNMS{}
	instance.SetName("test")
	ic := NewImageChecker(nil, 60)
	service := appsv1.Deployment{}
	service.SetName("service")
	services := []client.Object{
		&service,
	}
	ic.RunningInstances = map[string][]client.Object{
		"test": services,
	}
	ic.StopImageCheckerForInstance(instance)
	emptyMap := map[string][]client.Object{}

	assert.Equal(t, emptyMap, ic.RunningInstances, "should have deleted the map entry for the instance")
}

func TestImageChecker_ImageCheckerForInstanceRunning(t *testing.T) {
	instance := v1alpha1.OpenNMS{}
	instance.SetName("test")
	ic := NewImageChecker(nil, 60)
	service := appsv1.Deployment{}
	service.SetName("service")
	services := []client.Object{
		&service,
	}
	ic.RunningInstances = map[string][]client.Object{
		"test": services,
	}
	assert.True(t, ic.ImageCheckerForInstanceRunning(instance), "should recognise that a checker is running for the instance")
}

func TestImageChecker_ServiceMarkedForImageCheck(t *testing.T) {
	service := appsv1.Deployment{}
	service.SetName("service")
	service.SetAnnotations(map[string]string{
		"autoupdate": "true",
	})
	ic := NewImageChecker(nil, 60)
	res := ic.ServiceMarkedForImageCheck(&service)
	assert.True(t, res, "should recognise that a service is marked for autoupdating")

	service.SetAnnotations(map[string]string{
		"autoupdate": "false",
	})
	res = ic.ServiceMarkedForImageCheck(&service)
	assert.False(t, res, "should recognise that a service is not marked for autoupdating")
}

func TestImageCheck_getImageForService(t *testing.T) {
	imageName := "thisistheimage"
	imageID := "imageID"
	serviceName := "serviceName"
	namespaceName := "namespace"
	namespace := corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespaceName,
		},
	}
	mockPod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"app.kubernetes.io/name": serviceName,
			},
			Name:      serviceName,
			Namespace: namespaceName,
		},
		Status: corev1.PodStatus{
			Phase: "Running",
			ContainerStatuses: []corev1.ContainerStatus{
				{
					Image:   imageName,
					ImageID: imageID,
				},
			},
		},
	}
	mockClient := fake.NewClientBuilder().WithObjects(&namespace, &mockPod).Build()

	ic := NewImageChecker(mockClient, 60)
	resImage, resImageID := ic.getImageForService(context.Background(), namespaceName, &mockPod)

	assert.Equal(t, imageName, resImage, "should return the pod's image name")
	assert.Equal(t, imageID, resImageID, "should return the pod's image ID")
}
