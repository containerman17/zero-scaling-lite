# Zero scaler for Kubernetes
Simple, straightforward, lightweight. Works on metrics from Nginx Ingress. Depends on Prometheus to retrieve the metrics. 


## How it works 

### Proxy component
On request received, holds it, asks scaler for scaling up to 1 replica.  

### Downscaler
Every minute:
Checks ingress with `zero-scaling\sleep_enabled` metrics every 60 seconds. If ingress has traffic after `zero-scaling\sleep_after` seconds, scales down deployment `zero-scaling\deployment_name` to zero and redirect all traffic to proxy.

### Upscaler
HTTP /wakeup/[domainName]
Rescales the deployment and sends traffic to `zero-scaling\service_name`.
