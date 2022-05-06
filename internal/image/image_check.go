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
	"github.com/OpenNMS/opennms-operator/internal/util/crd"
	"github.com/go-co-op/gocron"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

const (
	None             = "none"
	Now              = "now"
	DefaultFrequency = 60
)

type ImageChecker struct {
	client.Client
	Log              logr.Logger
	Frequency        int
	RunningInstances map[string][]client.Object
	Scheduler        *gocron.Scheduler
}

func NewImageChecker(k8sClient client.Client, freq int) ImageChecker {
	if freq == 0 {
		freq = DefaultFrequency
	}
	return ImageChecker{
		Client:           k8sClient,
		Log:              ctrl.Log.WithName("imageCheck").WithName("OpenNMS"),
		Frequency:        freq,
		RunningInstances: map[string][]client.Object{},
	}
}

func (ic *ImageChecker) InitScheduler() {
	ic.Scheduler = gocron.NewScheduler(time.UTC)
	ic.Scheduler.SetMaxConcurrentJobs(1, gocron.WaitMode)
	ic.Scheduler.StartAsync()
	ic.Scheduler.Every(ic.Frequency).Minute().Do(ic.checkAll)
	//TODO capture job returned by Do() and stop it gracefully when the operator shuts down
}

//StartImageCheckerForInstance - start the recurrent image check for the given instance and services
func (ic *ImageChecker) StartImageCheckerForInstance(instance v1alpha1.OpenNMS, services []client.Object) {
	if ic.Scheduler == nil { //lazy scheduler start
		ic.InitScheduler()
	}
	ic.RunningInstances[instance.Name] = services
}

//StopImageCheckerForInstance - stop the recurrent image check for the given instance and services
func (ic *ImageChecker) StopImageCheckerForInstance(instance v1alpha1.OpenNMS) {
	delete(ic.RunningInstances, instance.Name)
}

//ImageCheckerForInstanceRunning - check if an image check is running for the given instance
func (ic *ImageChecker) ImageCheckerForInstanceRunning(instance v1alpha1.OpenNMS) bool {
	_, ok := ic.RunningInstances[instance.Name]
	return ok
}

//ServiceMarkedForImageCheck - check if a given service is marked for auto updating
func (ic *ImageChecker) ServiceMarkedForImageCheck(service client.Object) bool {
	if autoupdate, ok := service.GetAnnotations()["autoupdate"]; ok {
		if autoupdate == "true" {
			return true
		}
	}
	return false
}

//checkAll - check all registered instances for image updates
func (ic *ImageChecker) checkAll() {
	for instance, services := range ic.RunningInstances {
		//instanceUpdate := false
		for _, service := range services {
			ctx := context.Background()
			imageName, imageId := ic.getImageForService(ctx, instance, service)
			_, err := ic.getLatestImageDigest(ctx, imageName, imageId)
			if err != nil { //if there's an error getting the image digest, skip this service
				continue
			}
		}
	}
}

//getImageForService - get the Image and ImageID for a given service in a given instance
func (ic *ImageChecker) getImageForService(ctx context.Context, instanceName string, service client.Object) (string, string) {
	pod := ic.getPodFromCluster(ctx, instanceName, service)
	//should only be one running container
	return pod.Status.ContainerStatuses[0].Image, pod.Status.ContainerStatuses[0].ImageID
}

//getPodFromCluster - get the first running pod for the given service from the given instance namespace
func (ic *ImageChecker) getPodFromCluster(ctx context.Context, instanceName string, service client.Object) corev1.Pod {
	serviceName := service.GetLabels()["app.kubernetes.io/name"]
	labelMap := map[string]string{
		"app.kubernetes.io/name": serviceName,
	}
	labelsSelector := labels.SelectorFromSet(labelMap)

	fieldSelector := fields.OneTermEqualSelector("status.phase", "Running")

	var pods corev1.PodList
	err := ic.Client.List(ctx, &pods, &client.ListOptions{LabelSelector: labelsSelector, FieldSelector: fieldSelector, Namespace: instanceName})
	if err != nil {
		ic.Log.Error(err, "could not get pod from cluster", "namespace", instanceName, "service", serviceName)
	}
	//only return the first pod of the list - all replicas of a service should be running the same image
	return pods.Items[0]
}

//markInstanceForUpdate -
func (ic *ImageChecker) markInstanceForUpdate(ctx context.Context, instanceName, oldDigest, newDigest string) {
	instance, err := crd.GetInstance(ctx, ic.Client, types.NamespacedName{Namespace: instanceName})
	if err != nil {
		if errors.IsNotFound(err) {
			ic.Log.Info("OpenNMS resource not found", "name", instanceName)
			return
		}
		// Error reading the object - requeue the request.
		ic.Log.Error(err, "Failed to get OpenNMS", "name", instanceName)
		return
	}
	now := time.Now().Format(time.Stamp)

	instance.Status.Image.IsLatest = !(oldDigest == newDigest)
	instance.Status.Image.CheckedAt = now
	if err := ic.Status().Update(ctx, &instance); err != nil {
		ic.Log.Error(err, "Failed to update OpenNMS status", "name", instanceName)
	}
}
