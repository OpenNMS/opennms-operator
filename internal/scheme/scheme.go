package scheme

import (
	opennmsv1alpha1 "github.com/OpenNMS/opennms-operator/api/v1alpha1"
	certv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	esv1 "github.com/elastic/cloud-on-k8s/pkg/apis/elasticsearch/v1"
	keycloak "github.com/keycloak/keycloak-operator/pkg/apis"
	olm "github.com/operator-framework/api/pkg/operators/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
)

//GetScheme - get the k8s scheme + any custom resources being used
func GetScheme() *runtime.Scheme {
	scheme := runtime.NewScheme()

	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(opennmsv1alpha1.AddToScheme(scheme))
	utilruntime.Must(esv1.AddToScheme(scheme))
	utilruntime.Must(certv1.AddToScheme(scheme))
	utilruntime.Must(keycloak.AddToScheme(scheme))
	utilruntime.Must(olm.AddToScheme(scheme))

	return scheme
}
