# podchaosmonkey
PodChaosMonkey is a simple controller created to run pod delete operations inside Kubernetes

## Description
The PodChaosMonkey simplifies randomised pod deletions utilising a controller to interact with the api-server in a given
cluster.  It provides the Monkey [CRDs](https://kubernetes.io/docs/tasks/access-kubernetes-api/extend-api-custom-resource-definitions) 
```yaml
apiVersion: podchaos.podchaosmonkey.pt/v1alpha1
kind: Monkey
metadata:
  name: monkey-sample
spec:
  noop: true # choose to log only or run the pod delete operation
  interval: 1m # choose to minimum interval to run operations default: 30s
  namespace: workloads # choose to minimum interval to run operations
  selector: # label selector for choosing the pods to delete
    matchLabels:
      chaosAllowed: "true" #example label
```
Once the **Monkey** Resource is loaded into the cluster, **podchaosmonkey** will add a status condition to indicate that the experiments are active from a given time. At every interval specified a pod matching the search criteria from the cluster will be deleted at random.

## Background
This controller was built using the kubernetes-sig project [kubebuilder](https://github.com/kubernetes-sigs/kubebuilder) project which is an SDK created as part of the kubernetes project as a means to simplify the creation of [custom resource definitions (CRDs)](https://kubernetes.io/docs/tasks/access-kubernetes-api/extend-api-custom-resource-definitions).  As part of this resources such as RBAC and deployment templates are generated to standardise and improve reliability.

This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/) and uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/) providing a reconcile function responsible for synchronizing resources and running the chaos experiments.


## Getting Started
You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

### Running on the cluster

1. Build and push your image to the location specified by `IMG`:
	
```sh
make docker-build docker-push IMG=<some-registry>/podchaosmonkey:tag
```
	
3. Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=<some-registry>/podchaosmonkey:tag
```

4. Deploy a sample workload
```sh
kubectl apply -f config/samples/sample-deployment.yaml
```

5. Create the Monkey resource:

```sh
kubectl apply -f config/samples/podchaos_v1alpha1_monkey.yaml
```

6. Watch the Monkey do his thing!

```sh
watch kubectl get pods -n workloads
```

### Uninstall CRDs
To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller
UnDeploy the controller to the cluster:

```sh
make undeploy
```

**NOTE:** You can also run this in one step by running: `make install run`

### Modifying the API definitions
If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

