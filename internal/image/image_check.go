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
	RunningInstances map[string]v1alpha1.OpenNMS
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
		RunningInstances: map[string]v1alpha1.OpenNMS{},
	}
}

func (ic ImageChecker) InitScheduler() {
	ic.Scheduler = gocron.NewScheduler(time.UTC)
	ic.Scheduler.SetMaxConcurrentJobs(1, gocron.WaitMode)
	ic.Scheduler.StartAsync()
	ic.Scheduler.Every(ic.Frequency).Minute().Do(ic.checkAll)
	//TODO capture job returned by Do() and stop it gracefully when the operator shuts down
}

func (ic *ImageChecker) StartImageCheckerForInstance(instance v1alpha1.OpenNMS) {
	if ic.Scheduler == nil { //lazy scheduler start
		ic.InitScheduler()
	}
	ic.RunningInstances[instance.Name] = instance
}

func (ic *ImageChecker) StopImageCheckerForInstance(instanceName string) {
	delete(ic.RunningInstances, instanceName)
}

func (ic *ImageChecker) checkAll() {
	for iName, instance := range ic.RunningInstances {

	}
}

func (ic *ImageChecker) getImageForInstance(instance v1alpha1.OpenNMS) string {

}
