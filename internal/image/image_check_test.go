//go:build unit
// +build unit

package image

import (
	"github.com/OpenNMS/opennms-operator/api/v1alpha1"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
	service := v1.Deployment{}
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
	service := v1.Deployment{}
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
	service := v1.Deployment{}
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
	service := v1.Deployment{}
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
