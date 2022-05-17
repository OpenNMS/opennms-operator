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
	helmclient "github.com/mittwald/go-helm-client"
	"strings"
)

//ApplyHelmDependencies - applies the list of helm dependencies to the cluster
func ApplyHelmDependencies() error {
	helmClient, err := helmclient.New(&helmclient.Options{})
	if err != nil {
		return err
	}
	ctx := context.Background()
	err = registerRepositories(helmClient)
	if err != nil {
		return err
	}
	err = applyCharts(ctx, helmClient)
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
func applyCharts(ctx context.Context, helmClient helmclient.Client) error {
	for _, chart := range charts {
		_, err := helmClient.InstallOrUpgradeChart(ctx, &chart)
		if err != nil {
			//same helm chart may have already been installed by someone else in this cluster
			if strings.Contains(err.Error(), "rendered manifests contain a resource that already exists") {
				continue //not a real error, defer to the existing resources
			}
			return err
		}
	}
	return nil
}
