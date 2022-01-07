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

package main

import (
	"github.com/OpenNMS/opennms-operator/config"
	"github.com/OpenNMS/opennms-operator/internal/reconciler"
	"github.com/OpenNMS/opennms-operator/internal/util/values"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	opennmsv1alpha1 "github.com/OpenNMS/opennms-operator/api/v1alpha1"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(opennmsv1alpha1.AddToScheme(scheme))
}

func main() {

	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	operatorConfig := config.LoadConfig()

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: ":9090",
		Port:               9443,
		LeaderElection:     false,
		LeaderElectionID:   "",
		Namespace:          "", // namespaced-scope when the value is not an empty string
	})
	if err != nil {
		setupLog.Error(err, "unable to define OpenNMS operator")
		os.Exit(1)
	}

	if err = (&reconciler.OpenNMSReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("OpenNMS"),
		Scheme: mgr.GetScheme(),
		CodecFactory: serializer.NewCodecFactory(mgr.GetScheme()),
		Config: operatorConfig,
		DefaultValues: values.LoadValues(operatorConfig.DefaultOpenNMSValuesFile),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create OpenNMS controller", "controller", "OpenNMS")
		os.Exit(1)
	}

	setupLog.Info("starting operator")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem starting operator")
		os.Exit(1)
	}
}