package image

import (
	"github.com/OpenNMS/opennms-operator/api/v1alpha1"
	"github.com/go-co-op/gocron"
	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

const (
	None             = "none"
	Now              = "now"
	Autoupdate       = "autoupdate"
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
			_ = ic.getImageForService(instance, service)
		}
	}
}

func (ic *ImageChecker) getImageForService(instanceName string, service client.Object) string {
	//ctx := context.Background()
	//var pod corev1.Pod
	//selector := types.Label{
	//	Pairs: map[string]string{
	//
	//	},
	//}
	//ic.Client.Get(ctx, selector, &pod)
	return ""
}
