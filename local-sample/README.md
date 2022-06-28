# Requirements

Kind is installed.

Kubectl is installed.

Running on Mac or Linux.

```
sudo vi /etc/hosts
```
Add with the following to /etc/hosts:
```
127.0.0.1 localhostui
127.0.0.1 localhostkey
127.0.0.1 localhostapi
127.0.0.1 localhostcore
```

# Deploy Update

```
kind create cluster --config=local-sample/config-kind.yaml
make local-docker
kind load docker-image opennms/operator:local-build
docker exec -it  kind-control-plane crictl images # Check images.

helm upgrade -i operator-local ./charts/opennms-operator -f local.yaml
# Wait until the operator is running before moving on.

kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
kubectl apply -f local-instance.yaml

cd local-sample/
kubectl -n local-instance apply -f services.yaml
kubectl -n local-instance apply -f ingress.yaml
kubectl -n local-instance get all

# Needed additional components to get running.  When the following leycloak
# operator is put into its own namespace, it cannot see the KeycloakRealmImport
# (in crd-keycloak.yaml), but if KeycloakRealmImport is put into keycloak
# namespace, it cannot see the deployment of keycloak, my-keycloak, that is in
# local-instance. Need to sort this out.
#kubectl -n keycloak apply -f https://raw.githubusercontent.com/keycloak/keycloak-k8s-resources/18.0.0/kubernetes/keycloaks.k8s.keycloak.org-v1.yml
#kubectl -n keycloak apply -f https://raw.githubusercontent.com/keycloak/keycloak-k8s-resources/18.0.0/kubernetes/keycloakrealmimports.k8s.keycloak.org-v1.yml
#kubectl -n keycloak apply -f https://raw.githubusercontent.com/keycloak/keycloak-k8s-resources/18.0.0/kubernetes/kubernetes.yml
#kubectl -n keycloak rollout restart deployment.apps/keycloak-operator
kubectl delete namespace keycloak
kubectl -n local-instance apply -f https://raw.githubusercontent.com/keycloak/keycloak-k8s-resources/18.0.0/kubernetes/keycloaks.k8s.keycloak.org-v1.yml
kubectl -n local-instance apply -f https://raw.githubusercontent.com/keycloak/keycloak-k8s-resources/18.0.0/kubernetes/keycloakrealmimports.k8s.keycloak.org-v1.yml
kubectl -n local-instance apply -f https://raw.githubusercontent.com/keycloak/keycloak-k8s-resources/18.0.0/kubernetes/kubernetes.yml

# Wait until the above keycloak operator is running. Then run this.
kubectl -n local-instance apply -f secret-postgres-admin.yaml
kubectl -n local-instance apply -f configmap-keycloak.yaml
kubectl -n local-instance rollout restart deployment.apps/postgres
kubectl -n local-instance apply -f deployment-keycloak.yaml
kubectl -n local-instance apply -f services-keycloak.yaml

# Wait until the above postgres is finished redeploying.
#kubectl -n keycloak apply -f ingress-keycloak.yaml # This only triggers the operator when in the same namespace as the operator. But, then it cannot get the instance.
kubectl -n local-instance apply -f ingress-keycloak.yaml
kubectl -n local-instance apply -f crd-keycloak.yaml # May need to rollout restart the my-keycloak deployment.
```

Go to http://localhostui and login with admin:admin.

Contents of local.yaml
```
Operator:
  image: opennms/operator:local-build
  resources:
    limits:
      cpu: 1
      memory: 500Mi
    requests:
      cpu: 1
      memory: 500Mi
TLS:
  Enabled: false
```

# Cleanup

```
# Delete cluster.
kind delete clusters kind

# Confirm all port-forwarding background processes are killed.
ps -axf | grep kubectl
```

