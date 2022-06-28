![](https://www.opennms.com/wp-content/uploads/2021/04/OpenNMS_Horizontal-Logo_Light-BG-retina-website-300x56.png)

# OpenNMS Operator

A Kubernetes operator for deploying and maintaining [Horizon by OpenNMS](https://github.com/OpenNMS/opennms) in the cloud.

### Versioning

This repository follows [Semantic Versioning](https://semver.org/)

### Quick install

Have some sort of local Kubernetes cluster running, i.e. [Docker Desktop](https://docs.docker.com/desktop/kubernetes/), [minikube](https://minikube.sigs.k8s.io/docs/start/), [kind](https://kind.sigs.k8s.io/docs/user/quick-start/), etc.  

[Install the OperatorSDK](https://sdk.operatorframework.io/docs/installation/)

Run the local install script
```
bash deploy-horizon-stream.yaml
```