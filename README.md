# Zero scaler for Kubernetes
Simple, straightforward, lightweight. Works on metrics from Nginx Ingress. Depends on Prometheus to retrieve the metrics. 

## How it works 

### Proxy component
On request received, holds it, asks scaler for scaling up to 1 replica.  

### Downscaler
Every minute:
Checks ingress with `zero-scaling/enabled = "true"` annotation every 60 seconds. If ingress has traffic after 15 minutes, scales down deployment `zero-scaling/deploymentName` to zero and redirect all traffic to proxy. Original proxy value saved in 

### proxy
Holds requests and redirects them to original service as long as its ready. 

### prometheus proxy

This setup requires prometheus proxy service

```yaml
kind: Service
apiVersion: v1
metadata:
    name: prometheus-proxy
    namespace: downscaler-kubebuider-system
spec:
    type: ExternalName
    # TODO set externalName to your prometheus location here
    # type "kubectl get svc --all-namespaces | grep prometheus" to fund out
    externalName: prometheus-server.ingress-nginx.svc.cluster.local  #this is just an example
    ports:
    - port: 80
```