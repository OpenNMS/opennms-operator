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

package reconciler

import (
	"context"
	"github.com/OpenNMS/opennms-operator/api/v1alpha1"
	"github.com/OpenNMS/opennms-operator/internal/model/values"
	valuesutil "github.com/OpenNMS/opennms-operator/internal/util/crd"
	"github.com/OpenNMS/opennms-operator/internal/util/security"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

//UpdateValues - update values for an instance based on it's crd
func (r *OpenNMSReconciler) UpdateValues(ctx context.Context, instance v1alpha1.OpenNMS) values.TemplateValues {
	if r.ValuesMap == nil {
		r.ValuesMap = map[string]values.TemplateValues{}
	}

	namespace := instance.Namespace

	templateValues, ok := r.ValuesMap[namespace]
	if !ok {
		templateValues = r.DefaultValues
	}

	templateValues = valuesutil.ConvertCRDToValues(instance, templateValues)
	//TODO checking the cluster and reloading creds every iteration doesn't feel great, maybe rework
	templateValues, existingCreds := r.CheckForExistingCreds(ctx, templateValues, namespace)
	if !existingCreds { // only set new passwords if they weren't already created by a previous operator
		templateValues = setPasswords(templateValues)
	}

	r.ValuesMap[namespace] = templateValues

	return templateValues
}

//CheckForExistingCreds - checks if core credentials already exist for a given namespace
func (r *OpenNMSReconciler) CheckForExistingCreds(ctx context.Context, v values.TemplateValues, namespace string) (values.TemplateValues, bool) {
	var credSecret v1.Secret
	err := r.Client.Get(ctx, types.NamespacedName{Namespace: namespace, Name: "onms-initial-creds"}, &credSecret)
	if err != nil {
		return v, false
	}
	v.Values.Auth.AdminPass = string(credSecret.Data["admin"])
	v.Values.Auth.AdminPass = string(credSecret.Data["minion"])

	err = r.Client.Get(ctx, types.NamespacedName{Namespace: namespace, Name: "postgres"}, &credSecret)
	if err != nil {
		return v, false
	}
	v.Values.Postgres.Password = string(credSecret.Data["password"])
	return v, true
}

//setPasswords - sets randomly generated passwords if not already set
func setPasswords(tv values.TemplateValues) values.TemplateValues {
	if tv.Values.Auth.AdminPass == "notset" {
		tv.Values.Auth.AdminPass = security.GeneratePassword()
	}
	if tv.Values.Auth.MinionPass == "notset" {
		tv.Values.Auth.MinionPass = security.GeneratePassword()
	}
	if tv.Values.Postgres.Password == "notset" {
		tv.Values.Postgres.Password = security.GeneratePassword()
	}
	return tv
}

