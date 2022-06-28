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

package reconciler

import (
	"context"
	"github.com/OpenNMS/opennms-operator/api/v1alpha1"
	"github.com/OpenNMS/opennms-operator/internal/model/values"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"
)

func TestUpdateValues(t *testing.T) {
	testNamespace := "testNamespace"
	testHost := "testingHost"
	testValues := values.TemplateValues{
		Values: values.Values{
			Host:      testHost,
			Namespace: testNamespace,
		},
	}

	k8sClient := fake.NewClientBuilder().Build()
	testRecon := OpenNMSReconciler{
		DefaultValues: testValues,
		Client:        k8sClient,
	}

	crd := v1alpha1.OpenNMS{
		ObjectMeta: v1.ObjectMeta{
			Name: testNamespace,
		},
		Spec: v1alpha1.OpenNMSSpec{
			Namespace: testNamespace,
			Host:      testHost,
		},
	}

	res := testRecon.UpdateValues(context.Background(), crd)

	assert.Equal(t, testNamespace, res.Values.Namespace, "should have populated values from reconcile request")
	assert.Equal(t, "testingHost", res.Values.Host, "should have used values from the default values")

	_, ok := testRecon.ValuesMap[testNamespace]
	assert.True(t, ok, "should have saved the created values to the reconciler's values map")
}

func TestCheckForExistingCoreCreds(t *testing.T) {
	testValues := values.TemplateValues{
		Values: values.Values{},
	}
	k8sClient := fake.NewClientBuilder().Build()
	testRecon := OpenNMSReconciler{
		DefaultValues: testValues,
		Client:        k8sClient,
	}
	ctx := context.Background()
	_, resbool := testRecon.CheckForExistingCoreCreds(ctx, testValues, "")
	assert.False(t, resbool, "should return that no core creds existed")

	_, resbool = testRecon.CheckForExistingPostgresCreds(ctx, testValues, "")
	assert.False(t, resbool, "should return that no postgres creds existed")

	adminPwd := "testadminpwd"
	minionPwd := "testminionpwd"
	coreSecret := corev1.Secret{
		ObjectMeta: v1.ObjectMeta{
			Name: "onms-initial-creds",
		},
		Data: map[string][]byte{
			"admin":  []byte(adminPwd),
			"minion": []byte(minionPwd),
		},
	}
	err := k8sClient.Create(ctx, &coreSecret)
	assert.Nil(t, err)

	res, resbool := testRecon.CheckForExistingCoreCreds(ctx, testValues, "")
	assert.True(t, resbool, "should return that there are existing creds")
	//TODO update for keycloak
	//assert.Equal(t, adminPwd, res.Values.Auth.AdminPass, "should return the expected admin password values")
	//assert.Equal(t, minionPwd, res.Values.Auth.MinionPass, "should return the expected admin password values")

	adminPglPwd := "testpostgresadminpwd"
	keycloakPwd := "testpostgreskeycloakpwd"
	pgSecret := corev1.Secret{
		ObjectMeta: v1.ObjectMeta{
			Name: "postgres",
		},
		Data: map[string][]byte{
			"adminPwd":    []byte(adminPglPwd),
			"keycloakPwd": []byte(keycloakPwd),
		},
	}
	err = k8sClient.Create(ctx, &pgSecret)
	assert.Nil(t, err)

	res, resbool = testRecon.CheckForExistingPostgresCreds(ctx, testValues, "")
	assert.True(t, resbool, "should return that there are existing creds")
	assert.Equal(t, adminPglPwd, res.Values.Postgres.AdminPassword, "should return the postgres expected values")
	assert.Equal(t, keycloakPwd, res.Values.Postgres.KeycloakPassword, "should return the postgres expected values")
}
