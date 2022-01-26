//go:build integration
// +build integration

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

package integration

import (
	"context"
	"github.com/OpenNMS/opennms-operator/internal/scheme"
	"github.com/mittwald/go-helm-client"
	"github.com/stretchr/testify/assert"
	ctrl "sigs.k8s.io/controller-runtime"
	"testing"
)

func TestIntegration(t *testing.T) {
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme.GetScheme(),
		MetricsBindAddress: ":9090",
		Port:               9443,
		LeaderElection:     false,
		LeaderElectionID:   "",
		Namespace:          "", // namespaced-scope when the value is not an empty string
	})
	assert.Nil(t, err, "init k8s client")
	k8sClient := mgr.GetClient()

	helmOptions := &helmclient.Options{}
	helmClient, err := helmclient.New(helmOptions)
	assert.Nil(t, err, "init helm client")

	ctx := context.Background()

	DeployOperatorAndCRD(t, ctx, helmClient)

	DeployInstance(t, ctx, k8sClient)
}
