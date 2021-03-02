# Kubernetes Watcher

Kubernetes Watcher is a simple API that returns every `Running` Pod , in every `Namespace`, with all `Containers` and `InitContainers` in the object. 


## Requeriments

Apply `example.yaml`, this was tested in `kubernetes v1.17`. 

## Libraries

I used the [client-go library](https://github.com/kubernetes/client-go), with In-Cluster configuration. 
