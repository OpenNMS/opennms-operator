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

package dependencies

import (
	"context"
	"github.com/OpenNMS/opennms-operator/internal/model/values"
	"github.com/go-logr/logr"
	helmclient "github.com/mittwald/go-helm-client"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

//ApplyDependencies - applies the list of helm dependencies to the cluster
func ApplyDependencies(setupLog logr.Logger, values values.TemplateValues, k8sClient client.Client) error {
	helmClient, err := helmclient.New(&helmclient.Options{})
	if err != nil {
		return err
	}
	ctx := context.Background()
	setupLog.Info("Registering Helm repositories")
	err = registerRepositories(helmClient)
	if err != nil {
		return err
	}
	setupLog.Info("Applying Helm charts")
	err = applyCharts(ctx, setupLog, helmClient)
	if err != nil {
		return err
	}
	err = applyConfigs(ctx, setupLog, values, k8sClient)
	if err != nil {
		return err
	}
	return nil
}

//registerRepositories - register dependencies repositories
func registerRepositories(helmClient helmclient.Client) error {
	for _, repository := range repositories {
		err := helmClient.AddOrUpdateChartRepo(repository)
		if err != nil {
			return err
		}
	}
	return nil
}

//applyCharts - apply a given helm dependency
func applyCharts(ctx context.Context, logger logr.Logger, helmClient helmclient.Client) error {
	for _, chart := range charts {
		logger.Info("Applying chart", "chart", chart.ChartName)
		_, err := helmClient.InstallOrUpgradeChart(ctx, &chart)
		if err != nil {
			//same helm chart may have already been installed by someone else in this cluster
			if strings.Contains(err.Error(), "rendered manifests contain a resource that already exists") {
				logger.Info("Resources for chart already exist, skipping", "chart", chart.ChartName)
				continue //not a real error, defer to the existing resources
			}g
			return err
		}
	}
	return nil
}

func applyConfigs(ctx context.Context, logger logr.Logger, v values.TemplateValues, k8sClient client.Client) error {
	for _, handler := range handlerslist {
		for _, resource := range handler.ProvideConfig(v) {
			err := k8sClient.Create(ctx, resource)
			if err != nil {
				logger.Error(err, "Error creating resource", "namespace", resource.GetNamespace(), "name", resource.GetName())
			}
		}
	}
}
